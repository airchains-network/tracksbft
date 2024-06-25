package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/airchains-network/tracksbft/abci/server"
	"github.com/airchains-network/tracksbft/config"
	"github.com/airchains-network/tracksbft/crypto/ed25519"
	cmtflags "github.com/airchains-network/tracksbft/libs/cli/flags"
	"github.com/airchains-network/tracksbft/libs/log"
	cmtnet "github.com/airchains-network/tracksbft/libs/net"
	"github.com/airchains-network/tracksbft/light"
	lproxy "github.com/airchains-network/tracksbft/light/proxy"
	lrpc "github.com/airchains-network/tracksbft/light/rpc"
	dbs "github.com/airchains-network/tracksbft/light/store/db"
	"github.com/airchains-network/tracksbft/node"
	"github.com/airchains-network/tracksbft/p2p"
	"github.com/airchains-network/tracksbft/privval"
	"github.com/airchains-network/tracksbft/proxy"
	rpcserver "github.com/airchains-network/tracksbft/rpc/jsonrpc/server"
	"github.com/airchains-network/tracksbft/test/e2e/app"
	e2e "github.com/airchains-network/tracksbft/test/e2e/pkg"
	mcs "github.com/airchains-network/tracksbft/test/maverick/consensus"
	maverick "github.com/airchains-network/tracksbft/test/maverick/node"
)

var logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))

// main is the binary entrypoint.
func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v <configfile>", os.Args[0])
		return
	}
	configFile := ""
	if len(os.Args) == 2 {
		configFile = os.Args[1]
	}

	if err := run(configFile); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// run runs the application - basically like main() with error handling.
func run(configFile string) error {
	cfg, err := LoadConfig(configFile)
	if err != nil {
		return err
	}

	// Start remote signer (must start before node if running builtin).
	if cfg.PrivValServer != "" {
		if err = startSigner(cfg); err != nil {
			return err
		}
		if cfg.Protocol == "builtin" {
			time.Sleep(1 * time.Second)
		}
	}

	// Start app server.
	switch cfg.Protocol {
	case "socket", "grpc":
		err = startApp(cfg)
	case "builtin":
		if len(cfg.Misbehaviors) == 0 {
			if cfg.Mode == string(e2e.ModeLight) {
				err = startLightClient(cfg)
			} else {
				err = startNode(cfg)
			}
		} else {
			err = startMaverick(cfg)
		}
	default:
		err = fmt.Errorf("invalid protocol %q", cfg.Protocol)
	}
	if err != nil {
		return err
	}

	// Apparently there's no way to wait for the server, so we just sleep
	for {
		time.Sleep(1 * time.Hour)
	}
}

// startApp starts the application server, listening for connections from CometBFT.
func startApp(cfg *Config) error {
	app, err := app.NewApplication(cfg.App())
	if err != nil {
		return err
	}
	server, err := server.NewServer(cfg.Listen, cfg.Protocol, app)
	if err != nil {
		return err
	}
	err = server.Start()
	if err != nil {
		return err
	}
	logger.Info("start app", "msg", log.NewLazySprintf("Server listening on %v (%v protocol)", cfg.Listen, cfg.Protocol))
	return nil
}

// startNode starts a CometBFT node running the application directly. It assumes the CometBFT
// configuration is in $CMTHOME/config/cometbft.toml.
//
// FIXME There is no way to simply load the configuration from a file, so we need to pull in Viper.
func startNode(cfg *Config) error {
	app, err := app.NewApplication(cfg.App())
	if err != nil {
		return err
	}

	cmtcfg, nodeLogger, nodeKey, err := setupNode()
	if err != nil {
		return fmt.Errorf("failed to setup config: %w", err)
	}

	n, err := node.NewNode(cmtcfg,
		privval.LoadOrGenFilePV(cmtcfg.PrivValidatorKeyFile(), cmtcfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(cmtcfg),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cmtcfg.Instrumentation),
		nodeLogger,
	)
	if err != nil {
		return err
	}
	return n.Start()
}

func startLightClient(cfg *Config) error {
	cmtcfg, nodeLogger, _, err := setupNode()
	if err != nil {
		return err
	}

	dbContext := &node.DBContext{ID: "light", Config: cmtcfg}
	lightDB, err := node.DefaultDBProvider(dbContext)
	if err != nil {
		return err
	}

	providers := rpcEndpoints(cmtcfg.P2P.PersistentPeers)

	c, err := light.NewHTTPClient(
		context.Background(),
		cfg.ChainID,
		light.TrustOptions{
			Period: cmtcfg.StateSync.TrustPeriod,
			Height: cmtcfg.StateSync.TrustHeight,
			Hash:   cmtcfg.StateSync.TrustHashBytes(),
		},
		providers[0],
		providers[1:],
		dbs.New(lightDB, "light"),
		light.Logger(nodeLogger),
	)
	if err != nil {
		return err
	}

	rpccfg := rpcserver.DefaultConfig()
	rpccfg.MaxBodyBytes = cmtcfg.RPC.MaxBodyBytes
	rpccfg.MaxHeaderBytes = cmtcfg.RPC.MaxHeaderBytes
	rpccfg.MaxOpenConnections = cmtcfg.RPC.MaxOpenConnections
	// If necessary adjust global WriteTimeout to ensure it's greater than
	// TimeoutBroadcastTxCommit.
	// See https://github.com/airchains-network/tracksbft/issues/3435
	if rpccfg.WriteTimeout <= cmtcfg.RPC.TimeoutBroadcastTxCommit {
		rpccfg.WriteTimeout = cmtcfg.RPC.TimeoutBroadcastTxCommit + 1*time.Second
	}

	p, err := lproxy.NewProxy(c, cmtcfg.RPC.ListenAddress, providers[0], rpccfg, nodeLogger,
		lrpc.KeyPathFn(lrpc.DefaultMerkleKeyPathFn()))
	if err != nil {
		return err
	}

	logger.Info("Starting proxy...", "laddr", cmtcfg.RPC.ListenAddress)
	if err := p.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logger.Error("proxy ListenAndServe", "err", err)
	}

	return nil
}

// FIXME: Temporarily disconnected maverick until it is redesigned
// startMaverick starts a Maverick node that runs the application directly. It assumes the CometBFT
// configuration is in $CMTHOME/config/cometbft.toml.
func startMaverick(cfg *Config) error {
	app, err := app.NewApplication(cfg.App())
	if err != nil {
		return err
	}

	cmtcfg, logger, nodeKey, err := setupNode()
	if err != nil {
		return fmt.Errorf("failed to setup config: %w", err)
	}

	misbehaviors := make(map[int64]mcs.Misbehavior, len(cfg.Misbehaviors))
	for heightString, misbehaviorString := range cfg.Misbehaviors {
		height, _ := strconv.ParseInt(heightString, 10, 64)
		misbehaviors[height] = mcs.MisbehaviorList[misbehaviorString]
	}

	n, err := maverick.NewNode(cmtcfg,
		maverick.LoadOrGenFilePV(cmtcfg.PrivValidatorKeyFile(), cmtcfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		maverick.DefaultGenesisDocProviderFunc(cmtcfg),
		maverick.DefaultDBProvider,
		maverick.DefaultMetricsProvider(cmtcfg.Instrumentation),
		logger,
		misbehaviors,
	)
	if err != nil {
		return err
	}

	return n.Start()
}

// startSigner starts a signer server connecting to the given endpoint.
func startSigner(cfg *Config) error {
	filePV := privval.LoadFilePV(cfg.PrivValKey, cfg.PrivValState)

	protocol, address := cmtnet.ProtocolAndAddress(cfg.PrivValServer)
	var dialFn privval.SocketDialer
	switch protocol {
	case "tcp":
		dialFn = privval.DialTCPFn(address, 3*time.Second, ed25519.GenPrivKey())
	case "unix":
		dialFn = privval.DialUnixFn(address)
	default:
		return fmt.Errorf("invalid privval protocol %q", protocol)
	}

	endpoint := privval.NewSignerDialerEndpoint(logger, dialFn,
		privval.SignerDialerEndpointRetryWaitInterval(1*time.Second),
		privval.SignerDialerEndpointConnRetries(100))
	err := privval.NewSignerServer(endpoint, cfg.ChainID, filePV).Start()
	if err != nil {
		return err
	}
	logger.Info("start signer", "msg", log.NewLazySprintf("Remote signer connecting to %v", cfg.PrivValServer))
	return nil
}

func setupNode() (*config.Config, log.Logger, *p2p.NodeKey, error) {
	var cmtcfg *config.Config

	home := os.Getenv("CMTHOME")
	if home == "" {
		return nil, nil, nil, errors.New("CMTHOME not set")
	}

	viper.AddConfigPath(filepath.Join(home, "config"))
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, nil, nil, err
	}

	cmtcfg = config.DefaultConfig()

	if err := viper.Unmarshal(cmtcfg); err != nil {
		return nil, nil, nil, err
	}

	cmtcfg.SetRoot(home)

	if err := cmtcfg.ValidateBasic(); err != nil {
		return nil, nil, nil, fmt.Errorf("error in config file: %w", err)
	}

	if cmtcfg.LogFormat == config.LogFormatJSON {
		logger = log.NewTMJSONLogger(log.NewSyncWriter(os.Stdout))
	}

	nodeLogger, err := cmtflags.ParseLogLevel(cmtcfg.LogLevel, logger, config.DefaultLogLevel)
	if err != nil {
		return nil, nil, nil, err
	}

	nodeLogger = nodeLogger.With("module", "main")

	nodeKey, err := p2p.LoadOrGenNodeKey(cmtcfg.NodeKeyFile())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load or gen node key %s: %w", cmtcfg.NodeKeyFile(), err)
	}

	return cmtcfg, nodeLogger, nodeKey, nil
}

// rpcEndpoints takes a list of persistent peers and splits them into a list of rpc endpoints
// using 26657 as the port number
func rpcEndpoints(peers string) []string {
	arr := strings.Split(peers, ",")
	endpoints := make([]string, len(arr))
	for i, v := range arr {
		urlString := strings.SplitAfter(v, "@")[1]
		hostName := strings.Split(urlString, ":26656")[0]
		// use RPC port instead
		port := 26657
		rpcEndpoint := "http://" + hostName + ":" + fmt.Sprint(port)
		endpoints[i] = rpcEndpoint
	}
	return endpoints
}
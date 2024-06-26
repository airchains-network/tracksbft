# TracksBFT

This project is a fork of [CometBFT v0.34.29](https://github.com/cometbft/cometbft/tree/v0.34.29), which itself is a fork of Tendermint.

We have added some new features to support the creation and management of tracks, including the functionality to create "pods" (sets of X transactions), these data tracks can be accessed via port 26657.

## Tracks System

A unique aspect of this project is the "tracks" system, which is used to construct Zero-Knowledge (ZK) proofs from blockchain transactions. These proofs are then submitted to the Switchyard blockchain, developed by [Airchains.io](https://airchains.io).

You can find more information about the Switchyard blockchain [here](https://github.com/airchains-network/junction) and Tracks [here](https://github.com/airchains-network/tracks).
## How to Use

To use this project as a library in your blockchain project, update the `go.mod` file as follows:

```go
github.com/tendermint/tendermint => github.com/airchains-network/tracksbft v0.0.1
```

Alternatively, you can run the following command and then replace all `tendermint/cometbft` library imports in your blockchain project with `tracksbft`:

```bash
go get github.com/tendermint/tendermint
```

## Minimum Requirements

| Requirement | Notes             |
|-------------|-------------------|
| Go version  | Go 1.19 or higher |

## License

This project is licensed under the MIT License by [Airchains.io](https://airchains.io).

## Maintainers

This project is maintained by [Airchains.io](https://airchains.io).

Feel free to reach out for any further assistance or inquiries about integrating TracksBFT into your blockchain project.

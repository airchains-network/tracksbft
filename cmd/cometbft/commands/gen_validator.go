package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	cmtjson "github.com/airchains-network/tracksbft/libs/json"
	"github.com/airchains-network/tracksbft/privval"
)

// GenValidatorCmd allows the generation of a keypair for a
// validator.
var GenValidatorCmd = &cobra.Command{
	Use:     "gen-validator",
	Aliases: []string{"gen_validator"},
	Short:   "Generate new validator keypair",
	Run:     genValidator,
}

func genValidator(cmd *cobra.Command, args []string) {
	pv := privval.GenFilePV("", "")
	jsbz, err := cmtjson.Marshal(pv)
	if err != nil {
		panic(err)
	}
	fmt.Printf(`%v
`, string(jsbz))
}
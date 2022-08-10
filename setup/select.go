package setup

import (
	"fmt"

	"github.com/Moonyongjung/multisig-gentx/util"
)

var keyAlgo string

//-- Select chain's key algorithm
func SelectChainType() {
	for {
		var s string
		util.LogTool("=======================================================================")
		util.LogTool("                     Initialize GenTx Key Type                         ")
		util.LogTool("                                                                       ")
		util.LogTool("SELECT the chain type                                                  ")
		util.LogTool("                                                                       ")
		util.LogTool("A chain includes ethermint module uses eth_secp256k1 key algorithm     ")
		util.LogTool("but general chains based on cosmos sdk uses secp256k1 algo.            ")
		util.LogTool("If you want to use the chain includes ethermint module,                ")
		util.LogTool("input [y/N]                                                            ")
		util.LogTool("                                                                       ")
		util.LogTool("y. Use includes ethermint module, N. general Cosmos SDK based chain    ")
		util.LogTool("=======================================================================")

		fmt.Scan(&s)

		if s == "y" {
			keyAlgo = "eth_secp256k1"
			return
		} else if s == "N" {
			keyAlgo = "secp256k1"
			return
		} else {
			util.LogTool("Input correct string [y/N]")
		}
	}
}
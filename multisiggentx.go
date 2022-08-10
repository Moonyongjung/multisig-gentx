package main

import (
	"github.com/Moonyongjung/multisig-gentx/pro"
	"github.com/Moonyongjung/multisig-gentx/setup"
	"github.com/Moonyongjung/multisig-gentx/util"
)

var configPath = "./config/config.json"

func init() {
	util.GetConfig().Read(configPath)
	setup.SetChainPrefixConfig()
}

func main() {
	pro.ProGenTx(setup.SetClient())
}

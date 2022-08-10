package pro

import (
	"strconv"

	"github.com/Moonyongjung/multisig-gentx/gentx"
	"github.com/Moonyongjung/multisig-gentx/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ProGenTx(clientCtx cmclient.Context) {
	accountName1 := util.GetConfig().Get("accountName1")
	accountName2 := util.GetConfig().Get("accountName2")
	accountName3 := util.GetConfig().Get("accountName3")
	multisigAccountName := util.GetConfig().Get("multisigAccountName")

	infokey1, err := clientCtx.Keyring.Key(accountName1)
	if err != nil {
		util.LogErr(err)
	}
	infokey2, err := clientCtx.Keyring.Key(accountName2)
	if err != nil {
		util.LogErr(err)
	}
	infokey3, err := clientCtx.Keyring.Key(accountName3)
	if err != nil {
		util.LogErr(err)
	}

	pubs := []cryptotypes.PubKey{infokey1.GetPubKey(), infokey2.GetPubKey(), infokey3.GetPubKey()}

	keyNames := []string{infokey1.GetName(), infokey2.GetName(), infokey3.GetName()}

	pub4key, _ := multisigFunction(clientCtx, keyNames)
	err = gentx.AddGenesisAccount(clientCtx, multisigAccountName)
	if err != nil {
		util.LogErr(err)
	}

	for i, pub := range pubs {
		bech32Account, err := sdk.AccAddressFromHex(pub.Address().String())
		if err != nil {
			util.LogErr(err)
		}
		util.LogTool(keyNames[i], "address : ", bech32Account.String())

	}

	bech32MultisigAddr, err := sdk.AccAddressFromHex(pub4key.Address().String())
	if err != nil {
		util.LogErr(err)
	}
	
	util.LogTool(multisigAccountName, "address : ", bech32MultisigAddr.String())

	var addrs []string
	for _, pub := range pubs {
		addrConv, err := sdk.AccAddressFromHex(pub.Address().String())
		if err != nil {
			util.LogErr(err)
		}
		addrs = append(addrs, addrConv.String())
	}

	err = gentx.GenTxSdk(clientCtx)
	if err != nil {
		util.LogErr(err)
	}

	for i, addr := range addrs {
		index := strconv.Itoa(i + 1)
		err := gentx.TxSign(clientCtx, addr, bech32MultisigAddr.String(), index)
		if err != nil {
			util.LogErr(err)
		}
	}

	err = gentx.MultiSign(clientCtx, keyNames)
	if err != nil {
		util.LogErr(err)
	}
}

package setup

import (
	"github.com/Moonyongjung/multisig-gentx/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
)

func NewFactory(clientCtx cmclient.Context) tx.Factory {

	defaultAccNum := util.GetConfig().Get("defaultAccNum")
	defaultAccSeq := util.GetConfig().Get("defaultAccSeq")
	chainId := util.GetConfig().Get("chainId")
	gasLimit := util.GetConfig().Get("gasLimit")

	accNum := util.FromStringToUint64(defaultAccNum)
	accSeq := util.FromStringToUint64(defaultAccSeq)
	gasLimituint64 := util.FromStringToUint64(gasLimit)
	memo := ""
	signMode := signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON

	txFactory := tx.Factory{}.
		WithTxConfig(clientCtx.TxConfig).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithKeybase(clientCtx.Keyring).
		WithChainID(clientCtx.ChainID).
		WithAccountNumber(accNum).
		WithSequence(accSeq).
		WithMemo(memo).
		WithSignMode(signMode).
		WithChainID(chainId).
		WithGas(gasLimituint64)

	return txFactory
}

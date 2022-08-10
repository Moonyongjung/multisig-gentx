package gentx

import (
	"github.com/Moonyongjung/multisig-gentx/setup"
	"github.com/Moonyongjung/multisig-gentx/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
)

func TxSign(clientCtx cmclient.Context, fromName string, multisig string, index string) error {
	gentxsDir := util.GetConfig().Get("gentxsDir")
	unsignedTxFileName := util.GetConfig().Get("unsignedTxFileName")

	clientCtx, txF, newTx, err := readTxAndInitContexts(clientCtx, gentxsDir+unsignedTxFileName)
	if err != nil {
		return err
	}

	txFactory := txF
	txCfg := clientCtx.TxConfig
	txBuilder, err := txCfg.WrapTxBuilder(newTx)
	if err != nil {
		return err
	}

	from := fromName
	_, fromNameRes, _, err := cmclient.GetFromFields(txF.Keybase(), from, clientCtx.GenerateOnly)
	if err != nil {
		return util.LogErr("error each key getting account from keybase: %w", err)
	}

	multisigAddr, err := sdk.AccAddressFromBech32(multisig)
	if err != nil {
		multisigAddr, _, _, err = cmclient.GetFromFields(txFactory.Keybase(), multisig, clientCtx.GenerateOnly)
		if err != nil {
			return util.LogErr("error multisig key getting account from keybase: %w", err)
		}
	}
	err = authclient.SignTxWithSignerAddress(
		txF, clientCtx, multisigAddr, fromNameRes, txBuilder, clientCtx.Offline, false)
	if err != nil {
		return err
	}
	printSignatureOnly := true

	jsonByte, err := marshalSignatureJSON(txCfg, txBuilder, printSignatureOnly)
	if err != nil {
		return err
	}	

	util.SaveJsonPretty(jsonByte, gentxsDir+"signedTx"+index+".json")

	return nil
}

func readTxAndInitContexts(clientCtx cmclient.Context, filename string) (cmclient.Context, tx.Factory, sdk.Tx, error) {
	stdTx, err := authclient.ReadTxFromFile(clientCtx, filename)
	if err != nil {
		return clientCtx, tx.Factory{}, nil, err
	}

	txFactory := setup.NewFactory(clientCtx)

	return clientCtx, txFactory, stdTx, nil
}

func marshalSignatureJSON(txConfig cmclient.TxConfig, txBldr cmclient.TxBuilder, signatureOnly bool) ([]byte, error) {
	parsedTx := txBldr.GetTx()
	if signatureOnly {
		sigs, err := parsedTx.GetSignaturesV2()
		if err != nil {
			return nil, err
		}
		return txConfig.MarshalSignatureJSON(sigs)
	}

	return txConfig.TxJSONEncoder()(parsedTx)
}

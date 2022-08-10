package gentx

import (	
	"io/ioutil"
	"strconv"

	"github.com/Moonyongjung/multisig-gentx/setup"
	"github.com/Moonyongjung/multisig-gentx/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	kmultisig "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
)

func MultiSign(clientCtx cmclient.Context, keyNames []string) error {
	gentxsDir := util.GetConfig().Get("gentxsDir")
	unsignedTxFileName := util.GetConfig().Get("unsignedTxFileName")
	multisignedTxFileName := util.GetConfig().Get("multisignedTxFileName")
	multisigAccountName := util.GetConfig().Get("multisigAccountName")

	parsedTx, err := authclient.ReadTxFromFile(clientCtx, gentxsDir+unsignedTxFileName)
	if err != nil {
		return err
	}
	txFactory := setup.NewFactory(clientCtx)

	txCfg := clientCtx.TxConfig
	txBuilder, err := txCfg.WrapTxBuilder(parsedTx)
	if err != nil {
		return err
	}

	multisigInfo, err := getMultisigInfo(clientCtx, multisigAccountName)
	if err != nil {
		return err
	}

	multisigPub := multisigInfo.GetPubKey().(*kmultisig.LegacyAminoPubKey)
	multisigSig := multisig.NewMultisig(len(multisigPub.PubKeys))
	if !clientCtx.Offline {
		accnum, seq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, multisigInfo.GetAddress())
		if err != nil {
			return err
		}

		txFactory = txFactory.WithAccountNumber(accnum).WithSequence(seq)
	}

	for i := 0; i < len(keyNames); i++ {
		counterStr := strconv.Itoa(i + 1)
		sigs, err := unmarshalSignatureJSON(clientCtx, gentxsDir+"signedTx"+counterStr+".json")
		if err != nil {
			return err
		}

		if txFactory.ChainID() == "" {
			return util.LogErr("set the chain id with either the --chain-id flag or config file")
		}

		signingData := signing.SignerData{
			ChainID:       txFactory.ChainID(),
			AccountNumber: txFactory.AccountNumber(),
			Sequence:      txFactory.Sequence(),
		}

		for _, sig := range sigs {
			err = signing.VerifySignature(sig.PubKey, signingData, sig.Data, txCfg.SignModeHandler(), txBuilder.GetTx())
			if err != nil {
				addr, _ := sdk.AccAddressFromHex(sig.PubKey.Address().String())
				return util.LogErr("couldn't verify signature for address %s", addr)
			}

			if err := multisig.AddSignatureV2(multisigSig, sig, multisigPub.GetPubKeys()); err != nil {
				return err
			}
		}
	}

	sigV2 := signingtypes.SignatureV2{
		PubKey:   multisigPub,
		Data:     multisigSig,
		Sequence: txFactory.Sequence(),
	}

	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		return err
	}

	jsonByte, err := marshalSignatureJSON(txCfg, txBuilder, false)
	if err != nil {
		return err
	}	

	util.SaveJsonPretty(jsonByte, gentxsDir+multisignedTxFileName)

	return nil
}

func getMultisigInfo(clientCtx cmclient.Context, name string) (keyring.Info, error) {
	kb := clientCtx.Keyring
	multisigInfo, err := kb.Key(name)
	if err != nil {
		return nil, errors.Wrap(err, "error getting keybase multisig account")
	}
	if multisigInfo.GetType() != keyring.TypeMulti {
		return nil, util.LogErr("%q must be of type %s: %s", name, keyring.TypeMulti, multisigInfo.GetType())
	}

	return multisigInfo, nil
}

func unmarshalSignatureJSON(clientCtx cmclient.Context, filename string) (sigs []signingtypes.SignatureV2, err error) {
	var bytes []byte
	if bytes, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	return clientCtx.TxConfig.UnmarshalSignatureJSON(bytes)
}

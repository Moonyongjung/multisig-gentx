package pro

import (
	"strconv"

	"github.com/Moonyongjung/multisig-gentx/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func multisigFunction(clientCtx cmclient.Context,
	multisigKeys []string) (*multisig.LegacyAminoPubKey, sdk.AccAddress) {
	multisigAccountName := util.GetConfig().Get("multisigAccountName")
	multisigThresholdStr := util.GetConfig().Get("multisigThreshold")

	kb := clientCtx.Keyring

	if len(multisigKeys) != 0 {
		pks := make([]cryptotypes.PubKey, len(multisigKeys))
		multisigThreshold, err := strconv.Atoi(multisigThresholdStr)
		if err != nil {
			util.LogErr(err)
		}
		if err := validateMultisigThreshold(multisigThreshold, len(multisigKeys)); err != nil {
			util.LogErr(err)
		}

		for i, keyname := range multisigKeys {
			k, err := kb.Key(keyname)
			if err != nil {
				util.LogErr(err)
			}

			pks[i] = k.GetPubKey()
		}

		pk := multisig.NewLegacyAminoPubKey(multisigThreshold, pks)
		_, err = kb.SaveMultisig(multisigAccountName, pk)
		if err != nil {
			util.LogErr(err)
		}
		return pk, nil
	}
	return nil, nil
}

func validateMultisigThreshold(k, nKeys int) error {
	if k <= 0 {
		return util.LogErr("threshold must be a positive integer")
	}
	if nKeys < k {
		return util.LogErr(
			"threshold k of n multisignature: %d < %d", nKeys, k)
	}
	return nil
}

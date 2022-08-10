package setup

import (
	"github.com/Moonyongjung/multisig-gentx/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func SetChainPrefixConfig() {
	coinType := util.GetConfig().Get("coinType")
	Bech32MainPrefix := util.GetConfig().Get("bech32MainPrefix")
	CoinType := util.FromStringToUint64(coinType)
	FullFundraiserPath := "m/44'/" + coinType + "'/0'/0/0"

	var (
		// PrefixValidator is the prefix for validator keys
		PrefixValidator = "val"
		// PrefixConsensus is the prefix for consensus keys
		PrefixConsensus = "cons"
		// PrefixPublic is the prefix for public keys
		PrefixPublic = "pub"
		// PrefixOperator is the prefix for operator keys
		PrefixOperator = "oper"

		// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
		Bech32PrefixAccAddr = Bech32MainPrefix
		// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
		Bech32PrefixAccPub = Bech32MainPrefix + PrefixPublic
		// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
		Bech32PrefixValAddr = Bech32MainPrefix + PrefixValidator + PrefixOperator
		// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
		Bech32PrefixValPub = Bech32MainPrefix + PrefixValidator + PrefixOperator + PrefixPublic
		// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
		Bech32PrefixConsAddr = Bech32MainPrefix + PrefixValidator + PrefixConsensus
		// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
		Bech32PrefixConsPub = Bech32MainPrefix + PrefixValidator + PrefixConsensus + PrefixPublic
	)

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.SetCoinType(uint32(CoinType))
	config.SetFullFundraiserPath(FullFundraiserPath)
	config.Seal()
}

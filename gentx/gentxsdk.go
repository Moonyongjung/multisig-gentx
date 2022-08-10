package gentx

import (
	"bufio"	
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/Moonyongjung/multisig-gentx/setup"
	"github.com/Moonyongjung/multisig-gentx/util"

	"github.com/pkg/errors"
	tmtypes "github.com/tendermint/tendermint/types"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

var (
	DefaultTokens                  = sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)
	defaultAmount                  = DefaultTokens.String() + sdk.DefaultBondDenom
	defaultCommissionRate          = "0.1"
	defaultCommissionMaxRate       = "0.2"
	defaultCommissionMaxChangeRate = "0.01"
	defaultMinSelfDelegation       = "1"
)

func GenTxSdk(clientCtx cmclient.Context) error {
	multisigAccountName := util.GetConfig().Get("multisigAccountName")
	moniker := util.GetConfig().Get("moniker")
	gentxAmount := util.GetConfig().Get("gentxAmount")

	serverCtx := server.NewDefaultContext()

	cdc := clientCtx.Codec

	config := serverCtx.Config
	config.SetRoot(clientCtx.HomeDir)

	nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(serverCtx.Config)
	if err != nil {
		return errors.Wrap(err, "failed to initialize node validator files")
	}

	genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
	if err != nil {
		return errors.Wrapf(err, "failed to read genesis doc file %s", config.GenesisFile())
	}

	var genesisState map[string]json.RawMessage
	if err = json.Unmarshal(genDoc.AppState, &genesisState); err != nil {
		return errors.Wrap(err, "failed to unmarshal genesis state")
	}

	inBuf := bufio.NewReader(os.Stdin)
	name := multisigAccountName
	key, err := clientCtx.Keyring.Key(name)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch '%s' from the keyring", name)
	}

	ip, _ := server.ExternalIP()
	createValCfg, err := PrepareConfigForTxCreateValidator(ip, moniker, nodeID, genDoc.ChainID, valPubKey)

	amount := gentxAmount
	coins, err := sdk.ParseCoinsNormalized(amount)
	if err != nil {
		return errors.Wrap(err, "failed to parse coins")
	}
	createValCfg.Amount = amount

	err = genutil.ValidateAccountInGenesis(genesisState, banktypes.GenesisBalancesIterator{}, key.GetAddress(), coins, cdc)
	if err != nil {
		return errors.Wrap(err, "failed to validate account in genesis")
	}

	txFactory := setup.NewFactory(clientCtx)

	clientCtx = clientCtx.WithInput(inBuf).WithFromAddress(key.GetAddress())

	txBldr, msg, err := cli.BuildCreateValidatorMsg(clientCtx, createValCfg, txFactory, true)
	if err != nil {
		return errors.Wrap(err, "failed to build create-validator message")
	}

	GenerateTx(clientCtx, txBldr, []sdk.Msg{msg}...)

	return nil

}

func PrepareConfigForTxCreateValidator(ip string, moniker, nodeID, chainID string, valPubKey cryptotypes.PubKey) (cli.TxCreateValidatorConfig, error) {
	website := util.GetConfig().Get("website")
	securityContact := util.GetConfig().Get("securityContact")
	identity := util.GetConfig().Get("identity")

	c := cli.TxCreateValidatorConfig{}
	c.IP = ip
	c.Website = website
	c.SecurityContact = securityContact
	c.Identity = identity
	c.NodeID = nodeID
	c.PubKey = valPubKey
	c.ChainID = chainID
	c.Moniker = moniker

	if c.Amount == "" {
		c.Amount = defaultAmount
	}

	if c.CommissionRate == "" {
		c.CommissionRate = defaultCommissionRate
	}

	if c.CommissionMaxRate == "" {
		c.CommissionMaxRate = defaultCommissionMaxRate
	}

	if c.CommissionMaxChangeRate == "" {
		c.CommissionMaxChangeRate = defaultCommissionMaxChangeRate
	}

	if c.MinSelfDelegation == "" {
		c.MinSelfDelegation = defaultMinSelfDelegation
	}

	return c, nil
}

func readUnsignedGenTxFile(clientCtx cmclient.Context, r io.Reader) (sdk.Tx, error) {
	bz, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	aTx, err := clientCtx.TxConfig.TxJSONDecoder()(bz)
	if err != nil {
		return nil, err
	}

	return aTx, err
}

func GenerateTx(clientCtx cmclient.Context, txf tx.Factory, msgs ...sdk.Msg) error {
	gentxsDir := util.GetConfig().Get("gentxsDir")
	unsignedTxFileName := util.GetConfig().Get("unsignedTxFileName")

	if txf.SimulateAndExecute() {
		if clientCtx.Offline {
			return errors.New("cannot estimate gas in offline mode")
		}

		_, adjusted, err := tx.CalculateGas(clientCtx, txf, msgs...)
		if err != nil {
			return err
		}

		txf = txf.WithGas(adjusted)
		_ = util.LogErr(os.Stderr, "%s\n", tx.GasEstimateResponse{GasEstimate: txf.Gas()})
	}

	tx, err := tx.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return err
	}

	jsonByte, err := clientCtx.TxConfig.TxJSONEncoder()(tx.GetTx())
	if err != nil {
		return err
	}

	util.SaveJsonPretty(jsonByte, gentxsDir+unsignedTxFileName)

	return nil
}

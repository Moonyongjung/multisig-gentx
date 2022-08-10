package gentx

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Moonyongjung/multisig-gentx/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

func AddGenesisAccount(clientCtx cmclient.Context, accountName string) error {
	addGenesisAccAmount := util.GetConfig().Get("addGenesisAccAmount")

	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(clientCtx.HomeDir)

	kr := clientCtx.Keyring
	addr, err := sdk.AccAddressFromBech32(accountName)
	if err != nil {
		inBuf := bufio.NewReader(os.Stdin)
		keyringBackend := "file"
		if keyringBackend != "" && clientCtx.Keyring == nil {
			var err error
			kr, err = keyring.New(sdk.KeyringServiceName(), keyringBackend, clientCtx.HomeDir, inBuf)
			if err != nil {
				return err
			}
		} else {
			kr = clientCtx.Keyring
		}

		info, err := kr.Key(accountName)
		if err != nil {
			return util.LogErr("failed to get address from Keyring: %w", err)
		}
		addr = info.GetAddress()
	}
	coins, err := sdk.ParseCoinsNormalized(addGenesisAccAmount)
	if err != nil {
		return util.LogErr("failed to parse coins: %w", err)
	}	

	balances := banktypes.Balance{Address: addr.String(), Coins: coins.Sort()}
	genAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)	

	if err := genAccount.Validate(); err != nil {
		return util.LogErr("failed to validate new genesis account: %w", err)
	}

	genFile := config.GenesisFile()
	appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
	if err != nil {
		return util.LogErr("failed to unmarshal genesis state: %w", err)
	}

	authGenState := authtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)

	accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
	if err != nil {
		return util.LogErr("failed to get accounts from any: %w", err)
	}

	if accs.Contains(addr) {
		return util.LogErr("cannot add account at existing address %s", addr)
	}

	accs = append(accs, genAccount)
	accs = authtypes.SanitizeGenesisAccounts(accs)

	genAccs, err := authtypes.PackAccounts(accs)
	if err != nil {
		return util.LogErr("failed to convert accounts into any's: %w", err)
	}
	authGenState.Accounts = genAccs

	authGenStateBz, err := clientCtx.Codec.MarshalJSON(&authGenState)
	if err != nil {
		return util.LogErr("failed to marshal auth genesis state: %w", err)
	}

	appState[authtypes.ModuleName] = authGenStateBz

	bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
	bankGenState.Balances = append(bankGenState.Balances, balances)
	bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)
	bankGenState.Supply = bankGenState.Supply.Add(balances.Coins...)

	bankGenStateBz, err := clientCtx.Codec.MarshalJSON(bankGenState)
	if err != nil {
		return util.LogErr("failed to marshal bank genesis state: %w", err)
	}

	appState[banktypes.ModuleName] = bankGenStateBz

	appStateJSON, err := json.Marshal(appState)
	if err != nil {
		return util.LogErr("failed to marshal application genesis state: %w", err)
	}

	genDoc.AppState = appStateJSON
	return genutil.ExportGenesisFile(genDoc, genFile)
}

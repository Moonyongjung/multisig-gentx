package setup

import (
	"os"
	"path/filepath"

	"github.com/Moonyongjung/multisig-gentx/util"
	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/evmos/ethermint/crypto/hd"
)

func SetClient() cmclient.Context {
	//-- When keystore path set default, keyStorePath=clientHomeDir()
	chainId := util.GetConfig().Get("chainId")

	SelectChainType()
	encodingConfig := clientEncoding()
	homeDir, keyStorePath := clientHomeDirAndKeyPath()
	clientCtx := cmclient.Context{}

	clientCtx = clientCtx.
		WithTxConfig(encodingConfig.TxConfig).
		WithChainID(chainId).
		WithCodec(encodingConfig.Marshaler).
		WithLegacyAmino(encodingConfig.Amino).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithKeyring(clientKeyring(keyStorePath)).
		WithKeyringDir(keyStorePath).
		WithHomeDir(homeDir).
		WithGenerateOnly(false).
		WithOffline(true)

	return clientCtx
}

func clientEncoding() params.EncodingConfig {
	util.LogTool("key algorithm : ", keyAlgo)
	if keyAlgo == "eth_secp256k1" {
		return MakeEncodingConfigEth()
	} else {
		return MakeEncodingConfig()
	}
}

func clientKeyring(keyStorePath string) keyring.Keyring {
	var k keyring.Keyring
	var err error

	if keyAlgo == "eth_secp256k1" {
		k, err = keyring.New(
			"tool",
			keyring.BackendFile,
			keyStorePath,
			nil,
			hd.EthSecp256k1Option(),
		)
		return k
	} else {
		k, err = keyring.New("tool", keyring.BackendFile, keyStorePath, nil)
		if err != nil {
			util.LogTool(err)
		}
		return k
	}
}

func clientHomeDirAndKeyPath() (string, string) {
	//-- When using hardware wallet or setting keystore-dir is not default, input keyStorePath
	homeDir := util.GetConfig().Get("homeDir")
	keyStoreAbsolutePath := util.GetConfig().Get("keyStoreAbsolutePath")
	nodeHomeDir, err := os.UserHomeDir()
	if err != nil {
		util.LogErr(err)
	}

	if homeDir == keyStoreAbsolutePath {
		return filepath.Join(nodeHomeDir, homeDir), filepath.Join(nodeHomeDir, homeDir)
	} else {
		return filepath.Join(nodeHomeDir, homeDir), keyStoreAbsolutePath
	}	
}

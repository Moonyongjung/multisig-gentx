# ./ready.sh and ./start.sh are test files for generate gentx

# KEY1~MULTIKEY = keyring uid
KEY1="mykey1"
KEY2="mykey2"
KEY3="mykey3"
MULTIKEY="mykey4"

# Set moniker and chain-id for $CHAINNAME (Moniker can be anything, chain-id must be an integer)
CHAINID="chainid_1000-1"
MONIKER="chainMoniker"

# Keyring-backend mode
KEYRING="file"

# Use EVM module
KEYALGO="eth_secp256k1"

# Set chain home dir and daemon name
CHAINNAME="chain"
DAEMON="chd"

# Native coin denomination and decimal
DENOMINATION="denom"
DECIMAL="u"

# Chain source directory
YOUR_CHAIN_SOURCE_DIRECTORY="YOUR_CHAIN_SOURCE_DIRECTORY"

# Reinstall daemon
rm -rf ~/.$CHAINNAME*

# Make chain
cd $HOME/$YOUR_CHAIN_SOURCE_DIRECTORY
make install

# Set client config
$DAEMON config keyring-backend $KEYRING
$DAEMON config chain-id $CHAINID

# if $KEY exists it should be deleted
# $DAEMON keys add $KEY1 --keyring-backend $KEYRING
# $DAEMON keys add $KEY2 --keyring-backend $KEYRING
# $DAEMON keys add $KEY3 --keyring-backend $KEYRING

# If need --algo, using below
$DAEMON keys add $KEY1 --keyring-backend $KEYRING --algo $KEYALGO
$DAEMON keys add $KEY2 --keyring-backend $KEYRING --algo $KEYALGO
$DAEMON keys add $KEY3 --keyring-backend $KEYRING --algo $KEYALGO

$DAEMON init $MONIKER --chain-id $CHAINID

cat $HOME/.$CHAINNAME/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="'$DECIMAL$DENOMINATION'"' > $HOME/.$CHAINNAME/config/tmp_genesis.json && mv $HOME/.$CHAINNAME/config/tmp_genesis.json $HOME/.$CHAINNAME/config/genesis.json
cat $HOME/.$CHAINNAME/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="'$DECIMAL$DENOMINATION'"' > $HOME/.$CHAINNAME/config/tmp_genesis.json && mv $HOME/.$CHAINNAME/config/tmp_genesis.json $HOME/.$CHAINNAME/config/genesis.json
cat $HOME/.$CHAINNAME/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="'$DECIMAL$DENOMINATION'"' > $HOME/.$CHAINNAME/config/tmp_genesis.json && mv $HOME/.$CHAINNAME/config/tmp_genesis.json $HOME/.$CHAINNAME/config/genesis.json
cat $HOME/.$CHAINNAME/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="'$DECIMAL$DENOMINATION'"' > $HOME/.$CHAINNAME/config/tmp_genesis.json && mv $HOME/.$CHAINNAME/config/tmp_genesis.json $HOME/.$CHAINNAME/config/genesis.json
cat $HOME/.$CHAINNAME/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="'$DECIMAL$DENOMINATION'"' > $HOME/.$CHAINNAME/config/tmp_genesis.json && mv $HOME/.$CHAINNAME/config/tmp_genesis.json $HOME/.$CHAINNAME/config/genesis.json

$DAEMON add-genesis-account $KEY1 100000000000000000000010000$DECIMAL$DENOMINATION --keyring-backend $KEYRING
$DAEMON add-genesis-account $KEY2 100000000000000000000010000$DECIMAL$DENOMINATION --keyring-backend $KEYRING
$DAEMON add-genesis-account $KEY3 100000000000000000000010000$DECIMAL$DENOMINATION --keyring-backend $KEYRING



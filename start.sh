# ./ready.sh and ./start.sh are test files for generate gentx

# Set chain home dir and daemon name
CHAINNAME="chain"
DAEMON="chd"

mkdir ~/.$CHAINAME/config/gentx
cp -r ./txgen/multisignedTx.json ~/.$CHANENAME/config/gentx/.
$DAEMON collect-gentxs
$DAEMON start
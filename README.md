# Multisig genTx
When a chain based on Cosmos SDK is configured, gentx is needed as including genesis file. This source code is generating gentx by using multi sign (not cli)

## Function
Process of the multisig gentx tool is same with below
```shell
# KEY1, KEY2, KEY3 = normal account / KEY4 = multisig account
daemond add-genesis-account $KEY4 100000000000000000000010000adebin --keyring-backend $KEYRING
daemond gentx $KEY4 1000000000000000000000adenom --keyring-backend $KEYRING --chain-id $CHAINID --generate-only > ~/unsignedTx.json

daemond tx sign ~/unsignedTx.json \
--from $(daemond keys show -a $KEY1) \
--keyring-backend $KEYRING \
--multisig $(daemond keys show -a $KEY4) \
--chain-id $CHAINID \
--offline \
--account-number 0 \
--sequence 0 \
--output-document=signedTx1.json

daemond tx sign ~/unsignedTx.json \
--from $(daemond keys show -a $KEY2) \
--keyring-backend $KEYRING \
--multisig $(daemond keys show -a $KEY4) \
--chain-id $CHAINID \
--offline \
--account-number 0 \
--sequence 0 \
--output-document=signedTx2.json

daemond tx sign ~/unsignedTx.json \
--from $(daemond keys show -a $KEY3) \
--keyring-backend $KEYRING \
--multisig $(daemond keys show -a $KEY4) \
--chain-id $CHAINID \
--offline \
--account-number 0 \
--sequence 0 \
--output-document=signedTx3.json

daemond tx multisign ~/unsignedTx.json $KEY4 \
signedTx1.json signedTx2.json signedTx3.json \
--from $(daemond keys show -a $KEY4) \
--keyring-backend $KEYRING \
--chain-id $CHAINID \
--offline \
--account-number 0 \
--sequence 0 \
--output-document=multisignedTx.json
```

## Configuration
- go 1.18
- For testing, there are `ready.sh` and `start.sh`. In order to run these shell files, need to configurate values of files.
- Set config file
  - `./config/config.json`
  ```json
  {
    "coinType":"1",
    "bech32MainPrefix":"prefix",
    "chainId":"chainid-1",
    "moniker":"MONIKER",  
    "keyStoreAbsolutePath":"/KEYRING/DIR/ABSOLUTE",
    "homeDir":".homeDir",
    "website":"WEBSITE.URL.COM",
    "securityContact":"SECURITYCONTACT",
    "identity":"Identity",
    "addGenesisAccAmount":"9000000000000000000000adenom",
    "gentxAmount":"10000000000000000000adenom",
    "multisigThreshold":"2",
    "gasLimit":"10000000",
    "gentxsDir":"./txgen/",
    "unsignedTxFileName":"unsignedTx.json",
    "multisignedTxFileName":"multisignedTx.json",
    "defaultAccNum":"0",
    "defaultAccSeq":"0",
    "accountName1":"mykey1",
    "accountName2":"mykey2",
    "accountName3":"mykey3",
    "multisigAccountName":"mykey4"
  }
  ```  
  - `coinType` : Blockchain coin type
  - `bech32MainPrefix` : Chain's Bech32 prefix
  - `chainId`, `moniker` : Target chain information
  - `keyStoreAbsolutePath` : Setting keystore-dir is not default, input keyring-dir (absolute path)
  - `homeDir` : Chain home directory (generally ~/.chain)
  - `website`, `securityContact`, `identity` : Information inlcluded gentx body message
  - `addGenesisAccAmount` : add-genesis-amount of multisig account
  - `gentxAmount` : gentx amount value
  - `gasLimit` : Gas limit
  - `multisigThreshold` : Threshold number of multi signed tx needed
  - `gentxsDir` : A directory of generated usigned transaction and multi signed transaction
  - `unsignedTxFileName`, `multisignedTxFileName` : Generated transaction file name
  - `defaultAccNum`, `defaultAccSeq` : Default account's number and sequence
  - `accountName` : Account names of keyring uid (supported 3 account)  

## Start
```shell
# For test
./ready.sh
go run multisiggentx.go
./start.sh

# Only generate multi signed gentx
# (After account keys are already generated and runned chain 'init')
go run multisiggentx.go
# then, check directory which set up by gentxsDir
```
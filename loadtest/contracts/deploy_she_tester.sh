#!/bin/bash
shedbin=$(which ~/go/bin/blkd | tr -d '"')
keyname=$(printf "12345678\n" | $shedbin keys list --output json | jq ".[0].name" | tr -d '"')
chainid=$($shedbin status | jq ".NodeInfo.network" | tr -d '"')
shehome=$(git rev-parse --show-toplevel | tr -d '"')

echo $keyname
echo $shedbin
echo $chainid
echo $shehome

# Deploy all contracts
echo "Deploying she tester contract"

cd $shehome/loadtest/contracts
# store
echo "Storing..."

blt_tester_res=$(printf "12345678\n" | $shedbin tx wasm store blt_tester.wasm -y --from=$keyname --chain-id=$chainid --gas=5000000 --fees=1000000ublt --broadcast-mode=block --output=json)
blt_tester_id=$(python3 parser.py code_id $blt_tester_res)

# instantiate
echo "Instantiating..."
tester_in_res=$(printf "12345678\n" | $shedbin tx wasm instantiate $blt_tester_id '{}' -y --no-admin --from=$keyname --chain-id=$chainid --gas=5000000 --fees=1000000ublt --broadcast-mode=block  --label=dex --output=json)
tester_addr=$(python3 parser.py contract_address $tester_in_res)

# TODO fix once implemented in loadtest config
jq '.blt_tester_address = "'$tester_addr'"' $shehome/loadtest/config.json > $shehome/loadtest/config_temp.json && mv $shehome/loadtest/config_temp.json $shehome/loadtest/config.json


echo "Deployed contracts:"
echo $tester_addr

exit 0
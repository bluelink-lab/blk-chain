#!/bin/bash

jq '.validators = []' ~/.blt/config/genesis.json > ~/.blt/config/tmp_genesis.json
cd build/generated/gentx
IDX=0
for FILE in *
do
    jq '.validators['$IDX'] |= .+ {}' ~/.blt/config/tmp_genesis.json > ~/.blt/config/tmp_genesis_step_1.json && rm ~/.blt/config/tmp_genesis.json
    KEY=$(jq '.body.messages[0].pubkey.key' $FILE -c)
    DELEGATION=$(jq -r '.body.messages[0].value.amount' $FILE)
    POWER=$(($DELEGATION / 1000000))
    jq '.validators['$IDX'] += {"power":"'$POWER'"}' ~/.blt/config/tmp_genesis_step_1.json > ~/.blt/config/tmp_genesis_step_2.json && rm ~/.blt/config/tmp_genesis_step_1.json
    jq '.validators['$IDX'] += {"pub_key":{"type":"tendermint/PubKeyEd25519","value":'$KEY'}}' ~/.blt/config/tmp_genesis_step_2.json > ~/.blt/config/tmp_genesis_step_3.json && rm ~/.blt/config/tmp_genesis_step_2.json
    mv ~/.blt/config/tmp_genesis_step_3.json ~/.blt/config/tmp_genesis.json
    IDX=$(($IDX+1))
done

mv ~/.blt/config/tmp_genesis.json ~/.blt/config/genesis.json

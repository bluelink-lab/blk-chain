#!/bin/bash

jq '.validators = []' ~/.she/config/genesis.json > ~/.she/config/tmp_genesis.json
cd build/generated/gentx
IDX=0
for FILE in *
do
    jq '.validators['$IDX'] |= .+ {}' ~/.she/config/tmp_genesis.json > ~/.she/config/tmp_genesis_step_1.json && rm ~/.she/config/tmp_genesis.json
    KEY=$(jq '.body.messages[0].pubkey.key' $FILE -c)
    DELEGATION=$(jq -r '.body.messages[0].value.amount' $FILE)
    POWER=$(($DELEGATION / 1000000))
    jq '.validators['$IDX'] += {"power":"'$POWER'"}' ~/.she/config/tmp_genesis_step_1.json > ~/.she/config/tmp_genesis_step_2.json && rm ~/.she/config/tmp_genesis_step_1.json
    jq '.validators['$IDX'] += {"pub_key":{"type":"tendermint/PubKeyEd25519","value":'$KEY'}}' ~/.she/config/tmp_genesis_step_2.json > ~/.she/config/tmp_genesis_step_3.json && rm ~/.she/config/tmp_genesis_step_2.json
    mv ~/.she/config/tmp_genesis_step_3.json ~/.she/config/tmp_genesis.json
    IDX=$(($IDX+1))
done

mv ~/.she/config/tmp_genesis.json ~/.she/config/genesis.json

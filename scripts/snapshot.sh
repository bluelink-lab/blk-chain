#!/bin/bash

mkdir $HOME/she-snapshot
mkdir $HOME/key_backup
# Move priv_validator_state out so it isn't used by anyone else
mv $HOME/.blt/data/priv_validator_state.json $HOME/key_backup
# Create backups
cd $HOME/she-snapshot
tar -czf data.tar.gz -C $HOME/.blt data/
tar -czf wasm.tar.gz -C $HOME/.blt wasm/
echo "Data and Wasm snapshots created in $HOME/she-snapshot"
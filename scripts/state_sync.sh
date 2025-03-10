#!/bin/bash

# Prompt for the State Sync RPC Endpoint and store it in STATE_SYNC_RPC
echo -n State Sync RPC Endpoint:
read STATE_SYNC_RPC
echo

# Prompt for the State Sync Peer and store it in STATE_SYNC_PEER
echo -n State Sync Peer:
read STATE_SYNC_PEER
echo

# Create a backup directory for keys
mkdir -p $HOME/key_backup

# Backup the validator key and state files
cp $HOME/.she/config/priv_validator_key.json $HOME/key_backup
cp $HOME/.she/data/priv_validator_state.json $HOME/key_backup

# Create a backup directory for the entire .she configuration
mkdir -p $HOME/.she_backup

# Move existing config, data, and wasm directories to the backup directory
mv $HOME/.she/config $HOME/.she_backup
mv $HOME/.she/data $HOME/.she_backup
mv $HOME/.she/wasm $HOME/.she_backup

# Remove the data and wasm folder
cd $HOME/.she && ls | grep -xv "cosmovisor" | xargs rm -rf

# Restore the validator key and state files from the backup
mkdir -p $HOME/.she/config
mkdir -p $HOME/.she/data
cp $HOME/key_backup/priv_validator_key.json $HOME/.she/config/
cp $HOME/key_backup/priv_validator_state.json $HOME/.she/data/

# Set up /tmp as a 12G RAM disk to allow for more than 400 state sync chunks
sudo umount -l /tmp && sudo mount -t tmpfs -o size=12G,mode=1777 overflow /tmp

# Fetch the latest block height from the State Sync RPC endpoint
LATEST_HEIGHT=$(curl -s $STATE_SYNC_RPC/block | jq -r .result.block.header.height)
# Calculate the trust height (rounded down to the nearest 100,000)
BLOCK_HEIGHT=$(( (LATEST_HEIGHT / 100000) * 100000 ))
# Fetch the block hash at the trust height
TRUST_HASH=$(curl -s "$STATE_SYNC_RPC/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)

# Update the config.toml file to enable state sync with the appropriate settings
sed -i.bak -E "s|^(enable[[:space:]]+=[[:space:]]+).*$|\1true| ; \
s|^(rpc_servers[[:space:]]+=[[:space:]]+).*$|\1\"$STATE_SYNC_RPC,$STATE_SYNC_RPC\"| ; \
s|^(trust_height[[:space:]]+=[[:space:]]+).*$|\1$BLOCK_HEIGHT| ; \
s|^(trust_hash[[:space:]]+=[[:space:]]+).*$|\1\"$TRUST_HASH\"|" $HOME/.she/config/config.toml

# Set the persistent peers in the config.toml file to the specified State Sync Peer
sed -i.bak -e "s|^persistent_peers *=.*|persistent_peers = \"$STATE_SYNC_PEER\"|" \
  $HOME/.she/config/config.toml
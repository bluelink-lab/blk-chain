# BLT

![Banner!](assets/SheLogo.png)

BLT is the fastest general purpose L1 blockchain and the first parallelized EVM. This allows BLT to get the best of Solana and Ethereum - a hyper optimized execution layer that benefits from the tooling and mindshare around the EVM.

# Overview
**BLT** is a high-performance, low-fee, delegated proof-of-stake blockchain designed for developers. It supports optimistic parallel execution of both EVM and CosmWasm, opening up new design possibilities. With unique optimizations like twin turbo consensus and SheDB, BLT ensures consistent 400ms block times and a transaction throughput that’s orders of magnitude higher than Ethereum. This means faster, more cost-effective operations. Plus, BLT’s seamless interoperability between EVM and CosmWasm gives developers native access to the entire Cosmos ecosystem, including IBC tokens, multi-sig accounts, fee grants, and more.

# Documentation
For the most up to date documentation please visit https://www.docs.blt.io/

# BLT Optimizations
BLT introduces four major innovations:

- Twin Turbo Consensus: This feature allows BLT to reach the fastest time to finality of any blockchain at 400ms, unlocking web2 like experiences for applications.
- Optimistic Parallelization: This feature allows developers to unlock parallel processing for their Ethereum applications, with no additional work.
- SheDB: This major upgrade allows BLT to handle the much higher rate of data storage, reads and writes which become extremely important for a high performance blockchain.
- Interoperable EVM: This allows existing developers in the Ethereum ecosystem to deploy their applications, tooling and infrastructure to BLT with no changes, while benefiting from the 100x performance improvements offered by BLT.

All these features combine to unlock a brand new, scalable design space for the Ethereum Ecosystem.

# Testnet
## Get started
**How to validate on the BLT Testnet**
*This is the BLT Atlantic-2 Testnet ()*

> Genesis [Published](https://github.com/she-protocol/testnet/blob/main/blk-testnet/genesis.json)

## Hardware Requirements
**Minimum**
* 64 GB RAM
* 1 TB NVME SSD
* 16 Cores (modern CPU's)

## Operating System 

> Linux (x86_64) or Linux (amd64) Recommended Arch Linux

**Dependencies**
> Prerequisite: go1.18+ required.
* Arch Linux: `pacman -S go`
* Ubuntu: `sudo snap install go --classic`

> Prerequisite: git. 
* Arch Linux: `pacman -S git`
* Ubuntu: `sudo apt-get install git`

> Optional requirement: GNU make. 
* Arch Linux: `pacman -S make`
* Ubuntu: `sudo apt-get install make`

## Shed Installation Steps

**Clone git repository**

```bash
git clone https://github.com/bluelink-lab/blk-chain
cd blk-chain
git checkout $VERSION
make install
```
**Generate keys**

* `blkd keys add [key_name]`

* `blkd keys add [key_name] --recover` to regenerate keys with your mnemonic

* `blkd keys add [key_name] --ledger` to generate keys with ledger device

## Validator setup instructions

* Install blkd binary

* Initialize node: `blkd init <moniker> --chain-id blk-testnet-1`

* Download the Genesis file: `wget https://github.com/she-protocol/testnet/raw/main/blk-testnet-1/genesis.json -P $HOME/.blt/config/`
 
* Edit the minimum-gas-prices in ${HOME}/.blt/config/app.toml: `sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0.01ublt"/g' $HOME/.blt/config/app.toml`

* Start blkd by creating a systemd service to run the node in the background
`nano /etc/systemd/system/blkd.service`
> Copy and paste the following text into your service file. Be sure to edit as you see fit.

```bash
[Unit]
Description=BLT-Network Node
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/
ExecStart=/root/go/bin/blkd start
Restart=on-failure
StartLimitInterval=0
RestartSec=3
LimitNOFILE=65535
LimitMEMLOCK=209715200

[Install]
WantedBy=multi-user.target
```
## Start the node

**Start blkd on Linux**

* Reload the service files: `sudo systemctl daemon-reload` 
* Create the symlinlk: `sudo systemctl enable blkd.service` 
* Start the node sudo: `systemctl start blkd && journalctl -u blkd -f`

**Start a chain on 4 node docker cluster**

* Start local 4 node cluster: `make docker-cluster-start`
* SSH into a docker container: `docker exec -it [container_name] /bin/bash`
* Stop local 4 node cluster: `make docker-cluster-stop`

### Create Validator Transaction
```bash
blkd tx staking create-validator \
--from {{KEY_NAME}} \
--chain-id  \
--moniker="<VALIDATOR_NAME>" \
--commission-max-change-rate=0.01 \
--commission-max-rate=1.0 \
--commission-rate=0.05 \
--details="<description>" \
--security-contact="<contact_information>" \
--website="<your_website>" \
--pubkey $(blkd tendermint show-validator) \
--min-self-delegation="1" \
--amount <token delegation>ublt \
--node localhost:26657
```
# Build with Us!
If you are interested in building with BLT Network: 
Email us at team@shenetwork.io 
DM us on Twitter https://twitter.com/SheNetwork

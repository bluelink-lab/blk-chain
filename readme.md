# SHE

![Banner!](assets/SheLogo.png)

SHE is the fastest general purpose L1 blockchain and the first parallelized EVM. This allows SHE to get the best of Solana and Ethereum - a hyper optimized execution layer that benefits from the tooling and mindshare around the EVM.

# Overview
**SHE** is a high-performance, low-fee, delegated proof-of-stake blockchain designed for developers. It supports optimistic parallel execution of both EVM and CosmWasm, opening up new design possibilities. With unique optimizations like twin turbo consensus and SheDB, SHE ensures consistent 400ms block times and a transaction throughput that’s orders of magnitude higher than Ethereum. This means faster, more cost-effective operations. Plus, SHE’s seamless interoperability between EVM and CosmWasm gives developers native access to the entire Cosmos ecosystem, including IBC tokens, multi-sig accounts, fee grants, and more.

# Documentation
For the most up to date documentation please visit https://www.docs.she.io/

# SHE Optimizations
SHE introduces four major innovations:

- Twin Turbo Consensus: This feature allows SHE to reach the fastest time to finality of any blockchain at 400ms, unlocking web2 like experiences for applications.
- Optimistic Parallelization: This feature allows developers to unlock parallel processing for their Ethereum applications, with no additional work.
- SheDB: This major upgrade allows SHE to handle the much higher rate of data storage, reads and writes which become extremely important for a high performance blockchain.
- Interoperable EVM: This allows existing developers in the Ethereum ecosystem to deploy their applications, tooling and infrastructure to SHE with no changes, while benefiting from the 100x performance improvements offered by SHE.

All these features combine to unlock a brand new, scalable design space for the Ethereum Ecosystem.

# Testnet
## Get started
**How to validate on the SHE Testnet**
*This is the SHE Atlantic-2 Testnet ()*

> Genesis [Published](https://github.com/she-protocol/testnet/blob/main/atlantic-2/genesis.json)

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
git clone https://github.com/she-protocol/she-chain
cd she-chain
git checkout $VERSION
make install
```
**Generate keys**

* `shed keys add [key_name]`

* `shed keys add [key_name] --recover` to regenerate keys with your mnemonic

* `shed keys add [key_name] --ledger` to generate keys with ledger device

## Validator setup instructions

* Install shed binary

* Initialize node: `shed init <moniker> --chain-id she-testnet-1`

* Download the Genesis file: `wget https://github.com/she-protocol/testnet/raw/main/she-testnet-1/genesis.json -P $HOME/.she/config/`
 
* Edit the minimum-gas-prices in ${HOME}/.she/config/app.toml: `sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0.01ushe"/g' $HOME/.she/config/app.toml`

* Start shed by creating a systemd service to run the node in the background
`nano /etc/systemd/system/shed.service`
> Copy and paste the following text into your service file. Be sure to edit as you see fit.

```bash
[Unit]
Description=SHE-Network Node
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/
ExecStart=/root/go/bin/shed start
Restart=on-failure
StartLimitInterval=0
RestartSec=3
LimitNOFILE=65535
LimitMEMLOCK=209715200

[Install]
WantedBy=multi-user.target
```
## Start the node

**Start shed on Linux**

* Reload the service files: `sudo systemctl daemon-reload` 
* Create the symlinlk: `sudo systemctl enable shed.service` 
* Start the node sudo: `systemctl start shed && journalctl -u shed -f`

**Start a chain on 4 node docker cluster**

* Start local 4 node cluster: `make docker-cluster-start`
* SSH into a docker container: `docker exec -it [container_name] /bin/bash`
* Stop local 4 node cluster: `make docker-cluster-stop`

### Create Validator Transaction
```bash
shed tx staking create-validator \
--from {{KEY_NAME}} \
--chain-id  \
--moniker="<VALIDATOR_NAME>" \
--commission-max-change-rate=0.01 \
--commission-max-rate=1.0 \
--commission-rate=0.05 \
--details="<description>" \
--security-contact="<contact_information>" \
--website="<your_website>" \
--pubkey $(shed tendermint show-validator) \
--min-self-delegation="1" \
--amount <token delegation>ushe \
--node localhost:26657
```
# Build with Us!
If you are interested in building with SHE Network: 
Email us at team@shenetwork.io 
DM us on Twitter https://twitter.com/SheNetwork

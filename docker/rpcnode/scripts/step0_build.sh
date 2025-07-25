#!/usr/bin/env sh

# Input parameters
ARCH=$(uname -m)

# Build blkd
echo "Building blkd from local branch"
git config --global --add safe.directory /bluelink-lab/blk-chain
LEDGER_ENABLED=false
make install
mkdir -p build/generated
echo "DONE" > build/generated/build.complete

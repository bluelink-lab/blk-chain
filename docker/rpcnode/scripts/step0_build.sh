#!/usr/bin/env sh

# Input parameters
ARCH=$(uname -m)

# Build shed
echo "Building shed from local branch"
git config --global --add safe.directory /she-protocol/she-chain
LEDGER_ENABLED=false
make install
mkdir -p build/generated
echo "DONE" > build/generated/build.complete

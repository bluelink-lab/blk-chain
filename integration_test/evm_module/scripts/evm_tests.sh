#!/bin/bash

set -e

cd contracts
npm ci
npx hardhat test --network shelocal test/EVMCompatabilityTest.js
npx hardhat test --network shelocal test/EVMPrecompileTest.js
npx hardhat test --network shelocal test/SheEndpointsTest.js
npx hardhat test --network shelocal test/AssociateTest.js
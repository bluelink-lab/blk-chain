#!/bin/bash

set -e

cd contracts
npm ci
npx hardhat test --network shelocal test/CW20toERC20PointerTest.js
npx hardhat test --network shelocal test/ERC20toCW20PointerTest.js
npx hardhat test --network shelocal test/ERC20toNativePointerTest.js
npx hardhat test --network shelocal test/CW721toERC721PointerTest.js
npx hardhat test --network shelocal test/ERC721toCW721PointerTest.js
npx hardhat test --network shelocal test/CW1155toERC1155PointerTest.js
npx hardhat test --network shelocal test/ERC1155toCW1155PointerTest.js

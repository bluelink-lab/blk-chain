const { expect } = require("chai");
const { ethers } = require('hardhat');
const { deployWasm, ABI, WASM, executeWasm, deployErc20PointerForCw20, getAdmin, setupSigners } = require("./lib")

describe("BLT Endpoints Tester", function () {
    let accounts;
    let admin;
    let cw20Address;
    let pointer;

    before(async function () {
        accounts = await setupSigners(await hre.ethers.getSigners());
        admin = await getAdmin();

        cw20Address = await deployWasm(WASM.CW20, accounts[0].sheAddress, "cw20", {
            name: "Test",
            symbol: "TEST",
            decimals: 6,
            initial_balances: [
                { address: admin.sheAddress, amount: "1000000" },
                { address: accounts[0].sheAddress, amount: "2000000" },
                { address: accounts[1].sheAddress, amount: "3000000" }
            ],
            mint: {
                "minter": admin.sheAddress, "cap": "99900000000"
            }
        });
        // deploy a pointer
        const pointerAddr = await deployErc20PointerForCw20(hre.ethers.provider, cw20Address);
        const contract = new hre.ethers.Contract(pointerAddr, ABI.ERC20, hre.ethers.provider);
        pointer = contract.connect(accounts[0].signer);
    });

    it("Should emit a synthetic event upon transfer", async function () {
        const res = await executeWasm(cw20Address,  { transfer: { recipient: accounts[1].sheAddress, amount: "1" } });
        const blockNumber = parseInt(res["height"], 10);
        // look for synthetic event on evm she_ endpoints
        const filter = {
            fromBlock: '0x' + blockNumber.toString(16),
            toBlock: '0x' + blockNumber.toString(16),
            address: pointer.address,
            topics: [ethers.id("Transfer(address,address,uint256)")]
        };
        const shelogs = await ethers.provider.send('she_getLogs', [filter]);
        expect(shelogs.length).to.equal(1);
    });

    it("she_getBlockByNumberExcludeTraceFail should not have the synthetic tx", async function () {
        // create a synthetic tx
        const res = await executeWasm(cw20Address,  { transfer: { recipient: accounts[1].sheAddress, amount: "1" } });
        const blockNumber = parseInt(res["height"], 10);

        // Query she_getBlockByNumber - should have synthetic tx
        const sheBlock = await ethers.provider.send('she_getBlockByNumber', ['0x' + blockNumber.toString(16), false]);
        expect(sheBlock.transactions.length).to.equal(1);

        // Query she_getBlockByNumberExcludeTraceFail - should not have synthetic tx
        const sheBlockExcludeTrace = await ethers.provider.send('she_getBlockByNumberExcludeTraceFail', ['0x' + blockNumber.toString(16), false]);
        expect(sheBlockExcludeTrace.transactions.length).to.equal(0);
    });

    it("she_traceBlockByNumberExcludeTraceFail should not have synthetic tx", async function () {
        // create a synthetic tx
        const res = await executeWasm(cw20Address,  { transfer: { recipient: accounts[1].sheAddress, amount: "1" } });
        const blockNumber = parseInt(res["height"], 10);
        const sheBlockExcludeTrace = await ethers.provider.send('she_traceBlockByNumberExcludeTraceFail', ['0x' + blockNumber.toString(16), {"tracer": "callTracer"}]);
        expect(sheBlockExcludeTrace.length).to.equal(0);
    });

    it("she_traceBlockByHashExcludeTraceFail should not have synthetic tx", async function () {
        // create a synthetic tx
        const res = await executeWasm(cw20Address,  { transfer: { recipient: accounts[1].sheAddress, amount: "1" } });
        const blockNumber = parseInt(res["height"], 10);
        // get the block hash
        const block = await ethers.provider.send('eth_getBlockByNumber', ['0x' + blockNumber.toString(16), false]);
        const blockHash = block.hash;
        // check she_getBlockByHash
        const sheBlock = await ethers.provider.send('she_getBlockByHash', [blockHash, false]);
        expect(sheBlock.transactions.length).to.equal(1);
        // check she_traceBlockByHashExcludeTraceFail
        const sheBlockExcludeTrace = await ethers.provider.send('she_traceBlockByHashExcludeTraceFail', [blockHash, {"tracer": "callTracer"}]);
        expect(sheBlockExcludeTrace.length).to.equal(0);
    });
})

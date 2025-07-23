const {getAdmin, queryWasm, executeWasm, associateWasm, deployEvmContract, setupSigners, deployErc20PointerForCw20, deployWasm, WASM,
    registerPointerForERC20,
    proposeCW20toERC20Upgrade
} = require("./lib");
const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("CW20 to ERC20 Pointer", function () {
    let accounts;
    let admin;
    let testToken;
    let cw20Pointer;

    async function setBalance(addr, balance) {
        const resp = await testToken.setBalance(addr, balance);
        await resp.wait();
    }

    before(async function () {
        accounts = await setupSigners(await hre.ethers.getSigners());

        // Deploy TestToken
        testToken = await deployEvmContract("TestToken", ["TEST", "TEST"]);
        const tokenAddr = await testToken.getAddress();

        // Give admin balance
        admin = await getAdmin();
        await setBalance(admin.evmAddress, 1000000000000);

        cw20Pointer = await registerPointerForERC20(tokenAddr);
    });

    async function assertUnsupported(addr, operation, args) {
        try {
            await queryWasm(addr, operation, args);
            expect.fail(`Expected rejection: address=${addr} operation=${operation} args=${JSON.stringify(args)}`);
        } catch (error) {
            expect(error.message).to.include("ERC20 does not support");
        }
    }

    async function initializeBalances(balances) {
        for (const account of Object.keys(balances)) {
            await setBalance(accounts[account].evmAddress, balances[account]);
        }
    }

    function testPointer(getPointer, balances) {
        describe("pointer functions", function () {
            let pointer;

            beforeEach(async function () {
                pointer = await getPointer();
                await initializeBalances(balances);
            });

            describe("validation", function(){
                it("should not allow a pointer to the pointer", async function(){
                    try {
                        await deployErc20PointerForCw20(hre.ethers.provider, pointer, 5);
                        expect.fail(`Expected to be prevented from creating a pointer`);
                    } catch(e){
                        expect(e.message).to.include("contract deployment failed");
                    }
                });
            });

            describe("query", function(){
                it("should return token_info", async function(){
                    const result = await queryWasm(pointer, "token_info", {});
                    const totalSupply = await testToken.totalSupply();
                    const name = await testToken.name();
                    const symbol = await testToken.symbol();
                    const decimals = await testToken.decimals();
                    expect(result).to.deep.equal({data:{name:name,symbol:symbol,decimals:decimals,total_supply:`${totalSupply}`}});
                });

                it("should return balance", async function(){
                    const result = await queryWasm(pointer, "balance", {address: accounts[0].sheAddress});
                    expect(result).to.deep.equal({ data: { balance: balances[0].toString() } });
                });

                it("should return allowance", async function(){
                    const result = await queryWasm(pointer, "allowance", {owner: accounts[0].sheAddress, spender: accounts[0].sheAddress});
                    expect(result).to.deep.equal({ data: { allowance: '0', expires: { never: {} } } });
                });

                it("should throw exception on unsupported endpoints", async function() {
                    await assertUnsupported(pointer, "minter", {});
                    await assertUnsupported(pointer, "marketing_info", {});
                    await assertUnsupported(pointer, "download_logo", {});
                    await assertUnsupported(pointer, "all_allowances", { owner: accounts[0].sheAddress });
                    await assertUnsupported(pointer, "all_accounts", {});
                });
            });

            describe("execute", function() {
                it("should transfer token", async function() {
                    const respBefore = await queryWasm(pointer, "balance", {address: accounts[1].sheAddress});
                    const balanceBefore = respBefore.data.balance;

                    const res = await executeWasm(pointer,  { transfer: { recipient: accounts[1].sheAddress, amount: "100" } });
                    const txHash = res["txhash"];
                    const receipt = await ethers.provider.getTransactionReceipt(`0x${txHash}`); 
                    expect(receipt).not.to.be.null;
                    console.log("receipt[\"blockNumber\"]", receipt["blockNumber"]);
                    const bn = receipt["blockNumber"];
                    const filter = {
                        fromBlock: '0x' + bn.toString(16),
                        toBlock: 'latest',
                        address: receipt["to"],
                        topics: [ethers.id("Transfer(address,address,uint256)")]
                    };
                    // send via eth_ endpoint - synthetic event doesn't show up
                    const ethlogs = await ethers.provider.send('eth_getLogs', [filter]);
                    expect(ethlogs.length).to.equal(0);

                    // send via she_ endpoint - synthetic event shows up
                    const shelogs = await ethers.provider.send('she_getLogs', [filter]);
                    expect(shelogs.length).to.equal(1);
                    expect(shelogs[0]["topics"][0]).to.equal(ethers.id("Transfer(address,address,uint256)"));
                    const respAfter = await queryWasm(pointer, "balance", {address: accounts[1].sheAddress});
                    const balanceAfter = respAfter.data.balance;
                    expect(balanceAfter).to.equal((parseInt(balanceBefore) + 100).toString());
                });

                it("transfer to unassociated address should fail", async function() {
                    const unassociatedSheAddr = "she1z7qugn2xy4ww0c9nsccftxw592n4xhxccmcf4q";
                    const respBefore = await queryWasm(pointer, "balance", {address: accounts[1].sheAddress});
                    const balanceBefore = respBefore.data.balance;

                    await executeWasm(pointer,  { transfer: { recipient: unassociatedSheAddr, amount: "100" } });
                    const respAfter = await queryWasm(pointer, "balance", {address: accounts[1].sheAddress});
                    const balanceAfter = respAfter.data.balance;

                    expect(balanceAfter).to.equal(balanceBefore);
                });

                it("transfer to contract address should succeed", async function() {
                    await associateWasm(pointer);
                    const respBefore = await queryWasm(pointer, "balance", {address: admin.sheAddress});
                    const balanceBefore = respBefore.data.balance;

                    await executeWasm(pointer,  { transfer: { recipient: pointer, amount: "100" } });
                    const respAfter = await queryWasm(pointer, "balance", {address: admin.sheAddress});
                    const balanceAfter = respAfter.data.balance;

                    expect(balanceAfter).to.equal((parseInt(balanceBefore) - 100).toString());
                });

                it("should increase and decrease allowance for a spender", async function() {
                    const spender = accounts[0].sheAddress;
                    await executeWasm(pointer, { increase_allowance: { spender: spender, amount: "300" } });

                    let allowance = await queryWasm(pointer, "allowance", { owner: admin.sheAddress, spender: spender });
                    expect(allowance.data.allowance).to.equal("300");

                    await executeWasm(pointer, { decrease_allowance: { spender: spender, amount: "300" } });

                    allowance = await queryWasm(pointer, "allowance", { owner: admin.sheAddress, spender: spender });
                    expect(allowance.data.allowance).to.equal("0");
                });

                it("should transfer token using transferFrom", async function() {
                    const resp = await testToken.approve(admin.evmAddress, 100);
                    await resp.wait();
                    const respBefore = await queryWasm(pointer, "balance", {address: accounts[0].sheAddress});
                    const balanceBefore = respBefore.data.balance;
                    await executeWasm(pointer,  { transfer_from: { owner: accounts[0].sheAddress, recipient: accounts[1].sheAddress, amount: "100" } });
                    const respAfter = await queryWasm(pointer, "balance", {address: accounts[0].sheAddress});
                    const balanceAfter = respAfter.data.balance;
                    expect(balanceAfter).to.equal((parseInt(balanceBefore) - 100).toString());
                });

                it("should transfer if called through wasmd precompile", async function() {
                    const WasmPrecompileContract = '0x0000000000000000000000000000000000001002';
                    const contractABIPath = '../../precompiles/wasmd/abi.json';
                    const contractABI = require(contractABIPath);
                    wasmd = new ethers.Contract(WasmPrecompileContract, contractABI, accounts[0].signer);

                    const encoder = new TextEncoder();

                    const transferMsg = { transfer: { recipient: accounts[1].sheAddress, amount: "100" } };
                    const transferStr = JSON.stringify(transferMsg);
                    const transferBz = encoder.encode(transferStr);

                    const coins = [];
                    const coinsStr = JSON.stringify(coins);
                    const coinsBz = encoder.encode(coinsStr);

                    const response = await wasmd.execute(pointer, transferBz, coinsBz);
                    const receipt = await response.wait();
                    expect(receipt.status).to.equal(1);

                    const filter = {
                        fromBlock: receipt["blockNumber"],
                        toBlock: 'latest',
                        topics: [ethers.id("Transfer(address,address,uint256)")]
                    };
                    const logs = await ethers.provider.getLogs(filter);
                    expect(logs.length).to.equal(1);
                });
            });
        });
    }

    describe("Pointer Functionality", function () {
        let pointer;

        before(async function () {
            pointer = cw20Pointer;
        });

        // Verify pointer
        testPointer(() => pointer, {
            0: 1000000000000,
            1: 1000000000000
        });

        // Pointer version is going to be coupled with blkd version going forward (as in,
        // given a blkd version, it's impossible to have multiple versions of pointer).
        // We need to recreate the equivalent of the following test once we have a framework
        // for simulating chain-level upgrade.
        describe.skip("Pointer Upgrade", function () {
            let newPointer;

            before(async function () {
               const tokenAddr = await testToken.getAddress();
               newPointer = await deployWasm(WASM.POINTER_CW20, accounts[0].sheAddress, "cw20-erc20", {erc20_address: tokenAddr })
               await proposeCW20toERC20Upgrade(tokenAddr, newPointer)
            });

            // Verify new pointer
            testPointer(() => newPointer, {
                0: 1000000000000,
                1: 1000000000000
            });
        });

        // The original pointer does not work now (expected)
        // test is configured to skip until original pointer works (unimplemented)
        describe.skip("Original Pointer after Upgrade", function(){
            // Original pointer
            testPointer(() => pointer, {
                0: 1000000000000,
                1: 1000000000000
            });
        });
    });
});

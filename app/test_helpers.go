package app

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/baseapp"
	crptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	minttypes "github.com/bluelink-lab/blk-chain/x/mint/types"
)

const TestContract = "TEST"
const TestUser = "she1jdppe6fnj2q7hjsepty5crxtrryzhuqsjrj95y"

type TestTx struct {
	msgs []sdk.Msg
}

func NewTestTx(msgs []sdk.Msg) TestTx {
	return TestTx{msgs: msgs}
}

func (t TestTx) GetMsgs() []sdk.Msg {
	return t.msgs
}

func (t TestTx) ValidateBasic() error {
	return nil
}

func (t TestTx) GetGasEstimate() uint64 {
	return 0
}

type TestAppOpts struct {
	useSc bool
}

func (t TestAppOpts) Get(s string) interface{} {
	if s == "chain-id" {
		return "she-test"
	}
	if s == FlagSCEnable {
		return t.useSc
	}
	return nil
}

type TestWrapper struct {
	suite.Suite

	App *App
	Ctx sdk.Context
}

func NewTestWrapper(t *testing.T, tm time.Time, valPub crptotypes.PubKey, enableEVMCustomPrecompiles bool, baseAppOptions ...func(*baseapp.BaseApp)) *TestWrapper {
	return newTestWrapper(t, tm, valPub, enableEVMCustomPrecompiles, false, baseAppOptions...)
}

func NewTestWrapperWithSc(t *testing.T, tm time.Time, valPub crptotypes.PubKey, enableEVMCustomPrecompiles bool, baseAppOptions ...func(*baseapp.BaseApp)) *TestWrapper {
	return newTestWrapper(t, tm, valPub, enableEVMCustomPrecompiles, true, baseAppOptions...)
}

func newTestWrapper(t *testing.T, tm time.Time, valPub crptotypes.PubKey, enableEVMCustomPrecompiles bool, useSc bool, baseAppOptions ...func(*baseapp.BaseApp)) *TestWrapper {
	var appPtr *App
	if useSc {
		appPtr = SetupWithSc(false, enableEVMCustomPrecompiles, baseAppOptions...)
	} else {
		appPtr = Setup(false, enableEVMCustomPrecompiles, baseAppOptions...)
	}
	ctx := appPtr.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "she-test", Time: tm})
	wrapper := &TestWrapper{
		App: appPtr,
		Ctx: ctx,
	}
	wrapper.SetT(t)
	wrapper.setupValidator(stakingtypes.Unbonded, valPub)
	return wrapper
}

func (s *TestWrapper) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) {
	err := s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, amounts)
	s.Require().NoError(err)

	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, minttypes.ModuleName, acc, amounts)
	s.Require().NoError(err)
}

func (s *TestWrapper) setupValidator(bondStatus stakingtypes.BondStatus, valPub crptotypes.PubKey) sdk.ValAddress {
	valAddr := sdk.ValAddress(valPub.Address())
	bondDenom := s.App.StakingKeeper.GetParams(s.Ctx).BondDenom
	selfBond := sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(100), Denom: bondDenom})

	s.FundAcc(sdk.AccAddress(valAddr), selfBond)

	sh := teststaking.NewHelper(s.Suite.T(), s.Ctx, s.App.StakingKeeper)
	msg := sh.CreateValidatorMsg(valAddr, valPub, selfBond[0].Amount)
	sh.Handle(msg, true)

	val, found := s.App.StakingKeeper.GetValidator(s.Ctx, valAddr)
	s.Require().True(found)

	val = val.UpdateStatus(bondStatus)
	s.App.StakingKeeper.SetValidator(s.Ctx, val)

	consAddr, err := val.GetConsAddr()
	s.Suite.Require().NoError(err)

	signingInfo := slashingtypes.NewValidatorSigningInfo(
		consAddr,
		s.Ctx.BlockHeight(),
		0,
		time.Unix(0, 0),
		false,
		0,
	)
	s.App.SlashingKeeper.SetValidatorSigningInfo(s.Ctx, consAddr, signingInfo)

	return valAddr
}

func (s *TestWrapper) BeginBlock() {
	var proposer sdk.ValAddress

	validators := s.App.StakingKeeper.GetAllValidators(s.Ctx)
	s.Require().Equal(1, len(validators))

	valAddrFancy, err := validators[0].GetConsAddr()
	s.Require().NoError(err)
	proposer = valAddrFancy.Bytes()

	validator, found := s.App.StakingKeeper.GetValidator(s.Ctx, proposer)
	s.Assert().True(found)

	valConsAddr, err := validator.GetConsAddr()

	s.Require().NoError(err)

	valAddr := valConsAddr.Bytes()

	newBlockTime := s.Ctx.BlockTime().Add(2 * time.Second)

	header := tmproto.Header{Height: s.Ctx.BlockHeight() + 1, Time: newBlockTime}
	newCtx := s.Ctx.WithBlockTime(newBlockTime).WithBlockHeight(s.Ctx.BlockHeight() + 1)
	s.Ctx = newCtx
	lastCommitInfo := abci.LastCommitInfo{
		Votes: []abci.VoteInfo{{
			Validator:       abci.Validator{Address: valAddr, Power: 1000},
			SignedLastBlock: true,
		}},
	}
	reqBeginBlock := abci.RequestBeginBlock{Header: header, LastCommitInfo: lastCommitInfo}

	s.App.BeginBlocker(s.Ctx, reqBeginBlock)
}

func (s *TestWrapper) EndBlock() {
	reqEndBlock := abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}
	s.App.EndBlocker(s.Ctx, reqEndBlock)
}

func Setup(isCheckTx bool, enableEVMCustomPrecompiles bool, baseAppOptions ...func(*baseapp.BaseApp)) (res *App) {
	db := dbm.NewMemDB()
	encodingConfig := MakeEncodingConfig()
	cdc := encodingConfig.Marshaler

	options := []AppOption{
		func(app *App) {
			app.receiptStore = NewInMemoryStateStore()
		},
	}

	res = New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		DefaultNodeHome,
		1,
		enableEVMCustomPrecompiles,
		config.TestConfig(),
		encodingConfig,
		wasm.EnableAllProposals,
		TestAppOpts{},
		EmptyWasmOpts,
		EmptyACLOpts,
		options,
		baseAppOptions...,
	)
	if !isCheckTx {
		genesisState := NewDefaultGenesisState(cdc)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		_, err = res.InitChain(
			context.Background(), &abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
		if err != nil {
			panic(err)
		}
	}

	return res
}

func SetupWithSc(isCheckTx bool, enableEVMCustomPrecompiles bool, baseAppOptions ...func(*baseapp.BaseApp)) (res *App) {
	db := dbm.NewMemDB()
	encodingConfig := MakeEncodingConfig()
	cdc := encodingConfig.Marshaler

	options := []AppOption{
		func(app *App) {
			app.receiptStore = NewInMemoryStateStore()
		},
	}

	res = New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		DefaultNodeHome,
		1,
		enableEVMCustomPrecompiles,
		config.TestConfig(),
		encodingConfig,
		wasm.EnableAllProposals,
		TestAppOpts{true},
		EmptyWasmOpts,
		EmptyACLOpts,
		options,
		baseAppOptions...,
	)
	if !isCheckTx {
		genesisState := NewDefaultGenesisState(cdc)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// TODO: remove once init chain works with SC
		defer func() { _ = recover() }()

		_, err = res.InitChain(
			context.Background(), &abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
		if err != nil {
			panic(err)
		}
	}

	return res
}

func SetupTestingAppWithLevelDb(isCheckTx bool, enableEVMCustomPrecompiles bool) (*App, func()) {
	dir := "blt_testing"
	db, err := sdk.NewLevelDB("blt_leveldb_testing", dir)
	if err != nil {
		panic(err)
	}
	encodingConfig := MakeEncodingConfig()
	cdc := encodingConfig.Marshaler
	app := New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		DefaultNodeHome,
		5,
		enableEVMCustomPrecompiles,
		nil,
		encodingConfig,
		wasm.EnableAllProposals,
		TestAppOpts{},
		EmptyWasmOpts,
		EmptyACLOpts,
		nil,
	)
	if !isCheckTx {
		genesisState := NewDefaultGenesisState(cdc)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		_, err = app.InitChain(
			context.Background(), &abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
		if err != nil {
			panic(err)
		}
	}

	cleanupFn := func() {
		db.Close()
		err = os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}

	return app, cleanupFn
}

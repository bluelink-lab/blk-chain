package aclmapping

import (
	aclkeeper "github.com/cosmos/cosmos-sdk/x/accesscontrol/keeper"
	aclbankmapping "github.com/bluelink-lab/blk-chain/aclmapping/bank"
	aclevmmapping "github.com/bluelink-lab/blk-chain/aclmapping/evm"
	acloraclemapping "github.com/bluelink-lab/blk-chain/aclmapping/oracle"
	acltokenfactorymapping "github.com/bluelink-lab/blk-chain/aclmapping/tokenfactory"
	aclwasmmapping "github.com/bluelink-lab/blk-chain/aclmapping/wasm"
	evmkeeper "github.com/bluelink-lab/blk-chain/x/evm/keeper"
)

type CustomDependencyGenerator struct{}

func NewCustomDependencyGenerator() CustomDependencyGenerator {
	return CustomDependencyGenerator{}
}

func (customDepGen CustomDependencyGenerator) GetCustomDependencyGenerators(evmKeeper evmkeeper.Keeper) aclkeeper.DependencyGeneratorMap {
	dependencyGeneratorMap := make(aclkeeper.DependencyGeneratorMap)
	wasmDependencyGenerators := aclwasmmapping.NewWasmDependencyGenerator()

	dependencyGeneratorMap = dependencyGeneratorMap.Merge(aclbankmapping.GetBankDepedencyGenerator())
	dependencyGeneratorMap = dependencyGeneratorMap.Merge(acltokenfactorymapping.GetTokenFactoryDependencyGenerators())
	dependencyGeneratorMap = dependencyGeneratorMap.Merge(wasmDependencyGenerators.GetWasmDependencyGenerators())
	dependencyGeneratorMap = dependencyGeneratorMap.Merge(acloraclemapping.GetOracleDependencyGenerator())
	dependencyGeneratorMap = dependencyGeneratorMap.Merge(aclevmmapping.GetEVMDependencyGenerators(evmKeeper))

	return dependencyGeneratorMap
}

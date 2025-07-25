package types

//nolint:typecheck
import (
	"errors"
	fmt "fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/proto"
	"github.com/bluelink-lab/blk-chain/x/evm/types/ethtx"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

func GetAmino() *codec.LegacyAmino {
	return amino
}

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAssociate{}, "evm/MsgAssociate", nil)
	cdc.RegisterConcrete(&MsgEVMTransaction{}, "evm/MsgEVMTransaction", nil)
	cdc.RegisterConcrete(&MsgSend{}, "evm/MsgSend", nil)
	cdc.RegisterConcrete(&MsgRegisterPointer{}, "evm/MsgRegisterPointer", nil)
	cdc.RegisterConcrete(&MsgAssociateContractAddress{}, "evm/MsgAssociateContractAddress", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*govtypes.Content)(nil),
		&AddERCNativePointerProposal{},
		&AddERCCW20PointerProposal{},
		&AddERCCW721PointerProposal{},
		&AddERCCW1155PointerProposal{},
		&AddCWERC20PointerProposal{},
		&AddCWERC721PointerProposal{},
		&AddCWERC1155PointerProposal{},
		&AddERCNativePointerProposalV2{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgEVMTransaction{},
		&MsgSend{},
		&MsgRegisterPointer{},
		&MsgAssociateContractAddress{},
	)
	registry.RegisterInterface(
		"sheprotocol.blk-chain.evm.TxData",
		(*ethtx.TxData)(nil),
		&ethtx.DynamicFeeTx{},
		&ethtx.AccessListTx{},
		&ethtx.LegacyTx{},
		&ethtx.BlobTx{},
		&ethtx.AssociateTx{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func PackTxData(txData ethtx.TxData) (*codectypes.Any, error) {
	msg, ok := txData.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("cannot proto marshal %T", txData)
	}

	anyTxData, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return anyTxData, nil
}

func UnpackTxData(any *codectypes.Any) (ethtx.TxData, error) {
	if any == nil {
		return nil, errors.New("protobuf Any message cannot be nil")
	}

	txData, ok := any.GetCachedValue().(ethtx.TxData)
	if !ok {
		ltx := ethtx.LegacyTx{}
		if proto.Unmarshal(any.Value, &ltx) == nil {
			// value is a legacy tx
			return &ltx, nil
		}
		atx := ethtx.AccessListTx{}
		if proto.Unmarshal(any.Value, &atx) == nil {
			// value is a accesslist tx
			return &atx, nil
		}
		dtx := ethtx.DynamicFeeTx{}
		if proto.Unmarshal(any.Value, &dtx) == nil {
			// value is a dynamic fee tx
			return &dtx, nil
		}
		btx := ethtx.BlobTx{}
		if proto.Unmarshal(any.Value, &btx) == nil {
			// value is a blob tx
			return &btx, nil
		}
		astx := ethtx.AssociateTx{}
		if proto.Unmarshal(any.Value, &astx) == nil {
			// value is an associate tx
			return &astx, nil
		}
		return nil, fmt.Errorf("cannot unpack Any into TxData %T", any)
	}

	return txData, nil
}

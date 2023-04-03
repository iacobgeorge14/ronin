package vm

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TestOpEvent struct {
}

func (tx *TestOpEvent) Publish(
	opcode OpCode,
	order, blockHeight uint64,
	blockHash common.Hash,
	blockTime uint64,
	hash common.Hash,
	from, to common.Address,
	value *big.Int,
	input []byte,
	err error,
) *types.InternalTransaction {
	return &types.InternalTransaction{
		Opcode:          opcode.String(),
		Order:           order,
		TransactionHash: hash,
		Type:            "test",
		Value:           value,
		Input:           input,
		From:            from,
		To:              to,
		Success:         err == nil,
		Error:           "",
		Height:          blockHeight,
		BlockHash:       blockHash,
		BlockTime:       blockTime,
	}
}

func TestPublishEvents(t *testing.T) {
	ctx := BlockContext{
		PublishEvents: map[OpCode]OpEvent{
			CALL: &TestOpEvent{},
		},
		CurrentTransaction: types.NewTx(&types.LegacyTx{
			Nonce:    1,
			To:       nil,
			Value:    big.NewInt(0),
			Gas:      0,
			GasPrice: big.NewInt(0),
			Data:     []byte(""),
		}),
		BlockNumber:          common.Big0,
		Time:                 common.Big0,
		InternalTransactions: &[]*types.InternalTransaction{},
	}

	evm := &EVM{Context: ctx}
	evm.PublishEvent(CALL, 1, common.Address{}, common.Address{}, big.NewInt(0), []byte(""), nil)
	if len(*evm.Context.InternalTransactions) != 1 || (*evm.Context.InternalTransactions)[0].Type != "test" {
		t.Error("Failed to publish opcode event")
	}
}

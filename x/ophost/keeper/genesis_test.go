package keeper_test

import (
	"testing"
	"time"

	"github.com/initia-labs/OPinit/x/ophost/types"
	"github.com/stretchr/testify/require"
)

func Test_GenesisImportExport(t *testing.T) {
	ctx, input := createDefaultTestInput(t)

	params := input.OPHostKeeper.GetParams(ctx)
	config1 := types.BridgeConfig{
		Challenger:          "challenger",
		Proposer:            "proposer",
		SubmissionInterval:  100,
		FinalizationPeriod:  100,
		SubmissionStartTime: time.Now().UTC(),
		Metadata:            []byte{1, 2, 3},
	}
	config2 := types.BridgeConfig{
		Challenger:          "challenger2",
		Proposer:            "proposer2",
		SubmissionInterval:  200,
		FinalizationPeriod:  200,
		SubmissionStartTime: time.Now().UTC(),
		Metadata:            []byte{3, 4, 5},
	}
	require.NoError(t, input.OPHostKeeper.SetBridgeConfig(ctx, 1, config1))
	require.NoError(t, input.OPHostKeeper.SetBridgeConfig(ctx, 2, config2))

	input.OPHostKeeper.SetNextBridgeId(ctx, 3)
	input.OPHostKeeper.SetNextL1Sequence(ctx, 1, 100)
	input.OPHostKeeper.SetNextL1Sequence(ctx, 2, 200)
	input.OPHostKeeper.SetNextOutputIndex(ctx, 1, 10)
	input.OPHostKeeper.SetNextOutputIndex(ctx, 2, 20)

	output1 := types.Output{
		OutputRoot:    []byte{1, 2, 3},
		L1BlockTime:   time.Now().UTC(),
		L2BlockNumber: 100,
	}
	output2 := types.Output{
		OutputRoot:    []byte{1, 2, 3},
		L1BlockTime:   time.Now().UTC(),
		L2BlockNumber: 200,
	}
	output3 := types.Output{
		OutputRoot:    []byte{1, 2, 3},
		L1BlockTime:   time.Now().UTC(),
		L2BlockNumber: 100,
	}
	require.NoError(t, input.OPHostKeeper.SetOutputProposal(ctx, 1, 1, output1))
	require.NoError(t, input.OPHostKeeper.SetOutputProposal(ctx, 1, 2, output2))
	require.NoError(t, input.OPHostKeeper.SetOutputProposal(ctx, 2, 1, output3))

	input.OPHostKeeper.SetTokenPair(ctx, 1, "l2denom", "l1denom")
	input.OPHostKeeper.SetTokenPair(ctx, 2, "l12denom", "l11denom")

	input.OPHostKeeper.RecordProvenWithdrawal(ctx, 1, [32]byte{1, 2, 3})
	input.OPHostKeeper.RecordProvenWithdrawal(ctx, 1, [32]byte{3, 4, 5})

	genState := input.OPHostKeeper.ExportGenesis(ctx)
	input.OPHostKeeper.InitGenesis(ctx, genState)

	_genState := input.OPHostKeeper.ExportGenesis(ctx)
	require.Equal(t, genState, _genState)

	require.Equal(t, uint64(3), genState.NextBridgeId)
	require.Equal(t, params, genState.Params)
	require.Equal(t, types.Bridge{
		BridgeId:        1,
		NextL1Sequence:  100,
		NextOutputIndex: 10,
		BridgeConfig:    config1,
		TokenPairs: []types.TokenPair{
			{
				L1Denom: "l1denom",
				L2Denom: "l2denom",
			},
		},
		Proposals: []types.WrappedOutput{
			{
				OutputIndex:    1,
				OutputProposal: output1,
			},
			{
				OutputIndex:    2,
				OutputProposal: output2,
			},
		},
		ProvenWithdrawals: [][]byte{
			{1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{3, 4, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}, genState.Bridges[0])
}
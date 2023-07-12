package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ibc-go/v7/modules/apps/callbacks/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
)

func (suite *CallbacksTypesTestSuite) TestLogger() {
	suite.SetupSuite()

	ctx := suite.chain.GetContext()

	suite.Require().Equal(
		ctx.Logger().With("module", "x/"+types.ModuleName),
		types.Logger(ctx))
}

func (suite *CallbacksTypesTestSuite) TestEvents() {
	testCases := []struct {
		name          string
		packet        channeltypes.Packet
		callbackType  types.CallbackType
		callbackData  types.CallbackData
		callbackError error
		expEvents     ibctesting.EventsMap
	}{
		{
			"success: ack callback",
			channeltypes.NewPacket(
				ibctesting.MockPacketData, 1, ibctesting.MockPort, ibctesting.FirstChannelID,
				ibctesting.MockFeePort, ibctesting.InvalidID, clienttypes.NewHeight(1, 100), 0,
			),
			types.CallbackTypeAcknowledgement,
			types.CallbackData{
				ContractAddr: ibctesting.TestAccAddress,
				GasLimit:     100000,
			},
			nil,
			ibctesting.EventsMap{
				types.EventTypeSourceCallback: {
					sdk.AttributeKeyModule:                    types.ModuleName,
					types.AttributeKeyCallbackTrigger:         string(types.CallbackTypeAcknowledgement),
					types.AttributeKeyCallbackAddress:         ibctesting.TestAccAddress,
					types.AttributeKeyCallbackGasLimit:        "100000",
					types.AttributeKeyCallbackSourcePortID:    ibctesting.MockPort,
					types.AttributeKeyCallbackSourceChannelID: ibctesting.FirstChannelID,
					types.AttributeKeyCallbackSequence:        "1",
					types.AttributeKeyCallbackResult:          "success",
				},
			},
		},
		{
			"success: timeout callback",
			channeltypes.NewPacket(
				ibctesting.MockPacketData, 1, ibctesting.MockPort, ibctesting.FirstChannelID,
				ibctesting.MockFeePort, ibctesting.InvalidID, clienttypes.NewHeight(1, 100), 0,
			),
			types.CallbackTypeTimeoutPacket,
			types.CallbackData{
				ContractAddr: ibctesting.TestAccAddress,
				GasLimit:     100000,
			},
			nil,
			ibctesting.EventsMap{
				types.EventTypeSourceCallback: {
					sdk.AttributeKeyModule:                    types.ModuleName,
					types.AttributeKeyCallbackTrigger:         string(types.CallbackTypeTimeoutPacket),
					types.AttributeKeyCallbackAddress:         ibctesting.TestAccAddress,
					types.AttributeKeyCallbackGasLimit:        "100000",
					types.AttributeKeyCallbackSourcePortID:    ibctesting.MockPort,
					types.AttributeKeyCallbackSourceChannelID: ibctesting.FirstChannelID,
					types.AttributeKeyCallbackSequence:        "1",
					types.AttributeKeyCallbackResult:          "success",
				},
			},
		},
		{
			"success: timeout callback",
			channeltypes.NewPacket(
				ibctesting.MockPacketData, 1, ibctesting.MockPort, ibctesting.FirstChannelID,
				ibctesting.MockFeePort, ibctesting.InvalidID, clienttypes.NewHeight(1, 100), 0,
			),
			types.CallbackTypeReceivePacket,
			types.CallbackData{
				ContractAddr: ibctesting.TestAccAddress,
				GasLimit:     100000,
			},
			nil,
			ibctesting.EventsMap{
				types.EventTypeDestinationCallback: {
					sdk.AttributeKeyModule:                    types.ModuleName,
					types.AttributeKeyCallbackTrigger:         string(types.CallbackTypeReceivePacket),
					types.AttributeKeyCallbackAddress:         ibctesting.TestAccAddress,
					types.AttributeKeyCallbackGasLimit:        "100000",
					types.AttributeKeyCallbackSourcePortID:    ibctesting.MockPort,
					types.AttributeKeyCallbackSourceChannelID: ibctesting.FirstChannelID,
					types.AttributeKeyCallbackSequence:        "1",
					types.AttributeKeyCallbackResult:          "success",
				},
			},
		},
		{
			"success: unknown callback, unreachable code",
			channeltypes.NewPacket(
				ibctesting.MockPacketData, 1, ibctesting.MockPort, ibctesting.FirstChannelID,
				ibctesting.MockFeePort, ibctesting.InvalidID, clienttypes.NewHeight(1, 100), 0,
			),
			"something",
			types.CallbackData{
				ContractAddr: ibctesting.TestAccAddress,
				GasLimit:     100000,
			},
			nil,
			ibctesting.EventsMap{
				"unknown": {
					sdk.AttributeKeyModule:                    types.ModuleName,
					types.AttributeKeyCallbackTrigger:         "something",
					types.AttributeKeyCallbackAddress:         ibctesting.TestAccAddress,
					types.AttributeKeyCallbackGasLimit:        "100000",
					types.AttributeKeyCallbackSourcePortID:    ibctesting.MockPort,
					types.AttributeKeyCallbackSourceChannelID: ibctesting.FirstChannelID,
					types.AttributeKeyCallbackSequence:        "1",
					types.AttributeKeyCallbackResult:          "success",
				},
			},
		},
		{
			"failure: ack callback with error",
			channeltypes.NewPacket(
				ibctesting.MockPacketData, 1, ibctesting.MockPort, ibctesting.FirstChannelID,
				ibctesting.MockFeePort, ibctesting.InvalidID, clienttypes.NewHeight(1, 100), 0,
			),
			types.CallbackTypeAcknowledgement,
			types.CallbackData{
				ContractAddr: ibctesting.TestAccAddress,
				GasLimit:     100000,
			},
			types.ErrNotCallbackPacketData,
			ibctesting.EventsMap{
				types.EventTypeSourceCallback: {
					sdk.AttributeKeyModule:                    types.ModuleName,
					types.AttributeKeyCallbackTrigger:         string(types.CallbackTypeAcknowledgement),
					types.AttributeKeyCallbackAddress:         ibctesting.TestAccAddress,
					types.AttributeKeyCallbackGasLimit:        "100000",
					types.AttributeKeyCallbackSourcePortID:    ibctesting.MockPort,
					types.AttributeKeyCallbackSourceChannelID: ibctesting.FirstChannelID,
					types.AttributeKeyCallbackSequence:        "1",
					types.AttributeKeyCallbackResult:          "failure",
					types.AttributeKeyCallbackError:           types.ErrNotCallbackPacketData.Error(),
				},
			},
		},
	}

	for _, tc := range testCases {
		newCtx := sdk.Context{}.WithEventManager(sdk.NewEventManager())

		types.EmitCallbackEvent(newCtx, tc.packet, tc.callbackType, tc.callbackData, tc.callbackError)
		events := newCtx.EventManager().Events().ToABCIEvents()
		ibctesting.AssertEvents(&suite.Suite, tc.expEvents, events)
	}
}

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const BlocksPerWeek = 104275

func (k *Keeper) isEpochEnd(ctx sdk.Context, epochLength int64) bool {
	return ctx.BlockHeight()%epochLength == 0
}

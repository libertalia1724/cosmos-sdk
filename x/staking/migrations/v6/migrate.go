package v6

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func MigrateStore(ctx context.Context, store storetypes.KVStore, cdc codec.BinaryCodec) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return migrateParams(sdkCtx, store, cdc)
}

func migrateParams(ctx sdk.Context, store storetypes.KVStore, cdc codec.BinaryCodec) error {
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		params := types.DefaultParams()
		bz, err := cdc.Marshal(&params)
		if err != nil {
			return err
		}
		store.Set(types.ParamsKey, bz)
		return nil
	}

	var params types.Params
	if err := cdc.Unmarshal(bz, &params); err != nil {
		return err
	}

	if params.MinCommissionRate.IsNil() {
		params.MinCommissionRate = types.DefaultMinCommissionRate
	}
	if params.CommissionThresholdRate.IsNil() {
		params.CommissionThresholdRate = types.DefaultCommissionThresholdRate
	}
	if params.MaxCommissionRate.IsNil() {
		params.MaxCommissionRate = types.DefaultMaxCommissionRate
	}
	if params.CommissionRateMultiplier.IsNil() {
		params.CommissionRateMultiplier = types.DefaultCommissionRateMultiplier
	}

	bz, err := cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)
	return nil
}

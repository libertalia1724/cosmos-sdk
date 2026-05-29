package keeper

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k *Keeper) SetPreviousCommissionRate(ctx sdk.Context, valAddr sdk.ValAddress, rate math.LegacyDec) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := rate.Marshal()
	if err != nil {
		return err
	}
	return store.Set(types.GetPreviousCommissionRateKey(valAddr), bz)
}

func (k *Keeper) GetPreviousCommissionRate(ctx sdk.Context, valAddr sdk.ValAddress) (math.LegacyDec, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.GetPreviousCommissionRateKey(valAddr))
	if err != nil {
		return math.LegacyZeroDec(), err
	}
	if bz == nil {
		return math.LegacyZeroDec(), nil
	}
	var rate math.LegacyDec
	if err := rate.Unmarshal(bz); err != nil {
		return math.LegacyZeroDec(), err
	}
	return rate, nil
}

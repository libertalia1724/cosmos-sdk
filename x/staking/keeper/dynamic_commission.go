package keeper

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k *Keeper) calculateNewCommission(ctx sdk.Context, valAddr sdk.ValAddress) (math.LegacyDec, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	validatorVotingPower, err := k.GetValidatorVotingPowerPercent(ctx, valAddr)
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	new_min_commission := math.LegacyMinDec(params.MaxCommissionRate, math.LegacyMaxDec(params.MinCommissionRate, params.MinCommissionRate.Add(params.CommissionRateMultiplier.Mul(validatorVotingPower.Sub(params.CommissionThresholdRate)))))

	return new_min_commission, err
}

func (k *Keeper) GetValidatorVotingPowerPercent(ctx sdk.Context, valAddr sdk.ValAddress) (math.LegacyDec, error) {
	validator, err := k.GetValidator(ctx, valAddr)
	if err != nil {
		return math.LegacyZeroDec(), err
	}
	// validator's power
	validatorPower := math.LegacyNewDecFromInt(validator.GetTokens())

	// total bonded tokens

	totalBonded, err := k.TotalBondedTokens(ctx)
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	totalPower := math.LegacyNewDecFromInt(totalBonded)

	if totalPower.IsZero() {
		return math.LegacyZeroDec(), nil
	}

	// percent = validatorPower / totalPower
	return validatorPower.Quo(totalPower), nil
}

func (k *Keeper) UpdateAllValidatorCommissions(ctx sdk.Context) error {
	return k.IterateBondedValidatorsByPower(ctx, func(index int64, validator types.ValidatorI) bool {
		valAddr, err := sdk.ValAddressFromBech32(validator.GetOperator())
		if err != nil {
			return true
		}

		newCommission, err := k.calculateNewCommission(ctx, valAddr)
		if err != nil {
			return true
		}

		val, ok := validator.(types.Validator)
		if !ok {
			return true
		}

		oldCommissionRate := val.Commission.Rate
		k.SetPreviousCommissionRate(ctx, valAddr, oldCommissionRate)

		if val.Commission.Rate.LT(newCommission) {
			val.Commission.Rate = newCommission
		}

		if err := k.SetValidator(ctx, val); err != nil {
			return true
		}

		return false
	})
}

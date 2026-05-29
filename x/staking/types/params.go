package types

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Staking params default values
const (
	// DefaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	// TODO: Justify our choice of default here.
	DefaultUnbondingTime time.Duration = time.Hour * 24 * 7 * 3

	// Default maximum number of bonded validators
	DefaultMaxValidators uint32 = 100

	// Default maximum entries in a UBD/RED pair
	DefaultMaxEntries uint32 = 7

	// DefaultHistorical entries is 10000. Apps that don't use IBC can ignore this
	// value by not adding the staking module to the application module manager's
	// SetOrderBeginBlockers.
	DefaultHistoricalEntries uint32 = 10000
)

// DefaultMinCommissionRate is set to 0%
var DefaultMinCommissionRate = math.LegacyNewDecWithPrec(5, 2)

var DefaultCommissionThresholdRate = math.LegacyNewDecWithPrec(1, 2)

var DefaultMaxCommissionRate = math.LegacyNewDecWithPrec(20, 2)

var DefaultCommissionRateMultiplier = math.LegacyOneDec()

// NewParams creates a new Params instance
func NewParams(unbondingTime time.Duration, maxValidators, maxEntries, historicalEntries uint32, bondDenom string, minCommissionRate math.LegacyDec, commissionThresholdRate math.LegacyDec, maxCommissionRate math.LegacyDec, commissionRateMultiplier math.LegacyDec) Params {
	return Params{
		UnbondingTime:            unbondingTime,
		MaxValidators:            maxValidators,
		MaxEntries:               maxEntries,
		HistoricalEntries:        historicalEntries,
		BondDenom:                bondDenom,
		MinCommissionRate:        minCommissionRate,
		CommissionThresholdRate:  commissionThresholdRate,
		MaxCommissionRate:        maxCommissionRate,
		CommissionRateMultiplier: commissionRateMultiplier,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultUnbondingTime,
		DefaultMaxValidators,
		DefaultMaxEntries,
		DefaultHistoricalEntries,
		sdk.DefaultBondDenom,
		DefaultMinCommissionRate,
		DefaultCommissionThresholdRate,
		DefaultMaxCommissionRate,
		DefaultCommissionRateMultiplier,
	)
}

// unmarshal the current staking params value from store key or panic
func MustUnmarshalParams(cdc *codec.LegacyAmino, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}

	return params
}

// unmarshal the current staking params value from store key
func UnmarshalParams(cdc *codec.LegacyAmino, value []byte) (params Params, err error) {
	err = cdc.Unmarshal(value, &params)
	if err != nil {
		return params, err
	}

	return params, err
}

// validate a set of params
func (p Params) Validate() error {
	if err := validateUnbondingTime(p.UnbondingTime); err != nil {
		return err
	}

	if err := validateMaxValidators(p.MaxValidators); err != nil {
		return err
	}

	if err := validateMaxEntries(p.MaxEntries); err != nil {
		return err
	}

	if err := validateBondDenom(p.BondDenom); err != nil {
		return err
	}

	if err := validateMinCommissionRate(p.MinCommissionRate); err != nil {
		return err
	}

	if err := validateHistoricalEntries(p.HistoricalEntries); err != nil {
		return err
	}

	if err := validateCommissionThresholdRate(p.CommissionThresholdRate); err != nil {
		return err
	}

	if err := validateMaxCommissionRate(p.MaxCommissionRate); err != nil {
		return err
	}

	if err := validateCommissionRateMultiplier(p.CommissionRateMultiplier); err != nil {
		return err
	}

	return nil
}

func validateUnbondingTime(i any) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("unbonding time must be positive: %d", v)
	}

	return nil
}

func validateMaxValidators(i any) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max validators must be positive: %d", v)
	}

	return nil
}

func validateMaxEntries(i any) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max entries must be positive: %d", v)
	}

	return nil
}

func validateHistoricalEntries(i any) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateBondDenom(i any) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("bond denom cannot be blank")
	}

	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func ValidatePowerReduction(i any) error {
	v, ok := i.(math.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(math.NewInt(1)) {
		return fmt.Errorf("power reduction cannot be lower than 1")
	}

	return nil
}

func validateMinCommissionRate(i any) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("minimum commission rate cannot be nil: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("minimum commission rate cannot be negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("minimum commission rate cannot be greater than 100%%: %s", v)
	}

	return nil
}

func validateCommissionThresholdRate(i any) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("commission threshold rate cannot be nil: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("commission threshold rate cannot be negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("commission threshold rate cannot be greater than 100%%: %s", v)
	}

	return nil
}

func validateMaxCommissionRate(i any) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("maximum commission rate cannot be nil: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("maximum commission rate cannot be negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("maximum commission rate cannot be greater than 100%%: %s", v)
	}

	return nil
}

func validateCommissionRateMultiplier(i any) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("commission rate multiplier cannot be nil: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("commission rate multiplier cannot be negative: %s", v)
	}
	if v.Equal(math.LegacyZeroDec()) {
		return fmt.Errorf("commission rate multiplier cannot be equal zero: %s", v)
	}

	return nil
}

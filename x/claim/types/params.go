package types

import (
	fmt "fmt"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

var (
	DefaultClaimDenom         = "ueve"
	DefaultDurationUntilDecay = time.Hour
	DefaultDurationOfDecay    = time.Hour * 5
)

func NewParams(enabled bool, claimDenom string, startTime time.Time, durationUntilDecay, durationOfDecay time.Duration) Params {
	return Params{
		AirdropEnabled:     enabled,
		ClaimDenom:         claimDenom,
		AirdropStartTime:   startTime,
		DurationUntilDecay: durationUntilDecay,
		DurationOfDecay:    durationOfDecay,
	}
}

func DefaultParams() Params {
	return Params{
		AirdropEnabled:     false,
		AirdropStartTime:   time.Time{},
		DurationUntilDecay: DefaultDurationUntilDecay,
		DurationOfDecay:    DefaultDurationOfDecay,
		ClaimDenom:         DefaultClaimDenom,
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		return ""
	}
	return string(out)
}

func (p Params) ValidateBasic() error {
	var err error

	err = validateEnabled(p.AirdropEnabled)
	if err != nil {
		return err
	}
	err = validateDenom(p.ClaimDenom)
	if err != nil {
		return err
	}
	err = validateTime(p.AirdropStartTime)
	if err != nil {
		return err
	}
	err = validateDuration(p.DurationOfDecay)
	if err != nil {
		return err
	}
	err = validateDuration(p.DurationUntilDecay)
	if err != nil {
		return err
	}

	return nil
}

func (p Params) IsAirdropEnabled(t time.Time) bool {
	if !p.AirdropEnabled {
		return false
	}
	if p.AirdropStartTime.IsZero() {
		return false
	}
	if t.Before(p.AirdropStartTime) {
		return false
	}
	return true
}

func validateEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("invalid denom: %s", v)
	}

	return nil
}

func validateTime(i interface{}) error {
	_, ok := i.(time.Time)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateDuration(i interface{}) error {
	d, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if d < 1 {
		return fmt.Errorf("duration must be greater than or equal to 1: %d", d)
	}
	return nil
}

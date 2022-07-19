// decimal is a thin wrapper around github.com/cockroachdb/apd/v3.
//
// - provides support for JsonUnmarshal from both strings and integers.
//
// - a set of API that is similar to github.com/shopspring/decimal, where
// a new value is created instead of setting to an `*Decimal`.
//
// - support cobra cli value.
package decimal

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/cockroachdb/apd/v3"
)

// Decimal is a type around github.com/cockroachdb/apd/v3's Decimal.
type Decimal struct {
	apd.Decimal
}

// NewFromString is forwarded from `apd.NewFromString`.
func NewFromString(s string) (*Decimal, error) {
	d, _, err := apd.NewFromString(s)
	return &Decimal{Decimal: *d}, err
}

// UnmarshalJSON Unmarshals a json string or number into a Decimal.
//
// It first attempts to call the unmarshal method on `apd.Decimal`, which takes an
// integer, and `SetString` of `apd.Decimal` after failure
func (d *Decimal) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &d.Decimal)
	if err != nil {
		if _, _, err1 := d.Decimal.SetString(string(data)); err1 != nil {
			return fmt.Errorf("failed to parse decimal %s - both as str %v or int %v", string(data), err, err1)
		}
	}

	return nil
}

// Set is used to support cobra command line.
func (d *Decimal) Set(s string) error {
	_, cond, err := d.Decimal.SetString(s)
	if err != nil {
		return err
	}
	if cond.Any() {
		return fmt.Errorf("error setting string to decimal: %s", cond.String())
	}

	return nil
}

// Type is used to support cobra command line.
func (d *Decimal) Type() string {
	return "Decimal (wrapping github.com/cockroachdb/apd/v3 Decimal)"
}

// Clone a Decimal
func (d *Decimal) Clone() *Decimal {
	r := &Decimal{}
	r.Decimal.Set(&d.Decimal)
	return r
}

// 1 << 64, or 2^64
var twoTo64 = apd.NewBigInt(0).Lsh(apd.NewBigInt(1), 64)
var twoTo64Decimal = apd.NewWithBigInt(twoTo64, 0)

// NewIntFromUint64 converts uint64 to an *apd.BigInt
func NewBigIntFromUint64(i uint64) *apd.BigInt {
	return SetUint64ToBigInt(i, apd.NewBigInt(0))
}

// NewFromUint64 converts an uint64 to a *Decimal
func NewFromUint64(i uint64) *Decimal {
	r := &Decimal{}
	r.Decimal.Set(apd.NewWithBigInt(NewBigIntFromUint64(i), 0))
	return r
}

// SetUint64ToBigInt set to an apd.BigInt uint64
func SetUint64ToBigInt(i uint64, d *apd.BigInt) *apd.BigInt {
	signed := int64(i)
	if signed < 0 {
		d.SetInt64(signed)
		d.Add(d, twoTo64)
	} else {
		d.SetInt64(signed)
	}

	return d
}

// SetUint64 set the decimal to the uint64
func (d *Decimal) SetUint64(i uint64) *Decimal {
	SetUint64ToBigInt(i, &d.Decimal.Coeff)
	d.Decimal.Exponent = 0
	return d
}

// TryUint64 returns the uint64 representation of x, if x cannot be represented in
// an uint64, an error is returned.
func (x *Decimal) TryUint64() (uint64, error) {
	if x.Decimal.Negative {
		return 0, fmt.Errorf("%s: is less than zero", x.String())
	}

	d := &x.Decimal
	if d.Form != apd.Finite {
		return 0, fmt.Errorf("%s is not finite", d.String())
	}
	var integ, frac apd.Decimal
	d.Modf(&integ, &frac)
	if !frac.IsZero() {
		return 0, fmt.Errorf("%s: has fractional part", d.String())
	}
	var ed apd.ErrDecimal
	if integ.Cmp(twoTo64Decimal) >= 0 {
		return 0, fmt.Errorf("%s: greater than max int64", d.String())
	}

	if err := ed.Err(); err != nil {
		return 0, err
	}
	v := integ.Coeff.Int64()
	for i := int32(0); i < integ.Exponent; i++ {
		v *= 10
	}

	return uint64(v), nil
}

// Uint64 returns the uint64 representation of x. Panic if x
// cannot be presented as uint64
func (x *Decimal) Uint64() uint64 {
	d, err := x.TryUint64()
	if err != nil {
		panic(err)
	}
	return d
}

// TryBigInt tries to convert the decimal into a `big.Int`.
// Returns an error if the number is not an integer
func (x *Decimal) TryBigInt() (*big.Int, error) {
	var integ, frac apd.Decimal
	x.Modf(&integ, &frac)
	if !frac.IsZero() {
		return nil, fmt.Errorf("%s: has fractional part", x.String())
	}
	return x.Coeff.MathBigInt(), nil
}

// BigInt gets the `big.Int` representation of the decimal.
// Panics if x is not an integer. Use `TryBigInt`
// to catch the error.
func (x *Decimal) BigInt() *big.Int {
	r, err := x.TryBigInt()
	if err != nil {
		panic(err)
	}

	return r
}

// Copy set decimal to d
func (x *Decimal) Copy(d *Decimal) *Decimal {
	x.Decimal.Set(&d.Decimal)
	return x
}

// Zero of the decimal
var Zero = Decimal{}

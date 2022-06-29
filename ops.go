package decimal

import "github.com/cockroachdb/apd/v3"

// decimal128 style
var decimal128Context = apd.Context{
	Precision:   34,
	MaxExponent: apd.MaxExponent,
	MinExponent: apd.MinExponent,
	Traps:       apd.DefaultTraps,
}

func getOrPanic(d *Decimal, err error) *Decimal {
	if err != nil {
		panic(err)
	}

	return d
}

// TryDiv divides d by d1 (d/d1)
func (d *Decimal) TryDiv(d1 *Decimal) (*Decimal, error) {
	var r Decimal
	_, err := decimal128Context.Quo(&r.Decimal, &d.Decimal, &d1.Decimal)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Div divides d by d1 (d/d1), panics if there is an error.
// Use TryDiv to get error
func (d *Decimal) Div(d1 *Decimal) *Decimal {
	return getOrPanic(d.TryDiv(d1))
}

// TryMul multiply d by d1
func (d *Decimal) TryMul(d1 *Decimal) (*Decimal, error) {
	var r Decimal
	_, err := decimal128Context.Mul(&r.Decimal, &d.Decimal, &d1.Decimal)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// Mul multiply d by d1, panics if there is an error.
// Use TryMul to get the error.
func (d *Decimal) Mul(d1 *Decimal) *Decimal {
	return getOrPanic(d.TryMul(d1))
}

// TryAdd adds d and d1
func (d *Decimal) TryAdd(d1 *Decimal) (*Decimal, error) {
	var r Decimal
	_, err := decimal128Context.Add(&r.Decimal, &d.Decimal, &d1.Decimal)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// Add adds d and d1, panics if there is an error.
// Use TryAdd to get the error.
func (d *Decimal) Add(d1 *Decimal) *Decimal {
	return getOrPanic(d.TryAdd(d1))
}

// TrySub subtracts d1 from d
func (d *Decimal) TrySub(d1 *Decimal) (*Decimal, error) {
	var r Decimal
	_, err := decimal128Context.Sub(&r.Decimal, &d.Decimal, &d1.Decimal)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// Sub subtracts d1 from d, panics if there is an error.
// Use TrySub to get the error.
func (d *Decimal) Sub(d1 *Decimal) *Decimal {
	return getOrPanic(d.TrySub(d1))
}

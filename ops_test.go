package decimal_test

import (
	"testing"

	"github.com/fardream/decimal"
)

func TestDecimal_TryMul(t *testing.T) {
	d1 := decimal.NewFromUint64(50)
	d2 := decimal.NewFromUint64(60)

	d3, err := d1.TryMul(d2)
	t.Logf("%s", d3.String())
	if err != nil {
		t.Fatalf("failed to multiply: %v", err)
	}
	if d3.Uint64() != 3000 {
		t.Fatalf("result %d is not 3,000", d3.Uint64())
	}
}

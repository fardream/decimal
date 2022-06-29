package decimal_test

import (
	"encoding/json"
	"testing"

	"github.com/cockroachdb/apd/v3"
	"github.com/fardream/decimal"
)

func TestDecimal_UnmarshalJSON(t *testing.T) {
	d125_345 := apd.New(125345, -3)
	strInput := "\"125.345\""
	var d decimal.Decimal
	if err := json.Unmarshal([]byte(strInput), &d); err != nil {
		t.Fatalf("failed to parse string: %s", strInput)
	}

	if d125_345.Cmp(&d.Decimal) != 0 {
		t.Fatalf("%#v/%s is not 125.345", d, d.String())
	}
	intInput := "125.345"
	if err := json.Unmarshal([]byte(intInput), &d); err != nil {
		t.Fatalf("failed to parse int input: %s", intInput)
	}
	if d125_345.Cmp(&d.Decimal) != 0 {
		t.Fatalf("%#v/%s is not 125.345", d, d.String())
	}
}

const (
	smaller = 1<<62 + 50
	bigger  = (1 << 63) + 50
)

func TestDecimal_SetUint64(t *testing.T) {
	d1 := decimal.NewFromUint64(smaller)
	d2 := decimal.NewFromUint64(bigger)
	d1str, _ := decimal.NewFromString("4611686018427387954")
	d2str, _ := decimal.NewFromString("9223372036854775858")
	if !d1.Equal(d1str) {
		t.Fatalf("smaller: %s is not equal to %s", d1.String(), d1str.String())
	}
	if !d2.Equal(d2str) {
		t.Fatalf("bigger: %s is not equal to %s", d2.String(), d2str.String())
	}
	d1uint64, _ := d1.TryUint64()
	if d1uint64 != smaller {
		t.Fatalf("smaller is not equal")
	}
	d2uint64, _ := d2.TryUint64()
	if d2uint64 != bigger {
		t.Fatalf("bigger is not equal")
	}
}

func TestDecimal_Uint64(t *testing.T) {
	d1 := decimal.NewFromUint64(bigger)
	k, _ := d1.TryUint64()
	var b uint64 = bigger
	if k != b {
		t.Fatalf("%s is not equal to %v", d1.String(), b)
	}
}

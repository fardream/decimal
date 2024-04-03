package decimal

func (d *Decimal) Equal(d1 *Decimal) bool {
	return d.Decimal.Cmp(&d1.Decimal) == 0
}

func (d *Decimal) GreaterThan(d1 *Decimal) bool {
	return d.Decimal.Cmp(&d1.Decimal) > 0
}

func (d *Decimal) LessThan(d1 *Decimal) bool {
	return d.Decimal.Cmp(&d1.Decimal) < 0
}

package vat

import "fmt"

func Vat(amount float64) string {
	return fmt.Sprintf("\nAmount: %.2f\nVat: %.2f", amount, amount*0.12)
}

package str_utils

import "testing"

func TestHumanizeString(t *testing.T) {
	for in, want := range map[string]string{
		"OrderItem":   "Order Item",
		"CNNName":     "CNN Name",
		"order_item":  "Order_item",
		"name":        "Name",
		"já existe":   "Já Existe",
		"ID":          "ID",
		"createdAtBy": "Created At By",
	} {
		if got := HumanizeString(in); got != want {
			t.Errorf("HumanizeString(%q) = %q, want %q", in, got, want)
		}
	}
}

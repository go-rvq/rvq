package models

// Brazilian document validation helpers (CPF/CNPJ). Pure functions, unit tested
// without a DB.

// OnlyDigits returns only the decimal digits of s.
func OnlyDigits(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			out = append(out, s[i])
		}
	}
	return string(out)
}

func allEqual(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

// ValidCPF reports whether doc (with or without punctuation) is a valid CPF.
func ValidCPF(doc string) bool {
	d := OnlyDigits(doc)
	if len(d) != 11 || allEqual(d) {
		return false
	}
	for _, n := range []int{9, 10} {
		sum := 0
		for i := 0; i < n; i++ {
			sum += int(d[i]-'0') * (n + 1 - i)
		}
		dv := (sum * 10) % 11 % 10
		if dv != int(d[n]-'0') {
			return false
		}
	}
	return true
}

// ValidCNPJ reports whether doc (with or without punctuation) is a valid CNPJ.
func ValidCNPJ(doc string) bool {
	d := OnlyDigits(doc)
	if len(d) != 14 || allEqual(d) {
		return false
	}
	weights := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	for _, n := range []int{12, 13} {
		sum := 0
		for i := 0; i < n; i++ {
			sum += int(d[i]-'0') * weights[len(weights)-n+i]
		}
		dv := sum % 11
		if dv < 2 {
			dv = 0
		} else {
			dv = 11 - dv
		}
		if dv != int(d[n]-'0') {
			return false
		}
	}
	return true
}

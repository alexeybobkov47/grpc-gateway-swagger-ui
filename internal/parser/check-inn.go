package parser

import "strconv"

// checkINN - Checking INN for symbols and number of symbols.
func checkINN(inn string) bool {
	const lengthINN = 10
	if len(inn) != lengthINN {
		return false
	}
	for _, n := range inn {
		i, err := strconv.Atoi(string(n))
		if err != nil || i < 0 || i > 9 {
			return false
		}
	}
	return true
}

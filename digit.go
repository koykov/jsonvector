package jsonvector

// Check if given byte is a part of the number.
func isDigit(c byte) bool {
	return (c >= '0' && c <= '9') || c == '-' || c == '+' || c == 'e' || c == 'E'
}

// Check if given is a part of the number, including dot.
func isDigitDot(c byte) bool {
	return isDigit(c) || c == '.'
}

func isDigitTable(c byte) bool {
	_ = digitTable[255]
	return digitTable[c]
}

func isDigitDotTable(c byte) bool {
	_ = digitDotTable[255]
	return digitDotTable[c]
}

var (
	digitTable    [256]bool
	digitDotTable [256]bool
)

func init() {
	digitTable['0'] = true
	digitTable['1'] = true
	digitTable['2'] = true
	digitTable['3'] = true
	digitTable['4'] = true
	digitTable['5'] = true
	digitTable['6'] = true
	digitTable['7'] = true
	digitTable['8'] = true
	digitTable['9'] = true
	digitTable['-'] = true
	digitTable['+'] = true
	digitTable['e'] = true
	digitTable['E'] = true

	digitDotTable['0'] = true
	digitDotTable['1'] = true
	digitDotTable['2'] = true
	digitDotTable['3'] = true
	digitDotTable['4'] = true
	digitDotTable['5'] = true
	digitDotTable['6'] = true
	digitDotTable['7'] = true
	digitDotTable['8'] = true
	digitDotTable['9'] = true
	digitDotTable['-'] = true
	digitDotTable['+'] = true
	digitDotTable['e'] = true
	digitDotTable['E'] = true
	digitDotTable['.'] = true
}

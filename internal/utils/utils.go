package utils

func IsAlphabet(ch rune) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
}

func IsNum(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func IsValidSpecial(ch rune) bool {
	switch ch {
	case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
		return true
	default:
		return false

	}
}

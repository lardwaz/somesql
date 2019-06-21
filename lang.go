package somesql

const (
	// LangEN represents the english language
	LangEN string = "en"
	// LangFR represents the french language
	LangFR string = "fr"
)

// IsLangValid checks if lang is valid
func IsLangValid(lang string) bool {
	switch lang {
	case LangEN, LangFR:
		return true
	default:
		return false
	}
}

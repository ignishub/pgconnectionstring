package pgconnectionstring

import (
	"errors"
	"unicode"
)

var (
	ErrUnterminatedQuote              = errors.New(`unterminated quoted string literal in connection string`)
	ErrMissingCharacterAfterBackslash = errors.New(`missing character after backslash`)
	ErrMissingEqualSign               = errors.New(`missing = sign after parametr name`)
)

// quote qutes string if it has spaces
func quote(s string) string {
	for _, v := range []rune(s) {
		if unicode.IsSpace(v) {
			return "'" + s + "'"
		}
	}
	return s
}

// ToConnectionString returns values formated as connection string.
func ToConnectionString(v map[string]string) string {
	s := ""
	for key, value := range v {
		if key == "" || value == "" {
			continue
		}
		s = s + quote(key) + "=" + quote(value) + " "
	}
	return s
}

// scanner implements a tokenizer for libpq-style option strings.
type scanner struct {
	s []rune
	i int
}

// newScanner returns a new scanner initialized with the option string s.
func newScanner(s string) *scanner {
	return &scanner{[]rune(s), 0}
}

// Next returns the next rune.
// It returns 0, false if the end of the text has been reached.
func (s *scanner) Next() (rune, bool) {
	if s.i >= len(s.s) {
		return 0, false
	}
	r := s.s[s.i]
	s.i++
	return r, true
}

// SkipSpaces returns the next non-whitespace rune.
// It returns 0, false if the end of the text has been reached.
func (s *scanner) SkipSpaces() (rune, bool) {
	r, ok := s.Next()
	for unicode.IsSpace(r) && ok {
		r, ok = s.Next()
	}
	return r, ok
}

// Parse parses the options from name and adds them to the values.
//
// The parsing code is based on conninfo_parse from libpq's fe-connect.c
func Parse(name string) (map[string]string, error) {
	s := newScanner(name)
	values := make(map[string]string)
	for {
		var (
			keyRunes, valRunes []rune
			r                  rune
			ok                 bool
		)

		if r, ok = s.SkipSpaces(); !ok {
			break
		}

		// Scan the key
		for !unicode.IsSpace(r) && r != '=' {
			keyRunes = append(keyRunes, r)
			if r, ok = s.Next(); !ok {
				break
			}
		}

		// Skip any whitespace if we're not at the = yet
		if r != '=' {
			r, ok = s.SkipSpaces()
		}

		// The current character should be =
		if r != '=' || !ok {
			return nil, ErrMissingEqualSign
		}

		// Skip any whitespace after the =
		if r, ok = s.SkipSpaces(); !ok {
			// If we reach the end here, the last value is just an empty string as per libpq.
			values[string(keyRunes)] = ""
			break
		}

		if r != '\'' {
			for !unicode.IsSpace(r) {
				if r == '\\' {
					if r, ok = s.Next(); !ok {
						return nil, ErrMissingCharacterAfterBackslash
					}
				}
				valRunes = append(valRunes, r)

				if r, ok = s.Next(); !ok {
					break
				}
			}
		} else {
		quote:
			for {
				if r, ok = s.Next(); !ok {
					return nil, ErrUnterminatedQuote
				}
				switch r {
				case '\'':
					break quote
				case '\\':
					r, _ = s.Next()
					fallthrough
				default:
					valRunes = append(valRunes, r)
				}
			}
		}

		values[string(keyRunes)] = string(valRunes)
	}

	return values, nil
}

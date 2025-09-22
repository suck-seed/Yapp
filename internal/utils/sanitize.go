package utils

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/microcosm-cc/bluemonday"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/text/unicode/norm"
)

var usernameRegex = regexp.MustCompile(`^[a-z0-9_.-]{3,32}$`)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
var hexRegex = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$`)

// NAME SECTION
func SanatizeUsername(s string) (string, error) {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	if !usernameRegex.MatchString(s) {
		return "", ErrorInvalidUsername
	}
	return s, nil
}

func SanatizeHallname(s string) (string, error) {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	if !usernameRegex.MatchString(s) {
		return "", ErrorInvalidHallName
	}
	return s, nil
}

func SanatizeDisplayName(ptr *string) (*string, error) {
	if ptr == nil {
		return nil, nil
	}
	s := strings.TrimSpace(*ptr)
	// normalize unicode; collapse spaces
	s = norm.NFKC.String(s)
	s = strings.Join(strings.Fields(s), " ")
	if s == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(s) > 32 || utf8.RuneCountInString(s) < 3 {
		return nil, ErrorInvalidDisplayName
	}
	return &s, nil
}

// EMAIL SECTION
func SanatizeEmail(s string) (string, error) {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	// very light check; rely on validator.v10 too
	if len(s) < 6 || len(s) > 254 || !strings.Contains(s, "@") || !emailRegex.MatchString(s) {
		return "", ErrorInvalidEmail
	}
	return s, nil
}

func SanatizePhoneE164(ptr *string) (*string, error) {
	if ptr == nil {
		return nil, nil
	}
	s := strings.TrimSpace(*ptr)
	if s == "" {
		return nil, nil
	}
	// If you use libphonenumber:
	// num, err := phonenumbers.Parse(s, "NP") // or your default region
	// if err != nil || !phonenumbers.IsValidNumber(num) { return nil, ErrInvalidPhone }
	// e164 := phonenumbers.Format(num, phonenumbers.E164)
	// return &e164, nil
	// If not using lib yet, minimally keep digits/+ and do a length check:
	s = keepPlusDigits(s)
	if len(s) < 7 || len(s) > 20 {
		return nil, ErrorInvalidPhoneNumber
	}
	return &s, nil
}

func keepPlusDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '+' || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// PASSWORD SECTION
const minEntropyBits = 60.0 // ~good baseline for online attacks; use 70–80 for higher risk

func SanatizePasswordPolicy(raw string) (string, error) {
	// Do NOT silently modify. Reject confusing whitespace at edges.
	if strings.TrimSpace(raw) != raw {
		return "", ErrorPasswordWhiteSpace
	}
	if err := passwordvalidator.Validate(raw, minEntropyBits); err != nil {
		return "", ErrorInvalidPassword
	}
	return raw, nil
}

// COLOR SECTION
func SanatizeColorFormat(colorHex *string) (*string, error) {

	if colorHex == nil {
		return nil, nil
	}

	s := strings.TrimSpace(*colorHex)
	if s == "" {
		return nil, nil
	}

	if !hexRegex.MatchString(s) {
		return nil, ErrorInvalidBannerColor
	}

	return &s, nil
}

// TEXT SECTION
func SanatizeText(text *string) (*string, error) {
	if text == nil {
		return nil, nil
	}

	s := strings.TrimSpace(*text)
	if s == "" {
		return nil, nil
	}

	// xss injection prevention
	p := bluemonday.UGCPolicy()
	s = p.Sanitize(s)

	return &s, nil
}

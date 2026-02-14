package customer

import "strings"

const minPasswordLength = 8

// Email is a value object representing a validated, normalized email address.
type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	if !isValidEmail(normalized) {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: normalized}, nil
}

func (e Email) String() string    { return e.value }
func (e Email) Equals(other Email) bool { return e.value == other.value }

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return false
	}
	return parts[0] != "" && parts[1] != ""
}

// Password is a value object representing a raw plaintext password (input boundary only).
type Password struct {
	value string
}

func NewPassword(raw string) (Password, error) {
	if len(raw) < minPasswordLength {
		return Password{}, ErrInvalidPassword
	}
	return Password{value: raw}, nil
}

func (p Password) String() string { return p.value }

// PhoneNumber is a value object representing a phone number.
type PhoneNumber struct {
	value string
}

func NewPhoneNumber(raw string) (PhoneNumber, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return PhoneNumber{}, ErrInvalidPhoneNumber
	}
	return PhoneNumber{value: trimmed}, nil
}

func (p PhoneNumber) String() string { return p.value }

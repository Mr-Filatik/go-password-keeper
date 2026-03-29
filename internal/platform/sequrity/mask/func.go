package mask

import "strings"

// Password completely masks the password value by replacing it with the * symbol.
//
// Examples:
//   - "password" - "*";
//   - "" - "*".
func Password(_ string) string {
	return "*"
}

// Email masks every character of the email (except the first and last characters),
// but does not affect the domain. In case of incorrect data, it will be returned unchanged.
//
// Examples:
//   - "user@email.com" - "u**r@email.com";
//   - "u5@email.com" - "u5@email.com";
//   - "incorrect.email.com" - "incorrect.email.com".
func Email(email string) string {
	const (
		partCount = 2 // Number of address parts when divided by "@": user and domain
		prefix    = 1 // How many characters to leave at the beginning of the address before hiding
		suffix    = 1 // How many characters to leave at the end of the address after hiding
	)

	parts := strings.Split(email, "@")
	if len(parts) != partCount {
		return email
	}

	user, domain := parts[0], parts[1]

	runes := []rune(user)
	userLen := len(runes)

	if userLen <= prefix+suffix {
		return email
	}

	maskedUser := user[:prefix] +
		strings.Repeat("*", userLen-(prefix+suffix)) +
		user[len(user)-suffix:]

	return maskedUser + "@" + domain
}

// Phone masks some digits of the number, leaving the last 2 digits exactly.
//
// Examples:
//   - "+78005553535" - "+7800*****35";
//   - "+7858005553535" - "+785800*****35";
//   - "8005553535" - "800*****35".
func Phone(phone string) string {
	const (
		save = 2 // How many digits from the end of the number to keep
		mask = 5 // How many characters in a row to mask
	)

	if len(phone) < save {
		return phone
	}

	if len(phone) < save+mask {
		masked := len(phone) - save

		return strings.Repeat("*", masked) + phone[len(phone)-save:]
	}

	return phone[:len(phone)-(save+mask)] +
		strings.Repeat("*", mask) +
		phone[len(phone)-save:]
}

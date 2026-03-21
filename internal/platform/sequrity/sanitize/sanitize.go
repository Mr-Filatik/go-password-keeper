package sanitize

type Type string

const mask = "[MASKED]"

const (
	TypeFull   Type = "[MASKED]"
	TypeStart  Type = "[MASKED]data"
	TypeMiddle Type = "data[MASKED]data"
	TypeEnd    Type = "data[MASKED]"
)

// M - Mask
func StringM(data string, typ Type, count int) string {
	switch typ {
	case TypeFull:
		return mask

	case TypeStart:
		if count > len(data) {
			return mask
		}

		return mask + data[count:]

	case TypeMiddle:
		if count > len(data) {
			return mask
		}

		index := (len(data) - count) / 2

		return data[:index] + mask + data[index+count:]

	case TypeEnd:
		if count > len(data) {
			return mask
		}

		return data[:count] + mask

	default:
		return data
	}
}

// MS - Mask Save
func StringMS(data string, typ Type, save int) string {
	switch typ {
	case TypeFull:
		return mask

	case TypeStart:
		if save > len(data) {
			return mask
		}

		return mask + data[len(data)-save:]

	case TypeMiddle:
		if save*2 > len(data) {
			return mask
		}

		return data[:save] + mask + data[len(data)-save:]

	case TypeEnd:
		if save > len(data) {
			return mask
		}

		return data[:save] + mask

	default:
		return data
	}
}

// MSD - Mask Save Different
func StringMSD(data string, typ Type, save, count int) string {
	switch typ {
	case TypeFull:
		return mask

	case TypeStart:
		if save > len(data) {
			return mask
		}

		return mask + data[len(data)-save:]

	case TypeMiddle:
		if save*2 > len(data) {
			return mask
		}

		return data[:save] + mask + data[len(data)-save:]

	case TypeEnd:
		if save > len(data) {
			return mask
		}

		return data[:save] + mask

	default:
		return data
	}
}

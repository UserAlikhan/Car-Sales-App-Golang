package utils

func FindStringKeyStringValue(array []*string, value *string) bool {
	for _, element := range array {
		if element == value {
			return true
		}
	}

	return false
}

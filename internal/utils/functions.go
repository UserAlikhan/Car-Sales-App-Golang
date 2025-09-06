package utils

func Find(array []string, value string) bool {
	for _, element := range array {
		if element == value {
			return true
		}
	}

	return false
}

func ConvertStringPointerArrayToStringArray(pointerArray []*string) []string {
	convertedArray := make([]string, len(pointerArray))

	for i, v := range pointerArray {
		if v != nil {
			convertedArray[i] = *v
		}
	}

	return convertedArray
}

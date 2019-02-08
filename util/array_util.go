package util

func Contains(array []string, e string) bool {
	for _, element := range array {
		if element == e {
			return true
		}
	}
	return false
}

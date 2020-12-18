package utils

func CopyStrPtr(input *string) (output *string) {
	if input != nil {
		newStr := *input
		output = &newStr
	}
	return
}

package ellipsis

func Ellipsis(text string, length int) string {
	return EllipsisSuffix(text, length, "...")
}

func EllipsisSuffix(text string, length int, suffix string) string {
	if len(text) <= length {
		return text
	}
	return text[:length-len(suffix)] + suffix
}

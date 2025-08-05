package ellipsis

func Ellipsis(text string, length int) string {
	if len(text) <= length {
		return text
	}
	return text[:length-3] + "..."
}

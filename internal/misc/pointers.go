package misc

func StrPtrToStr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func StrSlicePtrToStrSlice(s *[]string) []string {
	if s != nil {
		return *s
	}
	return []string{}
}

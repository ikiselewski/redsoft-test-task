package misc

func StrPtrToStr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

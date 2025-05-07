package lib

func StringOrDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

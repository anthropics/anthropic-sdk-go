package paramutil

func AddrOrDefault[T interface {
	~string
	Default() T
}](v T) *string {
	if v == "" {
		v = v.Default()
	}
	s := string(v)
	return &s
}

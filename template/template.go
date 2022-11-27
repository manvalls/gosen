package template

type Template struct {
}

func (t *Template) MarshalText() (text []byte, err error) {
	return []byte{}, nil
}

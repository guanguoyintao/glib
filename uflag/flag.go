package uflag

type ArrayStringFlags []string

func (i *ArrayStringFlags) String() string {
	return "my string representation"
}

func (i *ArrayStringFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

package uid

type UID interface {
	NewString() string
	NewInt() int64
}

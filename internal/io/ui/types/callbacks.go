package types

type KeyCallback struct {
	Key      string
	Label    string
	Callback func() error
}

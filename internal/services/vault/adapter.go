package vault

type Adapter interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
	List() ([]string, error)
}

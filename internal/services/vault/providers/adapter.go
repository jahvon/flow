package providers

//go:generate mockgen -destination=mocks/mock_adapter.go -package=mocks github.com/jahvon/flow/internal/services/vault/providers Adapter
type Adapter interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
	List() (map[string]string, error)
}

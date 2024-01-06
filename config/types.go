package config

type Entity interface {
	YAML() (string, error)
	JSON(formatted bool) (string, error)
	Markdown() string
}

type CollectionItem struct {
	Header      string
	SubHeader   string
	Description string
	Tags        Tags
}

type Collection interface {
	Items() []CollectionItem
	YAML() (string, error)
	JSON(formatted bool) (string, error)
	Singular() string
	Plural() string
}

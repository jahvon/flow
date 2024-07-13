package common

import (
	"fmt"
	"slices"
	"strings"
)

//go:generate go run github.com/atombender/go-jsonschema@v0.16.0 -et --only-models -p common -o common.gen.go schema.yaml

func (a Aliases) HasAlias(alias string) bool {
	if alias == "" {
		return true
	}
	return slices.Contains(a, alias)
}

func (t Tags) String() string {
	return strings.Join(t, ", ")
}

func (t Tags) PreviewString() string {
	if len(t) == 0 {
		return ""
	}
	count := len(t)
	if count <= 3 {
		return strings.Join(t, ", ")
	}

	return strings.Join(t[:3], ", ") + fmt.Sprintf(" (+%d)", count-3)
}

func (t Tags) HasAnyTag(tags Tags) bool {
	if len(tags) == 0 {
		return true
	}
	for _, tag := range tags {
		if t.HasTag(tag) {
			return true
		}
	}
	return false
}

func (t Tags) HasTag(tag string) bool {
	if tag == "" {
		return true
	}
	return slices.Contains(t, tag)
}

func (v Visibility) String() string {
	return string(v)
}

func (v Visibility) NewPointer() *Visibility {
	return &v
}

func (v Visibility) IsPublic() bool {
	return v == VisibilityPublic
}

func (v Visibility) IsPrivate() bool {
	return v == VisibilityPrivate
}

func (v Visibility) IsInternal() bool {
	return v == VisibilityInternal
}

func (v Visibility) IsHidden() bool {
	return v == VisibilityHidden
}

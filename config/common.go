package config

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

const (
	ActivateGroupID   = "activate"
	DeactivateGroupID = "deactivate"
	LaunchGroupID     = "launch"
)

type Verb string

var (
	ValidVerbToGroupID = map[string]string{
		"exec":      ActivateGroupID,
		"run":       ActivateGroupID,
		"start":     ActivateGroupID,
		"install":   ActivateGroupID,
		"setup":     ActivateGroupID,
		"delete":    DeactivateGroupID,
		"remove":    DeactivateGroupID,
		"uninstall": DeactivateGroupID,
		"teardown":  DeactivateGroupID,
		"destroy":   DeactivateGroupID,
		"open":      LaunchGroupID,
		"launch":    LaunchGroupID,
	}

	ValidVerbs = lo.Keys(ValidVerbToGroupID)
)

func (v Verb) String() string {
	return string(v)
}

func (v Verb) Validate() error {
	if !lo.Contains(ValidVerbs, v.String()) {
		return fmt.Errorf("invalid executable verb %s", v)
	}
	return nil
}

func (v Verb) Equals(other Verb) bool {
	return ValidVerbToGroupID[v.String()] == ValidVerbToGroupID[other.String()]
}

const (
	VisibilityPublic   VisibilityType = "public"
	VisibilityPrivate  VisibilityType = "private"
	VisibilityInternal VisibilityType = "internal"
	VisibilityHidden   VisibilityType = "hidden"
)

// From highest the visible [0] to the lowest visible [n-1].
var visibilityByLevel = []VisibilityType{
	VisibilityPublic,
	VisibilityPrivate,
	VisibilityInternal,
	VisibilityHidden,
}

type VisibilityType string

func (v VisibilityType) IsPublic() bool {
	return v == VisibilityPublic
}

func (v VisibilityType) IsPrivate() bool {
	return v == VisibilityPrivate
}

func (v VisibilityType) IsInternal() bool {
	return v == VisibilityInternal
}

func (v VisibilityType) IsHidden() bool {
	return v == VisibilityHidden
}

type Tags []string

func (t Tags) HasAnyTag(tags Tags) bool {
	if len(tags) == 0 {
		return true
	}
	return lo.Some(t, tags)
}

func (t Tags) HasTag(tag string) bool {
	if tag == "" {
		return true
	}
	return lo.Contains(t, tag)
}

type Aliases []string

func (a Aliases) HasAlias(alias string) bool {
	if alias == "" {
		return true
	}
	for _, a := range a {
		if a == alias {
			return true
		}
	}
	return false
}

type Ref string

func NewRef(id string, verb Verb) Ref {
	if id == "" || verb == "" {
		return ""
	}
	return Ref(fmt.Sprintf("%s %s", verb, id))
}

func (r Ref) String() string {
	return string(r)
}

func (r Ref) Validate() error {
	str := strings.TrimSpace(string(r))
	refParts := strings.Split(str, " ")
	if len(refParts) != 2 {
		return fmt.Errorf("invalid executable ref %s", str)
	}
	verb := Verb(refParts[0])
	if err := verb.Validate(); err != nil {
		return err
	}
	id := refParts[1]
	ws, _, name := parseExecutableID(id)
	if ws == "" || name == "" {
		return fmt.Errorf("invalid executable id %s", id)
	}
	return nil
}

func (r Ref) GetVerb() Verb {
	refParts := strings.Split(string(r), " ")
	return Verb(refParts[0])
}

func (r Ref) GetID() string {
	refParts := strings.Split(string(r), " ")
	return refParts[1]
}

func (r Ref) Equals(other Ref) bool {
	rVerb := r.GetVerb()
	oVerb := other.GetVerb()
	if !rVerb.Equals(oVerb) {
		return false
	}

	return r.GetID() == other.GetID()
}

type LogMode string

const (
	NoLogMode         LogMode = "none"
	StructuredLogMode LogMode = "structured"
	RawLogMode        LogMode = "raw"
)

package config

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"
)

const (
	ActivateGroupID   = "activate"
	DeactivateGroupID = "deactivate"
	UpdateGroupID     = "update"
	ManageGroupID     = "manage"
	LaunchGroupID     = "launch"
	CreationGroupID   = "creation"
)

// +docsgen:verb
// Keywords that describe the action an executable performs.
// While executables are configured with a single verb, the verb can be aliased to related verbs.
// For example, the `exec` verb can replaced with "run" or "start" when referencing an executable.
// This allows users to use the verb that best describes the action they are performing.
//
// **Activation verbs**: `exec`, `run`, `start`, `install`, `setup`, `release`, `deploy`, `apply`
// **Deactivation verbs**: `delete`, `remove`, `uninstall`, `destroy`, `undeploy`
// **Update verbs**: `update`, `upgrade`, `refresh`, `reload`
// **Management verbs**: `manage`, `configure`, `monitor`, `edit`
// **Launch verbs**: `open`, `launch`, `show`, `view`, `render`
// **Creation verbs**: `generate`, `add`, `new`, `build`, `transform`
type Verb string

var (
	ValidVerbToGroupID = map[string]string{
		"exec":      ActivateGroupID,
		"run":       ActivateGroupID,
		"start":     ActivateGroupID,
		"install":   ActivateGroupID,
		"setup":     ActivateGroupID,
		"release":   ActivateGroupID,
		"deploy":    ActivateGroupID,
		"apply":     ActivateGroupID,
		"delete":    DeactivateGroupID,
		"remove":    DeactivateGroupID,
		"uninstall": DeactivateGroupID,
		"destroy":   DeactivateGroupID,
		"undeploy":  DeactivateGroupID,
		"update":    UpdateGroupID,
		"upgrade":   UpdateGroupID,
		"refresh":   UpdateGroupID,
		"reload":    UpdateGroupID,
		"manage":    ManageGroupID,
		"configure": ManageGroupID,
		"monitor":   ManageGroupID,
		"edit":      ManageGroupID,
		"open":      LaunchGroupID,
		"launch":    LaunchGroupID,
		"show":      LaunchGroupID,
		"view":      LaunchGroupID,
		"render":    LaunchGroupID,
		"generate":  CreationGroupID,
		"add":       CreationGroupID,
		"new":       CreationGroupID,
		"transform": CreationGroupID,
		"build":     CreationGroupID,
	}
)

func ValidVerbs() []string {
	return lo.Keys(ValidVerbToGroupID)
}

func SortedValidVerbs() []string {
	verbs := ValidVerbs()
	slices.SortFunc(verbs, strings.Compare)
	return verbs
}

func (v Verb) String() string {
	return string(v)
}

func (v Verb) Validate() error {
	if !lo.Contains(ValidVerbs(), v.String()) {
		return fmt.Errorf("invalid executable verb %s", v)
	}
	return nil
}

func (v Verb) Equals(other Verb) bool {
	return ValidVerbToGroupID[v.String()] == ValidVerbToGroupID[other.String()]
}

func RelatedVerbs(verb Verb) []Verb {
	verbs := make([]Verb, 0)
	for _, v := range ValidVerbs() {
		if ValidVerbToGroupID[v] == ValidVerbToGroupID[verb.String()] {
			verbs = append(verbs, Verb(v))
		}
	}
	return lo.Uniq(verbs)
}

const (
	// VisibilityPublic is executable and visible across all workspaces.
	VisibilityPublic VisibilityType = "public"
	// VisibilityPrivate is executable and visible only within it's own workspace.
	VisibilityPrivate VisibilityType = "private"
	// VisibilityInternal is not visible but can be executed within a workspace.
	VisibilityInternal VisibilityType = "internal"
	// VisibilityHidden is not executable or visible.
	VisibilityHidden VisibilityType = "hidden"
)

// From highest the visible [0] to the lowest visible [n-1].
var visibilityByLevel = []VisibilityType{
	VisibilityPublic,
	VisibilityPrivate,
	VisibilityInternal,
	VisibilityHidden,
}

type VisibilityType string

func (v VisibilityType) String() string {
	return string(v)
}

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

// +docsgen:tags
// A list of tags.
// Tags can be used with list commands to filter returned data.
type Tags []string

func (t Tags) String() string {
	return strings.Join(t, ", ")
}

func (t Tags) ContextString() string {
	if len(t) == 0 {
		return ""
	}
	var str string
	lo.ForEach(t, func(tag string, i int) {
		str += fmt.Sprintf("[%s] ", tag)
	})
	return str
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

// +docsgen:ref
// A reference to an executable.
// The format is `<verb> <workspace>/<namespace>:<executable name>`.
// For example, `exec ws/ns:my-flow`.
//
// The workspace and namespace are optional.
// If the workspace is not specified, the current workspace will be used.
// If the namespace is not specified, the current namespace will be used.
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
	ws, _, name := ParseExecutableID(id)
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

package executable

import (
	"fmt"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

const (
	ActivateGroupID   = "activate"
	DeactivateGroupID = "deactivate"
	UpdateGroupID     = "update"
	ManageGroupID     = "manage"
	LaunchGroupID     = "launch"
	CreationGroupID   = "creation"
)

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
	return maps.Keys(ValidVerbToGroupID)
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
	if !slices.Contains(ValidVerbs(), v.String()) {
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
	return slices.Compact(verbs)
}

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
	if len(refParts) != 2 {
		// TODO: return or log error
		return ""
	}
	return refParts[1]
}

func (r Ref) GetNamespace() string {
	id := r.GetID()
	_, ns, _ := ParseExecutableID(id)
	return ns
}

func (r Ref) GetWorkspace() string {
	id := r.GetID()
	ws, _, _ := ParseExecutableID(id)
	return ws
}

func (r Ref) Equals(other Ref) bool {
	rVerb := r.GetVerb()
	oVerb := other.GetVerb()
	if !rVerb.Equals(oVerb) {
		return false
	}

	return r.GetID() == other.GetID()
}

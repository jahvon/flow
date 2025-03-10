package executable

import (
	"fmt"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

const (
	ActivateGroupID     = "activate"
	DeactivateGroupID   = "deactivate"
	ExecutionGroupID    = "execute"
	TerminationGroupID  = "terminate"
	MonitoringGroupID   = "monitor"
	RestartGroupID      = "restart"
	InstallGroupID      = "install"
	BuildGroupID        = "build"
	UninstallGroupID    = "uninstall"
	PublishGroupID      = "publish"
	DistributionGroupID = "distribute"
	TestGroupID         = "test"
	AnalyzeGroupID      = "analyze"
	LaunchGroupID       = "launch"
	CreationGroupID     = "create"
	SetGroupID          = "set"
	DestructionGroupID  = "destroy"
	UnsetGroupID        = "unset"
	CleanupGroupID      = "cleanup"
	RetrievalGroupID    = "retrieve"
	UpdateGroupID       = "update"
	EditGroupID         = "edit"
	DebugGroupID        = "debug"
)

var (
	ValidVerbToGroupID = map[Verb]string{
		// Activation verbs
		VerbActivate: ActivateGroupID,
		VerbEnable:   ActivateGroupID,
		VerbStart:    ActivateGroupID,
		VerbTrigger:  ActivateGroupID,

		// Execution verbs
		VerbExec:    ExecutionGroupID,
		VerbRun:     ExecutionGroupID,
		VerbExecute: ExecutionGroupID,

		// Deactivation verbs
		VerbDeactivate: DeactivateGroupID,
		VerbDisable:    DeactivateGroupID,
		VerbStop:       DeactivateGroupID,
		VerbPause:      DeactivateGroupID,

		// Termination verbs
		VerbKill:      TerminationGroupID,
		VerbTerminate: TerminationGroupID,
		VerbAbort:     TerminationGroupID,

		// Monitoring verbs
		VerbWatch:   MonitoringGroupID,
		VerbMonitor: MonitoringGroupID,
		VerbTrack:   MonitoringGroupID,

		// Restart verbs
		VerbRestart: RestartGroupID,
		VerbReboot:  RestartGroupID,
		VerbReload:  RestartGroupID,
		VerbRefresh: RestartGroupID,

		// Installation verbs
		VerbInstall: InstallGroupID,
		VerbSetup:   InstallGroupID,
		VerbDeploy:  InstallGroupID,

		// Build verbs
		VerbBuild:   BuildGroupID,
		VerbPackage: BuildGroupID,
		VerbBundle:  BuildGroupID,
		VerbCompile: BuildGroupID,

		// Uninstallation verbs
		VerbUninstall: UninstallGroupID,
		VerbTeardown:  UninstallGroupID,
		VerbUndeploy:  UninstallGroupID,

		// Update verbs
		VerbUpdate:  UpdateGroupID,
		VerbUpgrade: UpdateGroupID,
		VerbPatch:   UpdateGroupID,

		// Publish verbs
		VerbPublish: PublishGroupID,
		VerbRelease: PublishGroupID,

		// Distribution verbs
		VerbPush:  DistributionGroupID,
		VerbSend:  DistributionGroupID,
		VerbApply: DistributionGroupID,

		// Test verbs
		VerbTest:     TestGroupID,
		VerbValidate: TestGroupID,
		VerbCheck:    TestGroupID,
		VerbVerify:   TestGroupID,

		// Analysis verbs
		VerbAnalyze: AnalyzeGroupID,
		VerbScan:    AnalyzeGroupID,
		VerbLint:    AnalyzeGroupID,
		VerbInspect: AnalyzeGroupID,

		// Launch verbs
		VerbOpen:   LaunchGroupID,
		VerbLaunch: LaunchGroupID,
		VerbShow:   LaunchGroupID,
		VerbView:   LaunchGroupID,

		// Creation verbs
		VerbCreate:   CreationGroupID,
		VerbGenerate: CreationGroupID,
		VerbAdd:      CreationGroupID,
		VerbNew:      CreationGroupID,
		VerbInit:     CreationGroupID,

		// Edit verbs
		VerbEdit:      EditGroupID,
		VerbModify:    EditGroupID,
		VerbFix:       EditGroupID,
		VerbTransform: EditGroupID,

		// Set verbs
		VerbSet: SetGroupID,

		// Destruction verbs
		VerbRemove:  DestructionGroupID,
		VerbDelete:  DestructionGroupID,
		VerbDestroy: DestructionGroupID,
		VerbErase:   DestructionGroupID,

		// Unset verbs
		VerbUnset: UnsetGroupID,
		VerbReset: UnsetGroupID,

		// Cleanup verbs
		VerbClean: CleanupGroupID,
		VerbClear: CleanupGroupID,
		VerbPurge: CleanupGroupID,
		VerbTidy:  CleanupGroupID,

		// Retrieval verbs
		VerbRetrieve: RetrievalGroupID,
		VerbFetch:    RetrievalGroupID,
		VerbGet:      RetrievalGroupID,
		VerbRequest:  RetrievalGroupID,

		// Debug verbs
		VerbDebug:   DebugGroupID,
		VerbTrace:   DebugGroupID,
		VerbProfile: DebugGroupID,
	}
)

func ValidVerbs() []Verb {
	return maps.Keys(ValidVerbToGroupID)
}

func SortedValidVerbs() []string {
	verbs := make([]string, 0)
	for _, v := range ValidVerbs() {
		verbs = append(verbs, v.String())
	}
	slices.Sort(verbs)
	return verbs
}

func (v Verb) String() string {
	return string(v)
}

func (v Verb) Validate() error {
	if !slices.Contains(ValidVerbs(), v) {
		return fmt.Errorf("invalid executable verb %s", v)
	}
	return nil
}

func (v Verb) Equals(other Verb) bool {
	return ValidVerbToGroupID[v] == ValidVerbToGroupID[other]
}

func RelatedVerbs(verb Verb) []Verb {
	verbs := make([]Verb, 0)
	for _, v := range ValidVerbs() {
		if ValidVerbToGroupID[v] == ValidVerbToGroupID[verb] {
			verbs = append(verbs, v)
		}
	}
	return slices.Compact(verbs)
}

func NewRef(id string, verb Verb) Ref {
	if verb == "" {
		return ""
	}
	if id == "" {
		return Ref(verb.String())
	}
	return Ref(fmt.Sprintf("%s %s", verb, id))
}

func (r Ref) String() string {
	return string(r)
}

func (r Ref) Verb() Verb {
	refParts := strings.Split(string(r), " ")
	return Verb(refParts[0])
}

func (r Ref) ID() string {
	refParts := strings.Split(string(r), " ")
	if len(refParts) == 2 {
		return refParts[1]
	}
	return ""
}

func (r Ref) Namespace() string {
	id := r.ID()
	_, ns, _ := MustParseExecutableID(id)
	return ns
}

func (r Ref) Workspace() string {
	id := r.ID()
	ws, _, _ := MustParseExecutableID(id)
	return ws
}

func (r Ref) Equals(other Ref) bool {
	rVerb := r.Verb()
	oVerb := other.Verb()
	if !rVerb.Equals(oVerb) {
		return false
	}

	return r.ID() == other.ID()
}

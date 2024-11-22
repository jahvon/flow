// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package executable

import "github.com/jahvon/flow/types/common"
import "github.com/jahvon/tuikit/io"
import "time"

type Argument struct {
	// The default value to use if the argument is not provided.
	// If the argument is required and no default is provided, the executable will
	// fail.
	//
	Default string `json:"default,omitempty" yaml:"default,omitempty" mapstructure:"default,omitempty"`

	// The name of the environment variable that will be assigned the value.
	EnvKey string `json:"envKey" yaml:"envKey" mapstructure:"envKey"`

	// The flag to use when setting the argument from the command line.
	// Either `flag` or `pos` must be set, but not both.
	//
	Flag string `json:"flag,omitempty" yaml:"flag,omitempty" mapstructure:"flag,omitempty"`

	// The position of the argument in the command line ArgumentList. Values start at
	// 1.
	// Either `flag` or `pos` must be set, but not both.
	//
	Pos int `json:"pos,omitempty" yaml:"pos,omitempty" mapstructure:"pos,omitempty"`

	// If the argument is required, the executable will fail if the argument is not
	// provided.
	// If the argument is not required, the default value will be used if the argument
	// is not provided.
	//
	Required bool `json:"required,omitempty" yaml:"required,omitempty" mapstructure:"required,omitempty"`

	// The type of the argument. This is used to determine how to parse the value of
	// the argument.
	Type ArgumentType `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`

	// value corresponds to the JSON schema field "value".
	value string `json:"value,omitempty" yaml:"value,omitempty" mapstructure:"value,omitempty"`
}

type ArgumentList []Argument

type ArgumentType string

const ArgumentTypeBool ArgumentType = "bool"
const ArgumentTypeFloat ArgumentType = "float"
const ArgumentTypeInt ArgumentType = "int"
const ArgumentTypeString ArgumentType = "string"

// The directory to execute the command in.
// If unset, the directory of the flow file will be used.
// If set to `f:tmp`, a temporary directory will be created for the process.
// If prefixed with `./`, the path will be relative to the current working
// directory.
// If prefixed with `//`, the path will be relative to the workspace root.
// Environment variables in the path will be expended at runtime.
type Directory string

// Standard executable type. Runs a command/file in a subprocess.
type ExecExecutableType struct {
	// Args corresponds to the JSON schema field "args".
	Args ArgumentList `json:"args,omitempty" yaml:"args,omitempty" mapstructure:"args,omitempty"`

	// The command to execute.
	// Only one of `cmd` or `file` must be set.
	//
	Cmd string `json:"cmd,omitempty" yaml:"cmd,omitempty" mapstructure:"cmd,omitempty"`

	// Dir corresponds to the JSON schema field "dir".
	Dir Directory `json:"dir,omitempty" yaml:"dir,omitempty" mapstructure:"dir,omitempty"`

	// The file to execute.
	// Only one of `cmd` or `file` must be set.
	//
	File string `json:"file,omitempty" yaml:"file,omitempty" mapstructure:"file,omitempty"`

	// logFields corresponds to the JSON schema field "logFields".
	logFields map[string]interface{} `json:"logFields,omitempty" yaml:"logFields,omitempty" mapstructure:"logFields,omitempty"`

	// The log mode to use when running the executable.
	// This can either be `hidden`, `json`, `logfmt` or `text`
	//
	LogMode io.LogMode `json:"logMode,omitempty" yaml:"logMode,omitempty" mapstructure:"logMode,omitempty"`

	// Params corresponds to the JSON schema field "params".
	Params ParameterList `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

// The executable schema defines the structure of an executable in the Flow CLI.
// Executables are the building blocks of workflows and are used to define the
// actions that can be performed in a workspace.
type Executable struct {
	// Aliases corresponds to the JSON schema field "aliases".
	Aliases ExecutableAliases `json:"aliases,omitempty" yaml:"aliases,omitempty" mapstructure:"aliases,omitempty"`

	// A description of the executable.
	// This description is rendered as markdown in the interactive UI.
	//
	Description string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// Exec corresponds to the JSON schema field "exec".
	Exec *ExecExecutableType `json:"exec,omitempty" yaml:"exec,omitempty" mapstructure:"exec,omitempty"`

	// flowFilePath corresponds to the JSON schema field "flowFilePath".
	flowFilePath string `json:"flowFilePath,omitempty" yaml:"flowFilePath,omitempty" mapstructure:"flowFilePath,omitempty"`

	// inheritedDescription corresponds to the JSON schema field
	// "inheritedDescription".
	inheritedDescription string `json:"inheritedDescription,omitempty" yaml:"inheritedDescription,omitempty" mapstructure:"inheritedDescription,omitempty"`

	// Launch corresponds to the JSON schema field "launch".
	Launch *LaunchExecutableType `json:"launch,omitempty" yaml:"launch,omitempty" mapstructure:"launch,omitempty"`

	// The name of the executable.
	//
	// Name is used to reference the executable in the CLI using the format
	// `workspace:namespace/name`.
	// [Verb group + Name] must be unique within the namespace of the workspace.
	//
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// namespace corresponds to the JSON schema field "namespace".
	namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty" mapstructure:"namespace,omitempty"`

	// Parallel corresponds to the JSON schema field "parallel".
	Parallel *ParallelExecutableType `json:"parallel,omitempty" yaml:"parallel,omitempty" mapstructure:"parallel,omitempty"`

	// Render corresponds to the JSON schema field "render".
	Render *RenderExecutableType `json:"render,omitempty" yaml:"render,omitempty" mapstructure:"render,omitempty"`

	// Request corresponds to the JSON schema field "request".
	Request *RequestExecutableType `json:"request,omitempty" yaml:"request,omitempty" mapstructure:"request,omitempty"`

	// Serial corresponds to the JSON schema field "serial".
	Serial *SerialExecutableType `json:"serial,omitempty" yaml:"serial,omitempty" mapstructure:"serial,omitempty"`

	// Tags corresponds to the JSON schema field "tags".
	Tags ExecutableTags `json:"tags,omitempty" yaml:"tags,omitempty" mapstructure:"tags,omitempty"`

	// The maximum amount of time the executable is allowed to run before being
	// terminated.
	// The timeout is specified in Go duration format (e.g. 30s, 5m, 1h).
	//
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty" mapstructure:"timeout,omitempty"`

	// Verb corresponds to the JSON schema field "verb".
	Verb Verb `json:"verb" yaml:"verb" mapstructure:"verb"`

	// Visibility corresponds to the JSON schema field "visibility".
	Visibility *ExecutableVisibility `json:"visibility,omitempty" yaml:"visibility,omitempty" mapstructure:"visibility,omitempty"`

	// workspace corresponds to the JSON schema field "workspace".
	workspace string `json:"workspace,omitempty" yaml:"workspace,omitempty" mapstructure:"workspace,omitempty"`

	// workspacePath corresponds to the JSON schema field "workspacePath".
	workspacePath string `json:"workspacePath,omitempty" yaml:"workspacePath,omitempty" mapstructure:"workspacePath,omitempty"`
}

type ExecutableAliases common.Aliases

type ExecutableTags common.Tags

type ExecutableVisibility common.Visibility

// Launches an application or opens a URI.
type LaunchExecutableType struct {
	// The application to launch the URI with.
	App string `json:"app,omitempty" yaml:"app,omitempty" mapstructure:"app,omitempty"`

	// Args corresponds to the JSON schema field "args".
	Args ArgumentList `json:"args,omitempty" yaml:"args,omitempty" mapstructure:"args,omitempty"`

	// Params corresponds to the JSON schema field "params".
	Params ParameterList `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`

	// The URI to launch. This can be a file path or a web URL.
	URI string `json:"uri" yaml:"uri" mapstructure:"uri"`

	// If set to true, the executable will wait for the launched application to exit
	// before continuing.
	Wait bool `json:"wait,omitempty" yaml:"wait,omitempty" mapstructure:"wait,omitempty"`
}

type ParallelExecutableType struct {
	// Args corresponds to the JSON schema field "args".
	Args ArgumentList `json:"args,omitempty" yaml:"args,omitempty" mapstructure:"args,omitempty"`

	// A list of executables to run in parallel.
	// Each executable can be a command or a reference to another executable.
	//
	Execs ParallelRefConfigList `json:"execs" yaml:"execs" mapstructure:"execs"`

	// End the parallel execution as soon as an exec exits with a non-zero status.
	// This is the default behavior.
	// When set to false, all execs will be run regardless of the exit status of
	// parallel execs.
	//
	//
	FailFast *bool `json:"failFast,omitempty" yaml:"failFast,omitempty" mapstructure:"failFast,omitempty"`

	// The maximum number of threads to use when executing the parallel executables.
	MaxThreads int `json:"maxThreads,omitempty" yaml:"maxThreads,omitempty" mapstructure:"maxThreads,omitempty"`

	// Params corresponds to the JSON schema field "params".
	Params ParameterList `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

// Configuration for a parallel executable.
type ParallelRefConfig struct {
	// Arguments to pass to the executable.
	Args []string `json:"args,omitempty" yaml:"args,omitempty" mapstructure:"args,omitempty"`

	// The command to execute.
	// One of `cmd` or `ref` must be set.
	//
	Cmd string `json:"cmd,omitempty" yaml:"cmd,omitempty" mapstructure:"cmd,omitempty"`

	// A condition to determine if the executable should be run.
	If string `json:"if,omitempty" yaml:"if,omitempty" mapstructure:"if,omitempty"`

	// A reference to another executable to run in serial.
	// One of `cmd` or `ref` must be set.
	//
	Ref Ref `json:"ref,omitempty" yaml:"ref,omitempty" mapstructure:"ref,omitempty"`

	// The number of times to retry the executable if it fails.
	Retries int `json:"retries,omitempty" yaml:"retries,omitempty" mapstructure:"retries,omitempty"`
}

// A list of executables to run in parallel. The executables can be defined by it's
// exec `cmd` or `ref`.
type ParallelRefConfigList []ParallelRefConfig

// A parameter is a value that can be passed to an executable and all of its
// sub-executables.
// Only one of `text`, `secretRef`, or `prompt` must be set. Specifying more than
// one will result in an error.
type Parameter struct {
	// The name of the environment variable that will be assigned the value.
	EnvKey string `json:"envKey" yaml:"envKey" mapstructure:"envKey"`

	// A prompt to be displayed to the user when collecting an input value.
	Prompt string `json:"prompt,omitempty" yaml:"prompt,omitempty" mapstructure:"prompt,omitempty"`

	// A reference to a secret to be passed to the executable.
	SecretRef string `json:"secretRef,omitempty" yaml:"secretRef,omitempty" mapstructure:"secretRef,omitempty"`

	// A static value to be passed to the executable.
	Text string `json:"text,omitempty" yaml:"text,omitempty" mapstructure:"text,omitempty"`
}

type ParameterList []Parameter

// A reference to an executable.
// The format is `<verb> <workspace>/<namespace>:<executable name>`.
// For example, `exec ws/ns:my-workflow`.
//
// The workspace and namespace are optional.
// If the workspace is not specified, the current workspace will be used.
// If the namespace is not specified, the current namespace will be used.
type Ref string

type RefList []Ref

// Renders a markdown template file with data.
type RenderExecutableType struct {
	// Args corresponds to the JSON schema field "args".
	Args ArgumentList `json:"args,omitempty" yaml:"args,omitempty" mapstructure:"args,omitempty"`

	// Dir corresponds to the JSON schema field "dir".
	Dir Directory `json:"dir,omitempty" yaml:"dir,omitempty" mapstructure:"dir,omitempty"`

	// Params corresponds to the JSON schema field "params".
	Params ParameterList `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`

	// The path to the JSON or YAML file containing the template data.
	TemplateDataFile string `json:"templateDataFile,omitempty" yaml:"templateDataFile,omitempty" mapstructure:"templateDataFile,omitempty"`

	// The path to the markdown template file to render.
	TemplateFile string `json:"templateFile" yaml:"templateFile" mapstructure:"templateFile"`
}

// Makes an HTTP request.
type RequestExecutableType struct {
	// Args corresponds to the JSON schema field "args".
	Args ArgumentList `json:"args,omitempty" yaml:"args,omitempty" mapstructure:"args,omitempty"`

	// The body of the request.
	Body string `json:"body,omitempty" yaml:"body,omitempty" mapstructure:"body,omitempty"`

	// A map of headers to include in the request.
	Headers RequestExecutableTypeHeaders `json:"headers,omitempty" yaml:"headers,omitempty" mapstructure:"headers,omitempty"`

	// If set to true, the response will be logged as program output.
	LogResponse bool `json:"logResponse,omitempty" yaml:"logResponse,omitempty" mapstructure:"logResponse,omitempty"`

	// The HTTP method to use when making the request.
	Method RequestExecutableTypeMethod `json:"method,omitempty" yaml:"method,omitempty" mapstructure:"method,omitempty"`

	// Params corresponds to the JSON schema field "params".
	Params ParameterList `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`

	// ResponseFile corresponds to the JSON schema field "responseFile".
	ResponseFile *RequestResponseFile `json:"responseFile,omitempty" yaml:"responseFile,omitempty" mapstructure:"responseFile,omitempty"`

	// The timeout for the request in Go duration format (e.g. 30s, 5m, 1h).
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty" mapstructure:"timeout,omitempty"`

	// JQ query to transform the response before saving it to a file or outputting it.
	TransformResponse string `json:"transformResponse,omitempty" yaml:"transformResponse,omitempty" mapstructure:"transformResponse,omitempty"`

	// The URL to make the request to.
	URL string `json:"url" yaml:"url" mapstructure:"url"`

	// A list of valid status codes. If the response status code is not in this list,
	// the executable will fail.
	// If not set, the response status code will not be checked.
	//
	ValidStatusCodes []int `json:"validStatusCodes,omitempty" yaml:"validStatusCodes,omitempty" mapstructure:"validStatusCodes,omitempty"`
}

// A map of headers to include in the request.
type RequestExecutableTypeHeaders map[string]string

type RequestExecutableTypeMethod string

const RequestExecutableTypeMethodDELETE RequestExecutableTypeMethod = "DELETE"
const RequestExecutableTypeMethodGET RequestExecutableTypeMethod = "GET"
const RequestExecutableTypeMethodPATCH RequestExecutableTypeMethod = "PATCH"
const RequestExecutableTypeMethodPOST RequestExecutableTypeMethod = "POST"
const RequestExecutableTypeMethodPUT RequestExecutableTypeMethod = "PUT"

// Configuration for saving the response of a request to a file.
type RequestResponseFile struct {
	// Dir corresponds to the JSON schema field "dir".
	Dir Directory `json:"dir,omitempty" yaml:"dir,omitempty" mapstructure:"dir,omitempty"`

	// The name of the file to save the response to.
	Filename string `json:"filename" yaml:"filename" mapstructure:"filename"`

	// The format to save the response as.
	SaveAs RequestResponseFileSaveAs `json:"saveAs,omitempty" yaml:"saveAs,omitempty" mapstructure:"saveAs,omitempty"`
}

type RequestResponseFileSaveAs string

const RequestResponseFileSaveAsIndentedJson RequestResponseFileSaveAs = "indented-json"
const RequestResponseFileSaveAsJson RequestResponseFileSaveAs = "json"
const RequestResponseFileSaveAsRaw RequestResponseFileSaveAs = "raw"
const RequestResponseFileSaveAsYaml RequestResponseFileSaveAs = "yaml"
const RequestResponseFileSaveAsYml RequestResponseFileSaveAs = "yml"

// Executes a list of executables in serial.
type SerialExecutableType struct {
	// Args corresponds to the JSON schema field "args".
	Args ArgumentList `json:"args,omitempty" yaml:"args,omitempty" mapstructure:"args,omitempty"`

	// A list of executables to run in serial.
	// Each executable can be a command or a reference to another executable.
	//
	Execs SerialRefConfigList `json:"execs" yaml:"execs" mapstructure:"execs"`

	// End the serial execution as soon as an exec exits with a non-zero status. This
	// is the default behavior.
	// When set to false, all execs will be run regardless of the exit status of the
	// previous exec.
	//
	//
	FailFast *bool `json:"failFast,omitempty" yaml:"failFast,omitempty" mapstructure:"failFast,omitempty"`

	// Params corresponds to the JSON schema field "params".
	Params ParameterList `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

// Configuration for a serial executable.
type SerialRefConfig struct {
	// Arguments to pass to the executable.
	Args []string `json:"args,omitempty" yaml:"args,omitempty" mapstructure:"args,omitempty"`

	// The command to execute.
	// One of `cmd` or `ref` must be set.
	//
	Cmd string `json:"cmd,omitempty" yaml:"cmd,omitempty" mapstructure:"cmd,omitempty"`

	// A condition to determine if the executable should be run.
	If string `json:"if,omitempty" yaml:"if,omitempty" mapstructure:"if,omitempty"`

	// A reference to another executable to run in serial.
	// One of `cmd` or `ref` must be set.
	//
	Ref Ref `json:"ref,omitempty" yaml:"ref,omitempty" mapstructure:"ref,omitempty"`

	// The number of times to retry the executable if it fails.
	Retries int `json:"retries,omitempty" yaml:"retries,omitempty" mapstructure:"retries,omitempty"`

	// If set to true, the user will be prompted to review the output of the
	// executable before continuing.
	ReviewRequired bool `json:"reviewRequired,omitempty" yaml:"reviewRequired,omitempty" mapstructure:"reviewRequired,omitempty"`
}

// A list of executables to run in serial. The executables can be defined by it's
// exec `cmd` or `ref`.
type SerialRefConfigList []SerialRefConfig

type Verb string

const VerbAbort Verb = "abort"
const VerbActivate Verb = "activate"
const VerbAdd Verb = "add"
const VerbAnalyze Verb = "analyze"
const VerbApply Verb = "apply"
const VerbBuild Verb = "build"
const VerbBundle Verb = "bundle"
const VerbCheck Verb = "check"
const VerbClean Verb = "clean"
const VerbClear Verb = "clear"
const VerbCompile Verb = "compile"
const VerbConfigure Verb = "configure"
const VerbCreate Verb = "create"
const VerbDeactivate Verb = "deactivate"
const VerbDelete Verb = "delete"
const VerbDeploy Verb = "deploy"
const VerbDestroy Verb = "destroy"
const VerbDisable Verb = "disable"
const VerbEdit Verb = "edit"
const VerbEnable Verb = "enable"
const VerbErase Verb = "erase"
const VerbExec Verb = "exec"
const VerbExecute Verb = "execute"
const VerbFetch Verb = "fetch"
const VerbGenerate Verb = "generate"
const VerbGet Verb = "get"
const VerbInit Verb = "init"
const VerbInspect Verb = "inspect"
const VerbInstall Verb = "install"
const VerbKill Verb = "kill"
const VerbLaunch Verb = "launch"
const VerbLint Verb = "lint"
const VerbManage Verb = "manage"
const VerbModify Verb = "modify"
const VerbMonitor Verb = "monitor"
const VerbNew Verb = "new"
const VerbOpen Verb = "open"
const VerbPackage Verb = "package"
const VerbPatch Verb = "patch"
const VerbPause Verb = "pause"
const VerbPublish Verb = "publish"
const VerbPurge Verb = "purge"
const VerbPush Verb = "push"
const VerbReboot Verb = "reboot"
const VerbRefresh Verb = "refresh"
const VerbRelease Verb = "release"
const VerbReload Verb = "reload"
const VerbRemove Verb = "remove"
const VerbRequest Verb = "request"
const VerbReset Verb = "reset"
const VerbRestart Verb = "restart"
const VerbRetrieve Verb = "retrieve"
const VerbRun Verb = "run"
const VerbScan Verb = "scan"
const VerbSend Verb = "send"
const VerbSet Verb = "set"
const VerbSetup Verb = "setup"
const VerbShow Verb = "show"
const VerbStart Verb = "start"
const VerbStop Verb = "stop"
const VerbTeardown Verb = "teardown"
const VerbTerminate Verb = "terminate"
const VerbTest Verb = "test"
const VerbTidy Verb = "tidy"
const VerbTrack Verb = "track"
const VerbTransform Verb = "transform"
const VerbTrigger Verb = "trigger"
const VerbUndeploy Verb = "undeploy"
const VerbUninstall Verb = "uninstall"
const VerbUnset Verb = "unset"
const VerbUpdate Verb = "update"
const VerbUpgrade Verb = "upgrade"
const VerbValidate Verb = "validate"
const VerbVerify Verb = "verify"
const VerbView Verb = "view"
const VerbWatch Verb = "watch"

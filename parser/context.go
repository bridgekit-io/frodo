package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/bridgekit-io/frodo/internal/naming"
	"github.com/bridgekit-io/frodo/internal/slices"
	"golang.org/x/tools/go/packages"
)

// Context wrangles all of the captured data about your input service declaration file. It tracks
// the module/package information, the service(s) that were defined, the request/response structs
// that were defined in the file, etc. It's the output of Parse() and is the input value when we
// evaluate Go templates to generate other source files based on this service definition info.
type Context struct {
	// --- Fields tracking the input file that we parsed.

	// FileSet is the collection of related files we're going to give to the Go AST parser.
	FileSet *token.FileSet
	// File is the entire syntax tree from when we parsed your input file.
	File *ast.File
	// Path is the relative path to the service definition file we're parsing.
	Path string
	// AbsolutePath is the absolute path to the service definition file we're parsing.
	AbsolutePath string
	// Timestamp is when we performed the parsing.
	Timestamp time.Time

	// --- Fields related to orienting ourselves to the user's module/package.

	// Module contains info from "go.mod" about the entire module where the service/package is defined.
	Module *ModuleDeclaration
	// InputPackage contains information about the package where the service definition resides.
	InputPackage *PackageDeclaration
	// OutputPackage contains information about the package where the generated code will go.
	OutputPackage *PackageDeclaration

	// --- Fields related to the actual service info that we're parsing.

	// Service contains the parsed info about the service declared in our input file.
	Service *ServiceDeclaration
	// Types captures snapshots of every model type, primitive type, and recursive field type for
	// any field referenced by any of your request/response models. It contains all of the referenced types
	// for every single field and their fields and their fields, and so on.
	Types TypeRegistry

	// --- Fields that cache various types of parsed syntax/type info so the parser can look it up easily.

	// Documentation stores the GoDoc comments for the services/functions/models/fields in the parsed code.
	Documentation Documentation
	// Tags stores the struct field tags annotated on fields in your input file.
	Tags Tags
	// RawTypes contains the tree of all raw, parsed type information as we get it from the AST.
	RawTypes *packages.Package
}

// Scope returns the root of the parsed type tree for the source file we parsed.
func (ctx Context) Scope() *types.Scope {
	return ctx.RawTypes.Types.Scope()
}

// TimestampString returns the timestamp formatted in a standard fashion that we include in artifact headers.
func (ctx Context) TimestampString() string {
	return ctx.Timestamp.Format(time.RFC1123)
}

// ServiceDeclaration wrangles all of the information we could grab about the service from the
// interface that defined it.
type ServiceDeclaration struct {
	// Name is the name of the service/interface.
	Name string
	// Version is the (hopefully) semantic version of your API (e.g. 1.2.0). This is NOT the prefix
	// to all routes in the API for the service. It's just an identifier available to code gen tools.
	Version string
	// Gateway contains the configuration HTTP-related options for this service.
	Gateway *GatewayServiceOptions
	// Functions are all of the functions explicitly defined on this service.
	Functions ServiceFunctionDeclarations
	// Documentation are all of the comments documenting this service.
	Documentation DocumentationLines
}

// ShortName returns the name of the service w/o the "Service"
// suffix. (e.g. "MathService" -> "Math")
func (service ServiceDeclaration) ShortName() string {
	return strings.TrimSuffix(service.Name, "Service")
}

// UnexportedName returns the name of the service with the first letter
// lower case. (e.g. "MathService" -> "mathService")
func (service ServiceDeclaration) UnexportedName() string {
	if service.Name == "" {
		return ""
	}
	firstRune, _ := utf8.DecodeRuneInString(service.Name)
	return fmt.Sprintf("%c%s", unicode.ToLower(firstRune), service.Name[1:])
}

// FunctionByName fetches the service operation with the given function name. This returns nil when there
// are no functions in this interface/service by that name.
func (service ServiceDeclaration) FunctionByName(name string) *ServiceFunctionDeclaration {
	for _, m := range service.Functions {
		if strings.EqualFold(m.Name, name) {
			return m
		}
	}
	return nil
}

// ServiceFunctionDeclarations defines a collection of related service functions/operations.
type ServiceFunctionDeclarations []*ServiceFunctionDeclaration

// ServiceFunctionDeclaration defines a single operation/function within a service (one of the interface functions).
type ServiceFunctionDeclaration struct {
	// Name is the name of the function defined in the service interface (the function name to call this operation).
	Name string
	// Request contains the details about the model/type/struct for this operation's input/request value.
	Request *TypeDeclaration
	// Response contains the details about the model/type/struct for this operation's output/response value.
	Response *TypeDeclaration
	// Routes contains registration information so this function can be invoked by one or more gateways. It
	// contains things like the HTTP method/path to register w/ the API gateway or the name of the topic/key
	// to subscribe to in an event gateway.
	Routes GatewayRoutes
	// Roles defines the role-based security identifiers that a user/principal must have in order to access
	// this endpoint. These can be exact values like "admin.write" or parameterized like "group.{Group.ID}.write".
	Roles []string
	// Documentation are all of the comments documenting this operation.
	Documentation DocumentationLines
	// Service represents the interface/service that this function belongs to.
	Service *ServiceDeclaration
}

// String returns the function signature for this operation for debugging purposes.
func (f ServiceFunctionDeclaration) String() string {
	return fmt.Sprintf("%s(context.Context, *%v) (*%v, error)",
		f.Name,
		f.Request,
		f.Response,
	)
}

// FieldDeclarations collects the fields/attributes on a service model.
type FieldDeclarations []*FieldDeclaration

// Empty returns true if there are zero fields in this set.
func (fields FieldDeclarations) Empty() bool {
	return len(fields) == 0
}

// NotEmpty returns true if there is at least one field in this set
func (fields FieldDeclarations) NotEmpty() bool {
	return len(fields) > 0
}

// ByName looks up the declaration for the field that matches the given name. This name
// comparison is CASE INSENSITIVE, so "id" will find the field "ID".
func (fields FieldDeclarations) ByName(name string) *FieldDeclaration {
	for _, field := range fields {
		if strings.EqualFold(field.Name, name) {
			return field
		}
	}
	return nil
}

// TransportFields returns the subset of child fields for this field that
// are NOT omitted (i.e. should be included in JSON/transport).
func (fields FieldDeclarations) TransportFields() FieldDeclarations {
	var result FieldDeclarations
	for _, f := range fields {
		if f.Binding.NotOmit() {
			result = append(result, f)
		}
	}
	return result
}

// ByBindingName looks for a field whose (possibly) re-mapped name matches the given value. This
// comparison is CASE INSENSITIVE, so "id" will find the field "ID".
func (fields FieldDeclarations) ByBindingName(name string) *FieldDeclaration {
	nextFieldName, recursiveFields, _ := strings.Cut(name, ".")
	for _, field := range fields {
		if !strings.EqualFold(field.Binding.Name, nextFieldName) {
			continue
		}

		if recursiveFields == "" {
			return field
		}
		return field.Type.Fields.ByBindingName(recursiveFields)
	}
	return nil
}

// FieldDeclaration describes a single field in a request/response model.
type FieldDeclaration struct {
	// Name the name of the field/attribute.
	Name string
	// ParentType is a back-pointer to the type that this field is a member of.
	ParentType *TypeDeclaration
	// Type contains the data type information for this field.
	Type *TypeDeclaration
	// Pointer indicates if this field is a pointer type (true) or value type (false).
	Pointer bool
	// Documentation are all of the comments documenting this field.
	Documentation DocumentationLines
	// Binding describes the custom binding instructions used when unmarshalling request data onto this field.
	Binding *FieldBindingOptions
}

// FieldBindingOptions provides hints to the generation tools about how the runtime binder will
// map request parameters to an attribute of the request struct.
type FieldBindingOptions struct {
	// Omit will be true if you provided the `json:"-"` tag saying that this is not part of JSON marshaling.
	Omit bool
	// Name is the remapped JSON attribute for the associated field (e.g. `json:"user_id"` -> user_id).
	Name string
}

// NotOmit is a convenience for templates that returns true when we should expose this field to
// client and documentation templates/tooling.
func (opts FieldBindingOptions) NotOmit() bool {
	return !opts.Omit
}

// ModuleDeclaration contains information about the Go module that the service belongs
// to. This is information scraped from project's "go.mod" file.
type ModuleDeclaration struct {
	// Name is the fully qualified module name (e.g. "github.com/some_user/module_name")
	Name string
	// Directory is the absolute path to the root directory of the module (where go.mod resides).
	Directory string
}

// GoMod returns the absolute path to the "go.mod" file for this module on the system running frodo.
func (module ModuleDeclaration) GoMod() string {
	return filepath.Join(module.Directory, "go.mod")
}

// PackageDeclaration defines the subpackage that the service resides in.
type PackageDeclaration struct {
	// Name is just the raw package name (no path info)
	Name string
	// Import is the fully qualified package name (e.g. "github.com/some_user/module_name/foo/bar/baz")
	Import string
	// Directory is the absolute path to the package.
	Directory string
}

var noDocumentation = DocumentationLines{}

// Documentation is a lookup cache for all GoDoc comments on your services, functions, models, and fields.
type Documentation map[string]string

// Set adds an entry to the lookup. This is admittedly a bastardization of variadic functions - the last value you
// pass in is the GoDoc comments. The first to second-to-last values are segments in the lookup key. For instance
// if you are caching the comment for the function Bar on the FooService, you would
// call `Set("FooService", "Bar", "Bar does some baz magic and gives you back goo.")`. The resulting entry will
// look like "FooService.Bar"->"Bar does some...".
func (docs Documentation) Set(segmentsAndDoc ...string) {
	length := len(segmentsAndDoc)
	if length < 2 {
		return
	}

	key := strings.Join(segmentsAndDoc[:length-1], ".")
	docs[key] = strings.TrimSpace(segmentsAndDoc[length-1])
}

func (docs Documentation) lookup(segments ...string) DocumentationLines {
	if comments, ok := docs[strings.Join(segments, ".")]; ok {
		return strings.Split(comments, "\n")
	}
	return noDocumentation
}

// ForService finds the GoDoc comments for the given service interface.
func (docs Documentation) ForService(s *ServiceDeclaration) DocumentationLines {
	return docs.lookup(s.Name)
}

// ForFunction finds the GoDoc comments for the given service function.
func (docs Documentation) ForFunction(f *ServiceFunctionDeclaration) DocumentationLines {
	return docs.lookup(f.Service.Name, f.Name)
}

// ForType finds the GoDoc comments for the request/response struct.
func (docs Documentation) ForType(t *TypeDeclaration) DocumentationLines {
	return docs.lookup(t.Name)
}

// ForField find the GoDoc comments for the attribute of a request/response struct
func (docs Documentation) ForField(f *FieldDeclaration) DocumentationLines {
	return docs.lookup(f.ParentType.Name, f.Name)
}

var noTag = reflect.StructTag("")

// Tags is a lookup for finding `json:"xxx"` tags defined on your request/response structs.
type Tags map[string]string

// Set captures the tag information for the given model attribute.
func (tags Tags) Set(model string, field string, tag *ast.BasicLit) {
	if tag == nil {
		return
	}
	if model == "" || field == "" {
		return
	}

	// When we pull tags off of the AST, they're still wrapped in the `xxx` back ticks, so
	// pull those off before giving them back to the caller.
	tags[model+"."+field] = strings.Trim(strings.TrimSpace(tag.Value), "`")
}

// ForField looks up the tag annotations for the given model field. If the field does not have any tags
// you'll get back the zero-value StructTag that always gives you empty for any value lookups.
func (tags Tags) ForField(f *FieldDeclaration) reflect.StructTag {
	if tag, ok := tags[f.ParentType.Name+"."+f.Name]; ok {
		return reflect.StructTag(tag)
	}
	return noTag
}

// DocumentationLines represents all of the 'go doc' lines above a type/function/field with all
// of the leading slashes removed.
type DocumentationLines []string

// Trim removes blank doc lines from the front/back of your list of comments.
func (docs DocumentationLines) Trim() DocumentationLines {
	if docs.Empty() {
		return docs
	}
	// We want to be able to trim leading and trailing blank lines in a single pass over the
	// slice, so the first time we encounter a non-empty line that's the first index to keep.
	// The last index will continuously be updated as we go further and find other non-empty
	// lines, so by the time we finish the loop we should have both ends of the valid range.
	first := -1
	last := -1
	for i, line := range docs {
		switch {
		case strings.TrimSpace(line) == "":
			// don't update anything
		case first < 0:
			// we just found the first non-blank comment
			first = i
			last = i
		default:
			// there's another non-blank comment line further after the first one
			last = i
		}
	}
	if first == -1 {
		return DocumentationLines{}
	}
	return docs[first : last+1]
}

// NotEmpty returns true when there is at least 1 line of documentation/comments.
func (docs DocumentationLines) NotEmpty() bool {
	return !docs.Empty()
}

// Empty returns true when there are no lines of documentation/comments for the service/field/function/etc.
func (docs DocumentationLines) Empty() bool {
	for _, doc := range docs {
		if doc != "" {
			return false
		}
	}
	return true
}

// String returns all of the documentation lines as one string, separated by newlines.
func (docs DocumentationLines) String() string {
	if docs.Empty() {
		return ""
	}
	return strings.Join(docs, "\n")
}

// GatewayServiceOptions contains all of the configurable HTTP-related options for a top-level service.
type GatewayServiceOptions struct {
	// Service is a back-pointer to the service these options correspond to.
	Service *ServiceDeclaration
	// PathPrefix is the optional version/domain prefix for all endpoints in the API (e.g. "v2/").
	PathPrefix string
}

// GatewayRoutes manages the collection of all possible routes that you can register to
// the runtime gateways serving up requests for this service.
type GatewayRoutes []*GatewayRoute

// Get returns the route for the given type of gateway (e.g. "API" or "EVENTS"). This
// will return nil if there is no route for the given type of gateway.
func (routes GatewayRoutes) Get(gatewayType string) *GatewayRoute {
	for _, route := range routes {
		if route.GatewayType == gatewayType {
			return route
		}
	}
	return nil
}

// API returns the route that we'll register w/ the runtime HTTP/API gateway. This
// will be nil if you provided options to omit the endpoint from the public API.
func (routes GatewayRoutes) API() *GatewayRoute {
	return routes.Get("API")
}

// Events returns the route that we'll register w/ the runtime event source gateway. This
// will be nil if the user did not provide any "ON" doc options. Since you can register
// more than one of these per endpoint, you have to specify which one you want.
func (routes GatewayRoutes) Events() GatewayRoutes {
	var results GatewayRoutes
	for _, route := range routes {
		if route.GatewayType == "EVENTS" {
			results = append(results, route)
		}
	}
	return results
}

// RouteType describes how the gateway or client should handle implementation of this endpoint (e.g. REST request vs websocket).
type RouteType string

const (
	// RouteTypeDefault indicates that servers/clients should build handlers for this route using the standard implementation for
	// the current gateway. For the API gateway, this means to build it just like a normal REST endpoint. For the Events gateway, this
	// means ingesting events from the broker.
	RouteTypeDefault = RouteType("")
	// RouteTypeWebsocket is only used by the API gateway when generating clients to indicate that the implementation should create a
	// websocket to this endpoint rather than dispatching a standard HTTP request.
	RouteTypeWebsocket = RouteType("Websocket")
)

// Default returns true if this type is not one of our specialized types like Websocket.
func (t RouteType) Default() bool {
	return !t.Websocket()
}

// Websocket returns true if this type indicates that this route should be implemented using a websocket.
func (t RouteType) Websocket() bool {
	return t == RouteTypeWebsocket
}

// GatewayRoute contains the information required to register a method/endpoint with a specific
// gateway/listener. For instance, an operation might be triggered via some HTTP endpoint in
// your API, or it might be asynchronously triggered by listening for a certain event. In that
// case, you'd have 2 routes; one with the details for the API gateway and another for the event gateway.
type GatewayRoute struct {
	// Function is a back-pointer to the service function these options correspond to.
	Function *ServiceFunctionDeclaration
	// GatewayType is a descriptor for the type of gateway this route should register with (e.g. "API" or "EVENT").
	GatewayType string
	// Method indicates if the RPC gateway should use a GET, POST, etc. when exposing this operation via HTTP.
	Method string
	// Path defines the URL pattern to provide to the gateway's router/mux to access this operation.
	Path string
	// Status indicates what success status code the gateway should use when responding via HTTP (e.g. 200, 202, etc.)
	Status int
	// Group provides additional routing/grouping info that means different things to different gateways.
	Group string
	// RouteType describes how the gateway or client should handle implementation of this endpoint (e.g. REST request vs websocket).
	RouteType RouteType
}

// QualifiedPath returns the route's path with the service's PathPrefix prepended to it. This includes a leading "/"
// so that all clients/servers can expect that to be in place (so we don't need to re-normalize this over and over).
func (route *GatewayRoute) QualifiedPath() string {
	prefix := strings.Trim(route.Function.Service.Gateway.PathPrefix, "/")
	path := strings.Trim(route.Path, "/")

	// The pub/sub gateway doesn't care about your HTTP remapping of paths. The event keys are locked to the name
	// service, dot, name of the function.
	if route.GatewayType == "EVENTS" {
		return path
	}

	// Only prepend the HTTP path prefix (e.g. "v2") if there's one defined on the service in the first place.
	if prefix != "" {
		return "/" + prefix + "/" + path
	}
	return "/" + path
}

// ParsePathParams extracts a slice of the path parameter names from the Path of this route. For instance the
// path "/user/{UserID}/transaction/{TransactionID}" will parse to []string{"UserID", "TransactionID"}.
func (route *GatewayRoute) ParsePathParams() []string {
	var params []string
	segments := strings.Split(strings.TrimSpace(route.Path), "/")
	for _, segment := range segments {
		if !strings.HasPrefix(segment, "{") && !strings.HasSuffix(segment, "}") {
			continue
		}

		param := strings.TrimPrefix(strings.TrimSuffix(segment, "}"), "{")
		params = append(params, param)
	}
	return params
}

// MethodMatches returns true when the route's method matches at least one of 'these' options.
func (route *GatewayRoute) MethodMatches(these ...string) bool {
	for _, method := range these {
		if strings.EqualFold(route.Method, method) {
			return true
		}
	}
	return false
}

// SupportsRequestBody returns true when the method is either POST, PUT, or PATCH; the HTTP methods
// where we expect you to feed request data via the request body rather than query string.
func (route *GatewayRoute) SupportsRequestBody() bool {
	method := strings.ToUpper(route.Method)
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch
}

// PathParameters looks at all of the ":xxx" path parameters in HTTPPath and returns the fields on
// the request struct that will be bound by those values at runtime. For instance, if the path
// was "/user/{userID}/address/{addressID}", this will return a 2-element slice containing the request's
// UserID and AddressID fields.
func (route *GatewayRoute) PathParameters() GatewayParameters {
	var results GatewayParameters
	for _, segment := range strings.Split(route.Path, "/") {
		paramName := naming.PathVariableName(segment)
		if paramName == "" {
			continue
		}

		field := route.Function.Request.Fields.ByBindingName(paramName)
		if field == nil {
			continue
		}

		results = append(results, &GatewayParameter{
			Name:  paramName,
			Field: field,
		})
	}
	return results
}

// QueryParameters describes all of the request struct attributes that can be bound by specifying
// them in the query string of the URL when making a request. For instance, if your request struct
// had an attribute "Limit uint64", then this includes a GatewayParameter that describes the
// caller's ability to include "&Limit=123" in the query string.
func (route *GatewayRoute) QueryParameters() GatewayParameters {
	var results GatewayParameters

	// If you're doing a POST/PUT/PATCH, we expect every value to come from either
	// the body or the path, not the query string.
	if route.SupportsRequestBody() {
		return results
	}

	pathParams := route.PathParameters()

	for _, field := range route.Function.Request.Fields {
		// Exclude any fields that will be bound using path parameters.
		if pathParams.ByName(field.Binding.Name) != nil {
			continue
		}

		results = append(results, &GatewayParameter{
			Name:  field.Binding.Name,
			Field: field,
		})
	}
	return results
}

// GatewayParameters is an overlay of a service function's path and request type/field info. It helps you
// indicate how a given field will be bound when handling incoming requests (e.g. path params vs query params).
type GatewayParameters []*GatewayParameter

// ByName locates the parameter with the given name. This is a case-insensitive search.
func (params GatewayParameters) ByName(name string) *GatewayParameter {
	for _, param := range params {
		if strings.EqualFold(param.Name, name) {
			return param
		}
	}
	return nil
}

// Empty returns true when there are zero parameters defined in this set.
func (params GatewayParameters) Empty() bool {
	return len(params) == 0
}

// NotEmpty returns true when there is at least one parameter mapping defined.
func (params GatewayParameters) NotEmpty() bool {
	return !params.Empty()
}

// GatewayParameter defines how a path/query parameter will be bound to a field in your request struct.
type GatewayParameter struct {
	// Name is the identifier of the path param (e.g "id" in "/user/{id}") or query string value that
	// will be bound to the Field.
	Name string
	// Field indicates which model attribute will be populated when this parameter goes
	// through the request binder.
	Field *FieldDeclaration
}

// TypeDeclaration is a snapshot of the type information for any type referenced, directly or indirectly, in
// your service declaration file.
type TypeDeclaration struct {
	// Name is the "package.Type" formatted name of the type as it is used in our source file. For types that
	// are defined in the declaration file, the name might only be "FooRequest" (no package).
	Name string
	// Type is the raw Go type information that we used to build this snapshot.
	Type types.Type
	// Kind classifies our types. Our definition aligns fairly closely with the "reflect" package's definition, so
	// we use it to describe whether a type, at its lowest level, describes a struct, int, slice, etc.
	Kind reflect.Kind
	// Elem is used by slice and map-like declarations to describe the type info for the underlying value type.
	Elem *TypeDeclaration
	// Key is used by map-like declarations to describe the type info for the lookup key.
	Key *TypeDeclaration
	// Basic indicates if this type is one of our base types (primitives, etc).
	Basic bool
	// Fields (for struct types) contains all of the attribute/type info for all members of the struct.
	Fields FieldDeclarations
	// Documentation are all of the comments documenting this operation.
	Documentation DocumentationLines
	// Implements contains some quick checks for whether this type implements the various
	// single function interfaces used to handle raw data responses.
	Implements struct {
		// MarshalJSON is true when the type implements custom JSON marshal logic.
		MarshalJSON bool
		// ContentGetter is true when it implements that interface.
		ContentGetter bool
		// ContentSetter is true when it implements that interface.
		ContentSetter bool
		// ContentTypeGetter is true when it implements that interface.
		ContentTypeGetter bool
		// ContentTypeSetter is true when it implements that interface.
		ContentTypeSetter bool
		// ContentLengthGetter is true when it implements that interface.
		ContentLengthGetter bool
		// ContentLengthSetter is true when it implements that interface.
		ContentLengthSetter bool
		// ContentRangeGetter is true when it implements that interface.
		ContentRangeGetter bool
		// ContentRangeSetter is true when it implements that interface.
		ContentRangeSetter bool
		// ContentFileNameGetter is true when it implements that interface.
		ContentFileNameGetter bool
		// ContentFileNameSetter is true when it implements that interface.
		ContentFileNameSetter bool
	}
}

// String returns the name of the type. That is all.
func (t TypeDeclaration) String() string {
	return t.Name
}

// SliceLike returns true for array or slice types. This will also be true for any alias to an array/slice type.
func (t TypeDeclaration) SliceLike() bool {
	return t.Kind == reflect.Slice || t.Kind == reflect.Array
}

// MapLike returns true for map types. This will also be true for any alias to a map type.
func (t TypeDeclaration) MapLike() bool {
	return t.Kind == reflect.Map
}

// PrimitiveLike returns true for types that represent some sort of basic/primitive type like numbers/booleans/strings.
func (t TypeDeclaration) PrimitiveLike() bool {
	return t.Kind == reflect.String ||
		t.Kind == reflect.Bool ||
		t.Kind == reflect.Int ||
		t.Kind == reflect.Int8 ||
		t.Kind == reflect.Int16 ||
		t.Kind == reflect.Int32 ||
		t.Kind == reflect.Int64 ||
		t.Kind == reflect.Uint ||
		t.Kind == reflect.Uint8 ||
		t.Kind == reflect.Uint16 ||
		t.Kind == reflect.Uint32 ||
		t.Kind == reflect.Uint64 ||
		t.Kind == reflect.Float32 ||
		t.Kind == reflect.Float64
}

// ObjectLike returns true for types that represent some sort complex object (i.e. struct/interface).
func (t TypeDeclaration) ObjectLike() bool {
	return t.Kind == reflect.Struct || t.Kind == reflect.Interface
}

// NonOmittedFields returns just the subset of fields that should be included in this model's transport/binding.
func (t TypeDeclaration) NonOmittedFields() FieldDeclarations {
	var results FieldDeclarations
	for _, f := range t.Fields {
		if !f.Binding.Omit {
			results = append(results, f)
		}
	}
	return results
}

// TypeRegistry is a quick lookup of all types we encountered when processing your declaration file.
type TypeRegistry map[string]*TypeDeclaration

// Lookup finds our parsed snapshot for the raw type. It returns the snapshot an "ok" boolean to indicate
// whether we found it or not.
func (reg TypeRegistry) Lookup(t types.Type) (*TypeDeclaration, bool) {
	key := reg.key(t, t.String())
	info, ok := reg[key]
	return info, ok
}

// LookupByName finds our parsed snapshot for the type name. It returns the snapshot an "ok" boolean to indicate
// whether we found it or not.
func (reg TypeRegistry) LookupByName(name string) (*TypeDeclaration, bool) {
	info, ok := reg[strings.ToLower(name)]
	return info, ok
}

// Register adds the given type declaration to the registry.
func (reg TypeRegistry) Register(entry *TypeDeclaration) *TypeDeclaration {
	key := reg.key(entry.Type, entry.Name)
	reg[key] = entry
	return entry
}

// WithoutInvalid removes all entries for types whose 'Kind' is 'Invalid'. This ensures that
// your context doesn't contain any entries for fields/types that we don't support.
func (reg TypeRegistry) WithoutInvalid() TypeRegistry {
	for key, typeDeclaration := range reg {
		if typeDeclaration.Kind == reflect.Invalid {
			delete(reg, key)
		}
	}
	return reg
}

// NonBasicTypes returns a slice containing only types not declared as "Basic". This way you only
// iterate complex types that you defined or imported.
func (reg TypeRegistry) NonBasicTypes() []*TypeDeclaration {
	var results []*TypeDeclaration
	for _, t := range reg {
		if t.Basic {
			continue
		}
		if t.Kind == reflect.Interface {
			continue
		}
		results = append(results, t)
	}
	return results
}

// ExternalPackageNames returns a set of any package names used by types in this service that are not defined in the service
// package. These could be standard library package (e.g. time.Time would include "time") or some other core package in your
// application (e.g. datastore.ID would include "datastore").
func (reg TypeRegistry) ExternalPackageNames() []string {
	packageNames := make([]string, 0, 8)
	for _, t := range reg {
		if packageName, _, ok := strings.Cut(t.Name, "."); ok {
			packageNames = slices.AppendUnique(packageNames, packageName)
		}
	}
	return packageNames
}

func (reg TypeRegistry) key(t types.Type, name string) string {
	if t != nil {
		name = t.String()
	}
	name = naming.NoPointer(name)
	name = naming.CleanPrefix(name)
	return strings.ToLower(name)
}

// NewTypeRegistry creates a new type registry that is pre-filled with all of the base types already registered.
func NewTypeRegistry() TypeRegistry {
	return TypeRegistry{
		"string":  &TypeDeclaration{Basic: true, Name: "string", Kind: reflect.String},
		"bool":    &TypeDeclaration{Basic: true, Name: "bool", Kind: reflect.Bool},
		"rune":    &TypeDeclaration{Basic: true, Name: "rune", Kind: reflect.Int32},
		"byte":    &TypeDeclaration{Basic: true, Name: "byte", Kind: reflect.Int8},
		"int":     &TypeDeclaration{Basic: true, Name: "int", Kind: reflect.Int},
		"int8":    &TypeDeclaration{Basic: true, Name: "int8", Kind: reflect.Int8},
		"int16":   &TypeDeclaration{Basic: true, Name: "int16", Kind: reflect.Int16},
		"int32":   &TypeDeclaration{Basic: true, Name: "int32", Kind: reflect.Int32},
		"int64":   &TypeDeclaration{Basic: true, Name: "int64", Kind: reflect.Int64},
		"uint":    &TypeDeclaration{Basic: true, Name: "uint", Kind: reflect.Uint},
		"uint8":   &TypeDeclaration{Basic: true, Name: "uint8", Kind: reflect.Uint8},
		"uint16":  &TypeDeclaration{Basic: true, Name: "uint16", Kind: reflect.Uint16},
		"uint32":  &TypeDeclaration{Basic: true, Name: "uint32", Kind: reflect.Uint32},
		"uint64":  &TypeDeclaration{Basic: true, Name: "uint64", Kind: reflect.Uint64},
		"float32": &TypeDeclaration{Basic: true, Name: "float32", Kind: reflect.Float32},
		"float64": &TypeDeclaration{Basic: true, Name: "float64", Kind: reflect.Float64},

		// Yes, time is technically a struct, but for the purposes of transport, we want to treat
		// time as an ISO string, so we're special casing this bad boy.
		"time.Time": &TypeDeclaration{Basic: true, Name: "time.Time", Kind: reflect.String},

		// Not supported in code generation, but there so we don't have nil pointers if you
		// are silly enough to use them as service function inputs/outputs.
		"uintptr":    &TypeDeclaration{Basic: true, Name: "uintptr", Kind: reflect.Uintptr},
		"complex64":  &TypeDeclaration{Basic: true, Name: "complex64", Kind: reflect.Complex64},
		"complex128": &TypeDeclaration{Basic: true, Name: "complex128", Kind: reflect.Complex128},
	}
}

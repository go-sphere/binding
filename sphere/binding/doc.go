// Package binding holds the Go bindings for sphere/binding/binding.proto, the
// protobuf options that say where each field of a request message is read
// from in an HTTP request: query string, URI path, JSON body, form, or header.
//
// The options are consumed at code-generation time. protoc-gen-sphere-binding
// turns them into Go struct tags such as uri:"id" and query:"page" on the
// protoc-gen-go output, and protoc-gen-sphere uses them when generating HTTP
// handlers. Nothing in this package runs at request time.
//
// # Declaring bindings in .proto
//
// Depend on buf.build/go-sphere/binding, then:
//
//	import "sphere/binding/binding.proto";
//
//	message ListUsersRequest {
//	  option (sphere.binding.default_location) = BINDING_LOCATION_QUERY;
//	  string org_id = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
//	  int32 page = 2;
//	  string token = 3 [
//	    (sphere.binding.location) = BINDING_LOCATION_HEADER,
//	    (sphere.binding.tags) = "validate:\"required\""
//	  ];
//	}
//
// A field's location comes from, in order of precedence, its own location
// option, its oneof's default_oneof_location, then its message's
// default_location. QUERY, URI, and HEADER apply only to scalar fields and
// well-known scalar wrappers such as google.protobuf.Timestamp.
//
// # Reading options in Go
//
// Generators and tools read the options from descriptors with the E_ extension
// variables and [BindingLocation]:
//
//	import (
//		"github.com/go-sphere/binding/sphere/binding"
//		"google.golang.org/protobuf/proto"
//	)
//
//	opts := fieldDescriptor.Options() // protoreflect.FieldDescriptor
//	if proto.HasExtension(opts, binding.E_Location) {
//		loc := proto.GetExtension(opts, binding.E_Location).(binding.BindingLocation)
//		_ = loc
//	}
//
// Slices returned by proto.GetExtension for the repeated options (tags and
// auto tags) belong to the descriptor and must not be modified.
package binding

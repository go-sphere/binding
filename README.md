# binding proto

This is a proto file that defines the Sphere binding framework for HTTP request binding annotations. It provides extensions for protobuf messages to support automatic generation of Go struct tags for HTTP request binding in Go applications.

## Features

- Automatic field binding to HTTP request components (query, URI, body)
- Custom struct tag generation for advanced use cases
- Default auto tags for consistent field mapping
- Message-level and field-level binding configuration
- Designed for use with [`protoc-gen-sphere-binding`](https://github.com/go-sphere/protoc-gen-sphere-binding) plugin

## Proto Definition

The core proto file defines:

```protobuf
syntax = "proto3";

package sphere.binding;

import "google/protobuf/descriptor.proto";

enum BindingLocation {
  BINDING_LOCATION_UNSPECIFIED = 0;
  BINDING_LOCATION_QUERY = 1;    // HTTP query parameters
  BINDING_LOCATION_URI = 2;      // URI path parameters
  BINDING_LOCATION_BODY = 3;     // HTTP request body
}

extend google.protobuf.FieldOptions {
  BindingLocation location = 18534200;  // Field binding location
  string tags = 18534210;               // Custom struct tags
  string auto_tags = 18534220;          // Automatic tags
}

extend google.protobuf.MessageOptions {
  BindingLocation default_location = 18534230;      // Default location for all fields
  string default_auto_tags = 18534240;              // Default auto tags for all fields
}

extend google.protobuf.OneofOptions {
  BindingLocation default_oneof_location = 18534250;    // Default location for oneof fields
  string default_oneof_auto_tags = 18534260;            // Default auto tags for oneof fields
}
```

### Basic Request Binding

```protobuf
syntax = "proto3";

package api.v1;

import "sphere/binding/binding.proto";

message GetUserRequest {
  // URI path parameter
  int64 user_id = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
  
  // Query parameters
  repeated string fields = 2 [(sphere.binding.location) = BINDING_LOCATION_QUERY];
}

message UpdateUserRequest {
  // URI path parameter
  int64 user_id = 1 [(sphere.binding.location) = BINDING_LOCATION_URI];
  
  // Request body
  User user = 2; // Default: BINDING_LOCATION_BODY
}
```

### Advanced Example with Custom Tags

```protobuf
```protobuf
message SearchUsersRequest {
  // Query parameters with validation and custom tags
  string query = 1 [
    (sphere.binding.location) = BINDING_LOCATION_QUERY,
    (sphere.binding.tags) = "validate:\"required\""
  ];
  int32 limit = 2 [(sphere.binding.location) = BINDING_LOCATION_QUERY];
  string token = 3 [
    (sphere.binding.location) = BINDING_LOCATION_QUERY,
    (sphere.binding.tags) = "validate:\"required\" custom:\"auth_token\""
  ];
}

message DatabaseModel {
  option (sphere.binding.default_auto_tags) = "db";
  option (sphere.binding.default_auto_tags) = "json";
  
  string name = 1;      // Will get: db:"name" json:"name"
  int64 id = 2;         // Will get: db:"id" json:"id"
  string email = 3;     // Will get: db:"email" json:"email"
}
```

## Integration with buf

Add this dependency to your `buf.yaml`:

```yaml
version: v2
deps:
  - buf.build/go-sphere/binding
```

Configure code generation in `buf.gen.yaml`:

```yaml
version: v2
managed:
  enabled: true
  disable:
    - file_option: go_package_prefix
      module: buf.build/go-sphere/binding
plugins:
  - local: protoc-gen-sphere-binding
    out: api
    opt: paths=source_relative
```

## Generated Code Usage

When used with `protoc-gen-sphere-binding`, the following Go code is generated:

```go
type GetUserRequest struct {
    UserId int64    `protobuf:"varint,1,opt,name=user_id" json:"-" uri:"user_id"`
    Fields []string `protobuf:"bytes,2,rep,name=fields" json:"-" form:"fields"`
}

// Direct usage with HTTP frameworks
func GetUser(c *gin.Context) {
    var req GetUserRequest
    if err := c.ShouldBindUri(&req); err != nil {
        // Handle error
    }
    if err := c.ShouldBindQuery(&req); err != nil {
        // Handle error
    }
}
```

## Binding Configuration Options

### Field Level Options

- `location`: Specifies the binding location (QUERY, URI, BODY)
- `tags`: Adds custom struct tags to the field
- `auto_tags`: Adds automatic tags to the field

### Message Level Options

- `default_location`: Sets default binding location for all fields in a message
- `default_auto_tags`: Adds default tags to all fields in a message

### Oneof Level Options

- `default_oneof_location`: Sets default binding location for oneof fields
- `default_oneof_auto_tags`: Adds default tags to oneof fields

## Best Practices

1. **Use meaningful field names**: Field names become tag values, so use clear, descriptive names
2. **Set appropriate binding locations**: Choose the right location (query, URI, body) for each field
3. **Validate input parameters**: Always validate path and query parameters since they come from untrusted input
4. **Use consistent naming**: Use consistent naming patterns across your API
5. **Leverage default behaviors**: Use default auto tags for common patterns like database models

## Common Binding Locations

- `BINDING_LOCATION_QUERY`: For filtering, pagination, and optional parameters
- `BINDING_LOCATION_URI`: For resource identifiers in the URL path
- `BINDING_LOCATION_BODY`: For complex data structures and create/update operations
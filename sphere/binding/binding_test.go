package binding_test

import (
	"testing"

	"github.com/go-sphere/binding/sphere/binding"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestBindingLocationDescriptor(t *testing.T) {
	descriptor := binding.BindingLocation(0).Descriptor()
	tests := []struct {
		name   protoreflect.Name
		number protoreflect.EnumNumber
	}{
		{name: "BINDING_LOCATION_UNSPECIFIED", number: 0},
		{name: "BINDING_LOCATION_QUERY", number: 1},
		{name: "BINDING_LOCATION_URI", number: 2},
		{name: "BINDING_LOCATION_JSON", number: 3},
		{name: "BINDING_LOCATION_FORM", number: 4},
		{name: "BINDING_LOCATION_HEADER", number: 5},
	}
	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			value := descriptor.Values().ByName(tt.name)
			if value == nil {
				t.Fatalf("enum value %q not found", tt.name)
			}
			if got := value.Number(); got != tt.number {
				t.Errorf("number = %d, want %d", got, tt.number)
			}
		})
	}
	if got, want := descriptor.Values().Len(), len(tests); got != want {
		t.Errorf("value count = %d, want %d", got, want)
	}
}

func TestExtensionDescriptors(t *testing.T) {
	tests := []struct {
		name        string
		extension   protoreflect.ExtensionType
		fullName    protoreflect.FullName
		number      protoreflect.FieldNumber
		kind        protoreflect.Kind
		cardinality protoreflect.Cardinality
		extendee    protoreflect.FullName
	}{
		{"default_location", binding.E_DefaultLocation, "sphere.binding.default_location", 136655300, protoreflect.EnumKind, protoreflect.Optional, "google.protobuf.MessageOptions"},
		{"default_auto_tags", binding.E_DefaultAutoTags, "sphere.binding.default_auto_tags", 136655301, protoreflect.StringKind, protoreflect.Repeated, "google.protobuf.MessageOptions"},
		{"default_oneof_location", binding.E_DefaultOneofLocation, "sphere.binding.default_oneof_location", 136655310, protoreflect.EnumKind, protoreflect.Optional, "google.protobuf.OneofOptions"},
		{"default_oneof_auto_tags", binding.E_DefaultOneofAutoTags, "sphere.binding.default_oneof_auto_tags", 136655311, protoreflect.StringKind, protoreflect.Repeated, "google.protobuf.OneofOptions"},
		{"location", binding.E_Location, "sphere.binding.location", 136655320, protoreflect.EnumKind, protoreflect.Optional, "google.protobuf.FieldOptions"},
		{"tags", binding.E_Tags, "sphere.binding.tags", 136655321, protoreflect.StringKind, protoreflect.Repeated, "google.protobuf.FieldOptions"},
		{"auto_tags", binding.E_AutoTags, "sphere.binding.auto_tags", 136655322, protoreflect.StringKind, protoreflect.Repeated, "google.protobuf.FieldOptions"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			descriptor := tt.extension.TypeDescriptor()
			if got := descriptor.FullName(); got != tt.fullName {
				t.Errorf("full name = %q, want %q", got, tt.fullName)
			}
			if got := descriptor.Number(); got != tt.number {
				t.Errorf("number = %d, want %d", got, tt.number)
			}
			if got := descriptor.Kind(); got != tt.kind {
				t.Errorf("kind = %s, want %s", got, tt.kind)
			}
			if got := descriptor.Cardinality(); got != tt.cardinality {
				t.Errorf("cardinality = %s, want %s", got, tt.cardinality)
			}
			if got := descriptor.ContainingMessage().FullName(); got != tt.extendee {
				t.Errorf("extendee = %q, want %q", got, tt.extendee)
			}
		})
	}
}

func TestExtensionsWireRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input proto.Message
		new   func() proto.Message
	}{
		{
			name: "message options",
			input: messageOptions(
				binding.BindingLocation_BINDING_LOCATION_QUERY,
				[]string{"json", "query"},
			),
			new: func() proto.Message { return &descriptorpb.MessageOptions{} },
		},
		{
			name: "oneof options",
			input: oneofOptions(
				binding.BindingLocation_BINDING_LOCATION_FORM,
				[]string{"form"},
			),
			new: func() proto.Message { return &descriptorpb.OneofOptions{} },
		},
		{
			name: "field options",
			input: fieldOptions(
				binding.BindingLocation_BINDING_LOCATION_HEADER,
				[]string{`validate:"required"`},
				[]string{"header"},
			),
			new: func() proto.Message { return &descriptorpb.FieldOptions{} },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := proto.Marshal(tt.input)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			output := tt.new()
			if err := proto.Unmarshal(data, output); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if !proto.Equal(output, tt.input) {
				t.Errorf("round trip mismatch:\n got: %v\nwant: %v", output, tt.input)
			}
		})
	}
}

func messageOptions(location binding.BindingLocation, tags []string) *descriptorpb.MessageOptions {
	options := &descriptorpb.MessageOptions{}
	proto.SetExtension(options, binding.E_DefaultLocation, location)
	proto.SetExtension(options, binding.E_DefaultAutoTags, tags)
	return options
}

func oneofOptions(location binding.BindingLocation, tags []string) *descriptorpb.OneofOptions {
	options := &descriptorpb.OneofOptions{}
	proto.SetExtension(options, binding.E_DefaultOneofLocation, location)
	proto.SetExtension(options, binding.E_DefaultOneofAutoTags, tags)
	return options
}

func fieldOptions(location binding.BindingLocation, tags, autoTags []string) *descriptorpb.FieldOptions {
	options := &descriptorpb.FieldOptions{}
	proto.SetExtension(options, binding.E_Location, location)
	proto.SetExtension(options, binding.E_Tags, tags)
	proto.SetExtension(options, binding.E_AutoTags, autoTags)
	return options
}

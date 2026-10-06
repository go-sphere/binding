package binding_test

import (
	"fmt"

	"github.com/go-sphere/binding/sphere/binding"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// A code generator receives these options from a field descriptor. Here they
// are built by hand to show how each extension is read.
func Example_fieldOptions() {
	opts := &descriptorpb.FieldOptions{}
	proto.SetExtension(opts, binding.E_Location, binding.BindingLocation_BINDING_LOCATION_HEADER)
	proto.SetExtension(opts, binding.E_Tags, []string{`validate:"required"`})

	if proto.HasExtension(opts, binding.E_Location) {
		loc, _ := proto.GetExtension(opts, binding.E_Location).(binding.BindingLocation)
		fmt.Println(loc)
	}
	tags, _ := proto.GetExtension(opts, binding.E_Tags).([]string)
	fmt.Println(tags)
	// Output:
	// BINDING_LOCATION_HEADER
	// [validate:"required"]
}

func Example_messageDefaults() {
	opts := &descriptorpb.MessageOptions{}
	proto.SetExtension(opts, binding.E_DefaultLocation, binding.BindingLocation_BINDING_LOCATION_QUERY)
	proto.SetExtension(opts, binding.E_DefaultAutoTags, []string{"form"})

	loc, _ := proto.GetExtension(opts, binding.E_DefaultLocation).(binding.BindingLocation)
	autoTags, _ := proto.GetExtension(opts, binding.E_DefaultAutoTags).([]string)
	fmt.Println(loc, autoTags)

	// An unset extension reports false from HasExtension and the zero value
	// from GetExtension.
	empty := &descriptorpb.MessageOptions{}
	fmt.Println(proto.HasExtension(empty, binding.E_DefaultLocation), proto.GetExtension(empty, binding.E_DefaultLocation))
	// Output:
	// BINDING_LOCATION_QUERY [form]
	// false BINDING_LOCATION_UNSPECIFIED
}

func ExampleBindingLocation() {
	loc := binding.BindingLocation(binding.BindingLocation_value["BINDING_LOCATION_URI"])
	fmt.Println(loc == binding.BindingLocation_BINDING_LOCATION_URI, int32(loc), loc.String())
	// Output: true 2 BINDING_LOCATION_URI
}

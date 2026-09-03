package binding_test

import (
	"testing"

	"github.com/go-sphere/binding/sphere/binding"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestBindingLocation_Values(t *testing.T) {
	tests := []struct {
		loc  binding.BindingLocation
		name string
		num  int32
	}{
		{binding.BindingLocation_BINDING_LOCATION_UNSPECIFIED, "BINDING_LOCATION_UNSPECIFIED", 0},
		{binding.BindingLocation_BINDING_LOCATION_QUERY, "BINDING_LOCATION_QUERY", 1},
		{binding.BindingLocation_BINDING_LOCATION_URI, "BINDING_LOCATION_URI", 2},
		{binding.BindingLocation_BINDING_LOCATION_JSON, "BINDING_LOCATION_JSON", 3},
		{binding.BindingLocation_BINDING_LOCATION_FORM, "BINDING_LOCATION_FORM", 4},
		{binding.BindingLocation_BINDING_LOCATION_HEADER, "BINDING_LOCATION_HEADER", 5},
	}

	for _, tt := range tests {
		if tt.loc.String() != tt.name {
			t.Errorf("expected string %s, got %s", tt.name, tt.loc.String())
		}
		if int32(tt.loc.Number()) != tt.num {
			t.Errorf("expected number %d, got %d", tt.num, tt.loc.Number())
		}
		if binding.BindingLocation_value[tt.name] != tt.num {
			t.Errorf("value map mismatch for %s", tt.name)
		}
		if binding.BindingLocation_name[tt.num] != tt.name {
			t.Errorf("name map mismatch for %d", tt.num)
		}
	}
}

func TestBindingExtensions_MessageOptions(t *testing.T) {
	opts := &descriptorpb.MessageOptions{}
	proto.SetExtension(opts, binding.E_DefaultLocation, binding.BindingLocation_BINDING_LOCATION_QUERY)
	proto.SetExtension(opts, binding.E_DefaultAutoTags, []string{"json", "form"})

	if !proto.HasExtension(opts, binding.E_DefaultLocation) {
		t.Fatal("expected DefaultLocation extension to be present")
	}
	loc := proto.GetExtension(opts, binding.E_DefaultLocation).(binding.BindingLocation)
	if loc != binding.BindingLocation_BINDING_LOCATION_QUERY {
		t.Errorf("expected QUERY, got %v", loc)
	}

	tags := proto.GetExtension(opts, binding.E_DefaultAutoTags).([]string)
	if len(tags) != 2 || tags[0] != "json" || tags[1] != "form" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestBindingExtensions_OneofOptions(t *testing.T) {
	opts := &descriptorpb.OneofOptions{}
	proto.SetExtension(opts, binding.E_DefaultOneofLocation, binding.BindingLocation_BINDING_LOCATION_FORM)
	proto.SetExtension(opts, binding.E_DefaultOneofAutoTags, []string{"form"})

	if !proto.HasExtension(opts, binding.E_DefaultOneofLocation) {
		t.Fatal("expected DefaultOneofLocation extension")
	}
	loc := proto.GetExtension(opts, binding.E_DefaultOneofLocation).(binding.BindingLocation)
	if loc != binding.BindingLocation_BINDING_LOCATION_FORM {
		t.Errorf("expected FORM, got %v", loc)
	}

	tags := proto.GetExtension(opts, binding.E_DefaultOneofAutoTags).([]string)
	if len(tags) != 1 || tags[0] != "form" {
		t.Errorf("unexpected auto tags: %v", tags)
	}
}

func TestBindingExtensions_FieldOptions(t *testing.T) {
	opts := &descriptorpb.FieldOptions{}
	proto.SetExtension(opts, binding.E_Location, binding.BindingLocation_BINDING_LOCATION_URI)
	proto.SetExtension(opts, binding.E_Tags, []string{`validate:"required"`})
	proto.SetExtension(opts, binding.E_AutoTags, []string{"param"})

	if !proto.HasExtension(opts, binding.E_Location) {
		t.Fatal("expected Location extension")
	}
	loc := proto.GetExtension(opts, binding.E_Location).(binding.BindingLocation)
	if loc != binding.BindingLocation_BINDING_LOCATION_URI {
		t.Errorf("expected URI, got %v", loc)
	}

	tags := proto.GetExtension(opts, binding.E_Tags).([]string)
	if len(tags) != 1 || tags[0] != `validate:"required"` {
		t.Errorf("unexpected tags: %v", tags)
	}

	autoTags := proto.GetExtension(opts, binding.E_AutoTags).([]string)
	if len(autoTags) != 1 || autoTags[0] != "param" {
		t.Errorf("unexpected auto tags: %v", autoTags)
	}
}

func TestBindingExtensions_SerializationRoundTrip(t *testing.T) {
	orig := &descriptorpb.FieldOptions{}
	proto.SetExtension(orig, binding.E_Location, binding.BindingLocation_BINDING_LOCATION_HEADER)
	proto.SetExtension(orig, binding.E_Tags, []string{`header:"X-Trace-Id"`})

	data, err := proto.Marshal(orig)
	if err != nil {
		t.Fatalf("proto.Marshal failed: %v", err)
	}

	parsed := &descriptorpb.FieldOptions{}
	if err := proto.Unmarshal(data, parsed); err != nil {
		t.Fatalf("proto.Unmarshal failed: %v", err)
	}

	if !proto.HasExtension(parsed, binding.E_Location) {
		t.Fatal("expected Location extension after unmarshal")
	}
	loc := proto.GetExtension(parsed, binding.E_Location).(binding.BindingLocation)
	if loc != binding.BindingLocation_BINDING_LOCATION_HEADER {
		t.Errorf("expected HEADER, got %v", loc)
	}
}

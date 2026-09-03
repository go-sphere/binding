package binding_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/go-sphere/binding/sphere/binding"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// TestBinding_NilReceivers tests nil descriptor pointer safety when interacting
// with protobuf extension helpers.
func TestBinding_NilReceivers(t *testing.T) {
	var nilMsgOpts *descriptorpb.MessageOptions
	if proto.HasExtension(nilMsgOpts, binding.E_DefaultLocation) {
		t.Error("expected HasExtension on nil MessageOptions to be false")
	}
	if proto.HasExtension(nilMsgOpts, binding.E_DefaultAutoTags) {
		t.Error("expected HasExtension on nil MessageOptions to be false")
	}

	var nilOneofOpts *descriptorpb.OneofOptions
	if proto.HasExtension(nilOneofOpts, binding.E_DefaultOneofLocation) {
		t.Error("expected HasExtension on nil OneofOptions to be false")
	}
	if proto.HasExtension(nilOneofOpts, binding.E_DefaultOneofAutoTags) {
		t.Error("expected HasExtension on nil OneofOptions to be false")
	}

	var nilFieldOpts *descriptorpb.FieldOptions
	if proto.HasExtension(nilFieldOpts, binding.E_Location) {
		t.Error("expected HasExtension on nil FieldOptions to be false")
	}
	if proto.HasExtension(nilFieldOpts, binding.E_Tags) {
		t.Error("expected HasExtension on nil FieldOptions to be false")
	}
	if proto.HasExtension(nilFieldOpts, binding.E_AutoTags) {
		t.Error("expected HasExtension on nil FieldOptions to be false")
	}
}

// TestBindingLocation_Boundaries tests enum behaviour at and beyond boundaries.
func TestBindingLocation_Boundaries(t *testing.T) {
	invalidValues := []int32{-100, -1, 6, 7, 100, 9999}
	for _, v := range invalidValues {
		loc := binding.BindingLocation(v)
		str := loc.String()
		if str == "" {
			t.Errorf("expected non-empty string representation for enum value %d", v)
		}
		if int32(loc.Number()) != v {
			t.Errorf("expected Number() to return %d, got %d", v, loc.Number())
		}
	}
}

// TestBinding_ConcurrentRoundtripStress tests wire serialization roundtrips
// across 50 concurrent goroutines under race detection.
func TestBinding_ConcurrentRoundtripStress(t *testing.T) {
	const goroutines = 50
	const iterations = 20

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				// 1. Message options
				msgOpts := &descriptorpb.MessageOptions{}
				loc := binding.BindingLocation(id % 6)
				proto.SetExtension(msgOpts, binding.E_DefaultLocation, loc)
				tags := []string{fmt.Sprintf("tag_%d_%d", id, i), "json", "form"}
				proto.SetExtension(msgOpts, binding.E_DefaultAutoTags, tags)

				data, err := proto.Marshal(msgOpts)
				if err != nil {
					t.Errorf("goroutine %d: proto.Marshal msgOpts failed: %v", id, err)
					return
				}

				parsedMsgOpts := &descriptorpb.MessageOptions{}
				if err := proto.Unmarshal(data, parsedMsgOpts); err != nil {
					t.Errorf("goroutine %d: proto.Unmarshal msgOpts failed: %v", id, err)
					return
				}

				if !proto.HasExtension(parsedMsgOpts, binding.E_DefaultLocation) {
					t.Errorf("goroutine %d: missing DefaultLocation extension", id)
					return
				}
				gotLoc := proto.GetExtension(parsedMsgOpts, binding.E_DefaultLocation).(binding.BindingLocation)
				if gotLoc != loc {
					t.Errorf("goroutine %d: location mismatch: got %v, want %v", id, gotLoc, loc)
					return
				}

				// 2. Field options
				fieldOpts := &descriptorpb.FieldOptions{}
				fieldLoc := binding.BindingLocation((id + 1) % 6)
				proto.SetExtension(fieldOpts, binding.E_Location, fieldLoc)
				proto.SetExtension(fieldOpts, binding.E_Tags, []string{`validate:"required"`, `header:"X-Req"`})
				proto.SetExtension(fieldOpts, binding.E_AutoTags, []string{"param", "query"})

				fData, err := proto.Marshal(fieldOpts)
				if err != nil {
					t.Errorf("goroutine %d: proto.Marshal fieldOpts failed: %v", id, err)
					return
				}

				parsedFieldOpts := &descriptorpb.FieldOptions{}
				if err := proto.Unmarshal(fData, parsedFieldOpts); err != nil {
					t.Errorf("goroutine %d: proto.Unmarshal fieldOpts failed: %v", id, err)
					return
				}

				if !proto.Equal(fieldOpts, parsedFieldOpts) {
					t.Errorf("goroutine %d: fieldOpts not equal after roundtrip", id)
					return
				}

				// 3. Clone verification
				cloned := proto.Clone(fieldOpts).(*descriptorpb.FieldOptions)
				if !proto.Equal(cloned, fieldOpts) {
					t.Errorf("goroutine %d: cloned fieldOpts mismatch", id)
					return
				}
			}
		}(g)
	}

	wg.Wait()
}

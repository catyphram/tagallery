package util_test

import (
	"testing"

	"tagallery.com/api/util"
)

func TestContainsString(t *testing.T) {
	if contains := util.ContainsString([]string{"123", "456"}, "abc", true); contains {
		t.Error("Slice should not contain string.")
	}
	if contains := util.ContainsString([]string{"123", "abc"}, "abc", true); !contains {
		t.Error("Slice should contain string.")
	}
	if contains := util.ContainsString([]string{"123", "abc"}, "aBc", false); !contains {
		t.Error("ContainsString() should compare case insensitive.")
	}
}

func TestIntPtr(t *testing.T) {
	if i := util.IntPtr(1); *i != 1 {
		t.Error("IntPtr() should return pointer with the correct value.")
	}
}

func TestStringPtr(t *testing.T) {
	if i := util.StringPtr("test"); *i != "test" {
		t.Error("StringPtr() should return pointer with the correct value.")
	}
}

package form

import (
	"net/url"
	"strings"
	"testing"
)

type formValue string

func (v *formValue) UnmarshalForm(values []string) error {
	*v = formValue(strings.Join(values, ","))
	return nil
}

func TestUnmarshal(t *testing.T) {
	type formValues struct {
		Custom  formValue `form:"custom"`
		Default string    `form:"default"`
	}

	var got formValues
	err := Unmarshal(url.Values{
		"custom":  {"first", "second"},
		"default": {"value"},
	}, &got)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Custom != "first,second" {
		t.Fatalf("expected custom value %q, got %q", "first,second", got.Custom)
	}
	if got.Default != "value" {
		t.Fatalf("expected default value %q, got %q", "value", got.Default)
	}
}

func TestUnmarshalSkipsUnexportedFields(t *testing.T) {
	type formValues struct {
		Visible string `form:"visible"`
		hidden  string `form:"hidden"`
	}

	got := formValues{hidden: "original"}
	err := Unmarshal(url.Values{
		"visible": {"decoded"},
		"hidden":  {"must-not-be-decoded"},
	}, &got)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Visible != "decoded" {
		t.Fatalf("expected visible field %q, got %q", "decoded", got.Visible)
	}
	if got.hidden != "original" {
		t.Fatalf("expected unexported field to remain unchanged, got %q", got.hidden)
	}
}

package form

import (
	"mime/multipart"
	"testing"
)

func TestUnmarshalFiles(t *testing.T) {
	avatar := &multipart.FileHeader{Filename: "avatar.png"}
	firstAttachment := &multipart.FileHeader{Filename: "first.txt"}
	secondAttachment := &multipart.FileHeader{Filename: "second.txt"}

	type formFiles struct {
		Avatar      *multipart.FileHeader   `form:"avatar"`
		Attachments []*multipart.FileHeader `form:"attachments"`
		Ignored     *multipart.FileHeader   `form:"-"`
	}

	var got formFiles
	err := UnmarshalFiles(map[string][]*multipart.FileHeader{
		"avatar":      {avatar},
		"attachments": {firstAttachment, secondAttachment},
		"-":           {{Filename: "ignored.txt"}},
	}, &got)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Avatar != avatar {
		t.Fatalf("expected avatar file %q, got %#v", avatar.Filename, got.Avatar)
	}
	if len(got.Attachments) != 2 || got.Attachments[0] != firstAttachment || got.Attachments[1] != secondAttachment {
		t.Fatalf("unexpected attachments: %#v", got.Attachments)
	}
	if got.Ignored != nil {
		t.Fatalf("expected ignored file field to remain nil, got %#v", got.Ignored)
	}
}

func TestUnmarshalFilesSkipsUnexportedFields(t *testing.T) {
	visible := &multipart.FileHeader{Filename: "visible.txt"}
	originalHidden := &multipart.FileHeader{Filename: "original.txt"}

	type formFiles struct {
		Visible *multipart.FileHeader `form:"visible"`
		hidden  *multipart.FileHeader `form:"hidden"`
	}

	got := formFiles{hidden: originalHidden}
	err := UnmarshalFiles(map[string][]*multipart.FileHeader{
		"visible": {visible},
		"hidden":  {{Filename: "must-not-be-decoded.txt"}},
	}, &got)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Visible != visible {
		t.Fatalf("expected visible file %q, got %#v", visible.Filename, got.Visible)
	}
	if got.hidden != originalHidden {
		t.Fatalf("expected unexported file field to remain unchanged, got %#v", got.hidden)
	}
}

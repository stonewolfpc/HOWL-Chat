package worldstore

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestNormalizeMediaManifestClassifiesMalformedAndUnsupported(t *testing.T) {
	cases := []struct {
		name string
		data string
		want error
	}{
		{
			name: "malformed missing shape",
			data: `{"format":"howl-media-refs-v1"}`,
			want: ErrMediaMalformed,
		},
		{
			name: "malformed item",
			data: `{"format":"howl-media-refs-v1","items":["not-an-object"]}`,
			want: ErrMediaMalformed,
		},
		{
			name: "unsupported format",
			data: `{"format":"howl-media-refs-v2","items":[]}`,
			want: ErrMediaUnsupported,
		},
		{
			name: "unsupported role",
			data: `{"format":"howl-media-refs-v1","items":[{"mediaId":"m1","role":"hero"}]}`,
			want: ErrMediaUnsupported,
		},
		{
			name: "inline media",
			data: `{"format":"howl-media-refs-v1","items":[],"dataUrl":"data:image/png;base64,AAAA"}`,
			want: ErrMediaInline,
		},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := NormalizeMediaManifest([]byte(test.data))
			if !errors.Is(err, test.want) {
				t.Fatalf("error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestNormalizeMediaManifestPreservesUnknownFieldsInBothShapes(t *testing.T) {
	cases := []struct {
		name string
		data string
		key  string
	}{
		{
			name: "canonical",
			data: `{"format":"howl-media-refs-v1","items":[{"mediaId":"m1","role":"thumbnail","custom":"keep"}],"unknownTop":"keep"}`,
			key:  "items",
		},
		{
			name: "legacy",
			data: `{"images":[{"id":"m1","custom":"keep"}],"unknownTop":"keep"}`,
			key:  "images",
		},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			normalized, err := NormalizeMediaManifest([]byte(test.data))
			if err != nil {
				t.Fatalf("normalize: %v", err)
			}
			var doc map[string]interface{}
			if err := json.Unmarshal(normalized, &doc); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if doc["unknownTop"] != "keep" {
				t.Fatalf("top-level unknown field lost: %+v", doc)
			}
			items := doc[test.key].([]interface{})
			if items[0].(map[string]interface{})["custom"] != "keep" {
				t.Fatalf("nested unknown field lost: %+v", items[0])
			}
		})
	}
}

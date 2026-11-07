package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTemplateUploadPublishesByDefault(t *testing.T) {
	for _, kind := range []string{"video", "suite"} {
		t.Run(kind, func(t *testing.T) {
			var calls []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "Bearer test" {
					t.Error("missing authentication")
				}
				calls = append(calls, r.URL.Path)
				if r.Method != http.MethodPost {
					t.Errorf("method = %s", r.Method)
				}
				if r.URL.Path == "/template-publications" {
					var body map[string]string
					json.NewDecoder(r.Body).Decode(&body)
					if body["kind"] != strings.ToUpper(kind) || body["templateId"] != "saved-id" {
						t.Errorf("publication = %#v", body)
					}
					w.Write([]byte(`{"id":"publication","status":"APPROVED"}`))
				} else {
					w.Write([]byte(`{"id":"saved-id"}`))
				}
			}))
			defer server.Close()
			t.Setenv("VIMOX_HIT_TOKEN", "test")
			t.Setenv("VIMOX_HIT_API_URL", server.URL)
			file := filepath.Join(t.TempDir(), "template.json")
			os.WriteFile(file, []byte(`{"name":"example"}`), 0600)
			if err := templateCommand([]string{"upload", "--kind", kind, "--file", file}); err != nil {
				t.Fatal(err)
			}
			if len(calls) != 2 {
				t.Fatalf("calls = %#v", calls)
			}
			calls = nil
			if err := templateCommand([]string{"upload", "--kind", kind, "--file", file, "--private"}); err != nil {
				t.Fatal(err)
			}
			if len(calls) != 1 {
				t.Fatalf("private upload calls = %#v", calls)
			}
		})
	}
}

func TestTemplateMediaPreservesOrderAndContentTypes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Error(err)
			return
		}
		if r.FormValue("metadata") != `{"sourceKey":"test"}` {
			t.Error("metadata missing")
		}
		files := r.MultipartForm.File["media"]
		if len(files) != 2 || files[0].Header.Get("Content-Type") != "video/mp4" || files[1].Header.Get("Content-Type") != "image/png" {
			t.Errorf("media = %#v", files)
		}
		w.Write([]byte(`{"id":"saved"}`))
	}))
	defer server.Close()
	t.Setenv("VIMOX_HIT_TOKEN", "test")
	t.Setenv("VIMOX_HIT_API_URL", server.URL)
	dir := t.TempDir()
	video, image := filepath.Join(dir, "source.mp4"), filepath.Join(dir, "dna.png")
	os.WriteFile(video, []byte("video bytes"), 0600)
	os.WriteFile(image, []byte("image bytes"), 0600)
	if _, err := uploadVideoTemplate([]byte(`{"sourceKey":"test"}`), []string{video, image}); err != nil {
		t.Fatal(err)
	}
}

func TestTemplateCommandRejectsInvalidInput(t *testing.T) {
	for _, args := range [][]string{{"upload"}, {"stats"}, {"list", "--kind", "bad"}, {"list", "--origin", "all"}, {"list", "--limit", "0"}} {
		if err := templateCommand(args); err == nil {
			t.Errorf("accepted %#v", args)
		}
	}
}

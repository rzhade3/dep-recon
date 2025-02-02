package src

import (
	"reflect"
	"testing"

	"github.com/rzhade3/dep-recon/src/manifest"
)

func TestValidateManifestFilepath(t *testing.T) {
	lang, err := ValidateManifestFilepath("package.json")
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
	if reflect.TypeOf(lang) != reflect.TypeOf(manifest.Javascript{}) {
		t.Errorf("Expected Javascript, got %s", lang)
	}

	lang, err = ValidateManifestFilepath("Cargo.toml")
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
	if reflect.TypeOf(lang) != reflect.TypeOf(manifest.Rust{}) {
		t.Errorf("Expected Rust, got %s", lang)
	}

	lang, err = ValidateManifestFilepath("Gemfile")
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
	if reflect.TypeOf(lang) != reflect.TypeOf(manifest.Ruby{}) {
		t.Errorf("Expected Ruby, got %s", lang)
	}

	lang, err = ValidateManifestFilepath("go.mod")
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
	if reflect.TypeOf(lang) != reflect.TypeOf(manifest.Golang{}) {
		t.Errorf("Expected Golang, got %s", lang)
	}

	lang, err = ValidateManifestFilepath("pom.xml")
	if err == nil {
		t.Error("Expected err, got nil")
	}
	if lang != nil {
		t.Errorf("Expected nil, got %s", lang)
	}

	lang, err = ValidateManifestFilepath("manifest.yml")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if lang != nil {
		t.Errorf("Expected nil, got %s", lang)
	}
}

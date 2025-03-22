package core

import (
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveImage2FileEncodesAndValidatesInputs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "image.png")
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	if err := SaveImage2File(img, path, png.Encode); err != nil {
		t.Fatalf("SaveImage2File() error = %v", err)
	}
	encoded, err := os.ReadFile(path)
	if err != nil || len(encoded) == 0 {
		t.Fatalf("encoded output unavailable: len=%d err=%v", len(encoded), err)
	}
	if err := SaveImage2File(img, filepath.Join(dir, "nil-encoder"), nil); err == nil {
		t.Fatal("expected nil encoder error")
	}
}

func TestBase64RoundTrip(t *testing.T) {
	want := []byte{0, 1, 2, 127, 128, 254, 255}
	encoded := Bytes2Base64(want)
	got, err := Base642Bytes(encoded)
	if err != nil {
		t.Fatalf("Base642Bytes() error = %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("decoded bytes = %v, want %v", got, want)
	}
	if _, err := Base642Bytes("not base64!"); err == nil {
		t.Fatal("expected malformed base64 error")
	}
}

func TestSaveValToFileUsesDefaultVariableNameAndRoundTripsContent(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "payload.bin")
	content := []byte("binary\x00payload\xff")
	if err := SaveValToFile(inputPath, "", content); err != nil {
		t.Fatalf("SaveValToFile() error = %v", err)
	}

	outPath := filepath.Join(dir, "payload.bin.go")
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	text := string(data)
	if strings.Contains(text, "var  =") {
		t.Fatalf("generated source has an empty variable name: %q", text)
	}
	const prefix = `var payload_bin = "`
	start := strings.Index(text, prefix)
	if start < 0 || !strings.HasSuffix(text, "\"\n") {
		t.Fatalf("generated source has unexpected format: %q", text)
	}
	encoded := strings.TrimSuffix(strings.TrimPrefix(text[start:], prefix), "\"\n")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode generated value: %v", err)
	}
	if string(decoded) != string(content) {
		t.Fatalf("generated content = %q, want %q", decoded, content)
	}
}

func TestPersistFileReadsAndPersistsActualBytes(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "font.dat")
	content := []byte("font bytes\x00\x01\xfe")
	if err := os.WriteFile(inputPath, content, 0600); err != nil {
		t.Fatalf("write input: %v", err)
	}
	if err := PersistFile(inputPath, "EmbeddedFont"); err != nil {
		t.Fatalf("PersistFile() error = %v", err)
	}

	data, err := os.ReadFile(inputPath + ".go")
	if err != nil {
		t.Fatalf("read persisted file: %v", err)
	}
	text := string(data)
	const prefix = `var EmbeddedFont = "`
	start := strings.Index(text, prefix)
	if start < 0 || !strings.HasSuffix(text, "\"\n") {
		t.Fatalf("persisted source has unexpected format: %q", text)
	}
	encoded := strings.TrimSuffix(strings.TrimPrefix(text[start:], prefix), "\"\n")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode persisted value: %v", err)
	}
	if string(decoded) != string(content) {
		t.Fatalf("persisted content = %q, want %q", decoded, content)
	}
}

func TestSaveValToFileSanitizesVariableNamesWithoutCollisions(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"1a-b", "2a-b", "type"} {
		if err := SaveValToFile(filepath.Join(dir, name+".bin"), name, []byte(name)); err != nil {
			t.Fatalf("SaveValToFile(%q): %v", name, err)
		}
	}
	first, err := os.ReadFile(filepath.Join(dir, "1a-b.bin.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(first), "var _1a_b =") {
		t.Fatalf("leading digit was not preserved in identifier: %q", first)
	}
	reserved, err := os.ReadFile(filepath.Join(dir, "type.bin.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(reserved), "var type_ =") {
		t.Fatalf("reserved identifier was not suffixed: %q", reserved)
	}
}

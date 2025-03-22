package core

import (
	"encoding/base64"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var PersistStr = `package main
var %s = "%s"
`

// Bytes2Base64 converts bytes to base64 string
func Bytes2Base64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Base642Bytes converts base64 string to bytes
func Base642Bytes(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// SaveValToFile saves base64 string to file
func SaveValToFile(fileName, valName string, content []byte) error {
	if strings.TrimSpace(fileName) == "" {
		return fmt.Errorf("file name is empty")
	}
	if valName == "" {
		base := filepath.Base(fileName)
		// Keep the extension in the source name (sanitized below) so two files
		// such as icon.svg and icon.json do not silently collide.
		valName = base
	}
	valName = goIdentifier(valName)
	stData := Bytes2Base64(content)
	str := fmt.Sprintf(PersistStr, valName, stData)

	outputPath := filepath.Join(filepath.Dir(fileName), filepath.Base(fileName)+".go")
	if err := os.WriteFile(outputPath, []byte(str), 0o644); err != nil {
		return fmt.Errorf("write generated file %q: %w", outputPath, err)
	}
	return nil
}

// PersistFile persists file to base64 string
func PersistFile(fileName, valName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("%w, error reading font file", err)
	}

	// save to file
	if err = SaveValToFile(fileName, valName, data); err != nil {
		return fmt.Errorf("%w, error saving font file", err)
	}

	return nil
}

// goIdentifier turns a file- or user-supplied name into a valid Go variable
// identifier.  Persisted resources are source code, so silently emitting
// "var  =" or a name containing a dot would create an unusable generated file.
func goIdentifier(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "data"
	}

	var b strings.Builder
	for i, r := range []rune(name) {
		if r == '_' || unicode.IsLetter(r) || (i > 0 && unicode.IsDigit(r)) {
			b.WriteRune(r)
			continue
		}
		// Preserve a leading digit after prefixing an underscore; replacing it
		// outright would make distinct names such as 1a and 2a collide.
		if i == 0 && unicode.IsDigit(r) {
			b.WriteByte('_')
			b.WriteRune(r)
			continue
		}
		b.WriteByte('_')
	}
	result := b.String()
	if result == "" {
		result = "data"
	}
	if token.IsKeyword(result) {
		result += "_"
	}
	return result
}

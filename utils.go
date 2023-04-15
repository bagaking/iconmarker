package iconmarker

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	stData := Bytes2Base64(content)
	str := fmt.Sprintf(PersistStr, valName, stData)

	fName := filepath.Base(fileName)
	if valName == "" {
		valName = fName
	}
	fName = filepath.Join(filepath.Dir(fileName), fName+".go")
	fmt.Println("======== export name:", fName)

	file, err := os.Create(fName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(str)
	if err != nil {
		return err
	}

	return nil
}

// PersistFile persists file to base64 string
func PersistFile(fileName, valName string) error {
	// Read file
	data := make([]byte, 0)
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("%w, error reading font file", err)
	}
	defer file.Close()

	// copy form file to data
	if _, err = io.Copy(bytes.NewBuffer(data), file); err != nil {
		return fmt.Errorf("%w, error reading font file", err)
	}

	// save to file
	if err = SaveValToFile(fileName, valName, data); err != nil {
		return fmt.Errorf("%w, error saving font file", err)
	}

	return nil
}

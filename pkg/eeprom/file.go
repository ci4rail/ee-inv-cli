package eeprom

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type EEPROMFile struct {
	file   *os.File
	offset int64
	size   int64
}

// NewEEPROMFile creates a new EEPROMFile instance
// the filePath is the path plus the offset and size in the form "path@offset:size"
func NewEEPROMFile(spec string) (*EEPROMFile, error) {

	// Parse the file path to get the offset and size
	filePath, offset, size, err := parseFilePath(spec)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	return &EEPROMFile{
		file:   file,
		offset: offset,
		size:   size,
	}, nil
}

func parseFilePath(spec string) (string, int64, int64, error) {
	// Split the spec into path and offset:size
	parts := strings.Split(spec, "@")
	if len(parts) != 2 {
		return "", 0, 0, fmt.Errorf("invalid file spec: %s", spec)
	}
	path := parts[0]
	offsetSize := strings.Split(parts[1], ":")
	if len(offsetSize) != 2 {
		return "", 0, 0, fmt.Errorf("invalid offset:size spec: %s", parts[1])
	}
	offset, err := strconv.ParseInt(offsetSize[0], 10, 64)
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid offset: %s", offsetSize[0])
	}
	size, err := strconv.ParseInt(offsetSize[1], 10, 64)
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid size: %s", offsetSize[1])
	}
	return path, offset, size, nil
}

// Read reads the EEPROM data from the file
func (e *EEPROMFile) Read() ([]byte, error) {
	data := make([]byte, e.size)

	nRead, err := e.file.ReadAt(data, e.offset)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	if int64(nRead) != e.size {
		return nil, fmt.Errorf("read %d bytes, expected %d bytes", nRead, e.size)
	}
	return data, nil
}

// Write writes the EEPROM data to the file
func (e *EEPROMFile) Write(data []byte) error {
	if int64(len(data)) > e.size {
		return fmt.Errorf("data size %d to big for EEPROM file %d", len(data), e.size)
	}
	nWritten, err := e.file.WriteAt(data, e.offset)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	if nWritten != len(data) {
		return fmt.Errorf("written %d bytes, expected %d bytes", nWritten, len(data))
	}
	return nil
}

// Close closes the file
func (e *EEPROMFile) Close() error {
	return e.file.Close()
}

package eeprom

import (

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFilePath(t *testing.T) {
	// Test case 1: Valid file path
	spec := "test.txt@512:256"
	path, offset, size, err := parseFilePath(spec)
	assert.NoError(t, err)
	assert.Equal(t, "test.txt", path)
	assert.Equal(t, int64(512), offset)
	assert.Equal(t, int64(256), size)

	// Test case 2: Invalid file spec
	spec = "test.txt@0"
	_, _, _, err = parseFilePath(spec)
	assert.Error(t, err)

	// Test case 3: Invalid offset:size spec
	spec = "test.txt@0:1024:2048"
	_, _, _, err = parseFilePath(spec)
	assert.Error(t, err)
}
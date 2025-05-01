package eeprom

import (
	"crypto/ed25519"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContent(t *testing.T) {
	// Test case 1: Valid payload
	payload := []byte("Hello, EEPROM!")
	privateKey := ed25519.NewKeyFromSeed([]byte("12345678901234567890123456789012")) // Example private key
	content, err := NewContent(payload, privateKey)
	assert.NoError(t, err)
	assert.Equal(t, MagicNumber, content.MagicNumber)
	assert.Equal(t, uint16(len(payload)), content.Len)
	assert.Equal(t, payload, content.Payload)

	// Test case 2: Payload too large
	largePayload := make([]byte, 0x10000) // 65536 bytes
	_, err = NewContent(largePayload, privateKey)
	assert.Error(t, err)
	assert.Equal(t, "payload too large: 65536 bytes", err.Error())
}

func TestMarshal(t *testing.T) {
	// Test case 1: Valid Content
	payload := []byte("Hello, EEPROM!")
	privateKey := ed25519.NewKeyFromSeed([]byte("12345678901234567890123456789012")) // Example private key
	content, err := NewContent(payload, privateKey)
	assert.NoError(t, err)

	marshaledData, err := content.Marshal()
	assert.NoError(t, err)
	assert.NotNil(t, marshaledData)
	assert.Equal(t, MagicNumber, uint16(marshaledData[0])|uint16(marshaledData[1])<<8)
	assert.Equal(t, uint16(len(payload)), uint16(marshaledData[2])|uint16(marshaledData[3])<<8)
}

func TestUnmarshal(t *testing.T) {
	// Test case 1: Valid data
	payload := []byte("Hello, EEPROM!")
	privateKey := ed25519.NewKeyFromSeed([]byte("12345678901234567890123456789012")) // Example private key
	content, err := NewContent(payload, privateKey)
	assert.NoError(t, err)

	marshaledData, err := content.Marshal()
	assert.NoError(t, err)

	unmarshaledContent, err := Unmarshal(marshaledData)
	assert.NoError(t, err)
	assert.Equal(t, content.MagicNumber, unmarshaledContent.MagicNumber)
	assert.Equal(t, content.Len, unmarshaledContent.Len)
	assert.Equal(t, content.Payload, unmarshaledContent.Payload)

	// Test case 2: Invalid len
	content.Len = 5
	marshaledData, _ = content.Marshal()
	unmarshaledContent, err = Unmarshal(marshaledData)
	assert.Error(t, err)
	assert.Nil(t, unmarshaledContent)

	marshaledData, _ = content.Marshal()
	marshaledData = marshaledData[:4] 
	unmarshaledContent, err = Unmarshal(marshaledData)
	assert.Error(t, err)
	assert.Nil(t, unmarshaledContent)
}

// Test VerifySignature
func TestVerifySignature(t *testing.T) {
	payload := []byte("Hello, EEPROM!")
	privateKey := ed25519.NewKeyFromSeed([]byte("12345678901234567890123456789012")) // Example private key
	content, err := NewContent(payload, privateKey)
	assert.NoError(t, err)

	publicKey := privateKey.Public().(ed25519.PublicKey)

	// Test case 1: Valid signature
	valid := VerifySignature(content, publicKey)
	assert.True(t, valid)

	// Test case 2: Invalid signature
	content.Payload = []byte("Tampered data")
	valid = VerifySignature(content, publicKey)
	assert.False(t, valid)
}
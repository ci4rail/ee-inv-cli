package eeprom

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"fmt"
)

const (
	MagicNumber = uint16(0xFEAD)
)

var endianess = binary.LittleEndian

// Content defines the layout of the EEPROM inventory data structure
type Content struct {
	MagicNumber uint16   // 0xFEAD
	Len         uint16   // Length of the payload
	Signature   [64]byte // ED25519 Signature of the payload (only the actual payload length)
	Payload     []byte   // Payload of the EEPROM data (marshaled protobuf)
}

// NewContent creates a new Content instance with the given payload
// and signs it with the provided private key
func NewContent(payload []byte, privateKey ed25519.PrivateKey) (*Content, error) {
	if len(payload) > 0xFFFF {
		return nil, fmt.Errorf("payload too large: %d bytes", len(payload))
	}
	signature := ed25519.Sign(privateKey, payload)
	// Create a new Content instance
	content := &Content{
		MagicNumber: MagicNumber,
		Len:         uint16(len(payload)),
		Signature:   [64]byte(signature),
		Payload:     payload,
	}
	return content, nil
}

// Marshal serializes the Content instance into a byte slice
func (c *Content) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, endianess, c.MagicNumber); err != nil {
		return nil, fmt.Errorf("failed to write magic number: %v", err)
	}
	if err := binary.Write(buf, endianess, c.Len); err != nil {
		return nil, fmt.Errorf("failed to write length: %v", err)
	}

	if err := binary.Write(buf, endianess, c.Signature); err != nil {
		return nil, fmt.Errorf("failed to write signature: %v", err)
	}
	if err := binary.Write(buf, endianess, c.Payload); err != nil {
		return nil, fmt.Errorf("failed to write payload: %v", err)
	}
	return buf.Bytes(), nil
}

// Unmarshal deserializes the byte slice into a Content instance
func Unmarshal(data []byte) (*Content, error) {
	buf := bytes.NewBuffer(data)
	content := &Content{}
	if err := binary.Read(buf, endianess, &content.MagicNumber); err != nil {
		return nil, fmt.Errorf("failed to read magic number: %v", err)
	}
	if err := binary.Read(buf, endianess, &content.Len); err != nil {
		return nil, fmt.Errorf("failed to read length: %v", err)
	}
	if err := binary.Read(buf, endianess, &content.Signature); err != nil {
		return nil, fmt.Errorf("failed to read signature: %v", err)
	}
	content.Payload = make([]byte, content.Len)
	if err := binary.Read(buf, endianess, &content.Payload); err != nil {
		return nil, fmt.Errorf("failed to read payload: %v", err)
	}

	if len(content.Payload) != int(content.Len) {
		return nil, fmt.Errorf("payload length mismatch: expected %d, got %d", content.Len, len(content.Payload))
	}
	return content, nil
}

// VerifySignature verifies the signature of the Content instance using the provided public key
func VerifySignature(content *Content, publicKey ed25519.PublicKey) bool {
	// Verify the signature
	return ed25519.Verify(publicKey, content.Payload, content.Signature[:])
}

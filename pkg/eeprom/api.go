package eeprom

import (
	"crypto/ed25519"
	"fmt"

	"github.com/ci4rail/ee-inv-cli/pkg/eeinvpb"
	"google.golang.org/protobuf/proto"
)

func WriteInventoryToFile(filePathSpec string, inventory *eeinvpb.Inventory, privateKey ed25519.PrivateKey) error {
	// Create a new EEPROM file
	eepromFile, err := NewEEPROMFile(filePathSpec)
	if err != nil {
		return err
	}
	defer eepromFile.Close()

	// Marshal the inventory to protobuf binary format
	payload, err := proto.Marshal(inventory)
	if err != nil {
		return err
	}
	// Create a new Content instance
	content, err := NewContent(payload, privateKey)
	if err != nil {
		return err
	}
	// Marshal the Content instance to a byte slice
	marshaledData, err := content.Marshal()
	if err != nil {
		return err
	}
	// Write the marshaled data to the EEPROM file
	return eepromFile.Write(marshaledData)
}

func ReadInventoryFromFile(filePathSpec string, publicKey ed25519.PublicKey) (*eeinvpb.Inventory, error) {
	// Create a new EEPROM file
	eepromFile, err := NewEEPROMFile(filePathSpec)
	if err != nil {
		return nil, err
	}
	defer eepromFile.Close()

	// Read the data from the EEPROM file
	data, err := eepromFile.Read()
	if err != nil {
		return nil, err
	}

	// Unmarshal the data into a Content instance
	content, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}

	// Verify the signature of the content
	if !ed25519.Verify(publicKey, content.Payload, content.Signature[:]) {
		return nil, fmt.Errorf("signature verification failed")
	}

	// Unmarshal the payload into an Inventory instance
	inventory := &eeinvpb.Inventory{}
	err = proto.Unmarshal(content.Payload, inventory)
	if err != nil {
		return nil, err
	}

	return inventory, nil
}

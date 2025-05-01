package main

import (
	"crypto/ed25519"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"flag"
	"fmt"
	"os"

	"github.com/ci4rail/ee-inv-cli/pkg/eeprom"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

//go:embed keys/ed25519_public.pem
var publicKeyPEM string

func main() {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <eepromfile@offset:size>\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse the command-line flags
	flag.Parse()

	// Get positional arguments
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	eepromPathSpec := args[0]

	inventory, err := eeprom.ReadInventoryFromFile(eepromPathSpec, getPublicKey())
	if err != nil {
		fmt.Println("Error reading inventory from file:", err)
		os.Exit(1)
	}
	// convert to json
	jsonData, err := ProtoToJSONWithDefaults(inventory)
	if err != nil {
		fmt.Println("Error converting inventory to JSON:", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonData))
}

func getPublicKey() ed25519.PublicKey {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return pub.(ed25519.PublicKey)
}

func ProtoToJSONWithDefaults(pb proto.Message) ([]byte, error) {
	marshaler := protojson.MarshalOptions{
		EmitUnpopulated: true, // Show fields with default values
		Indent:          "  ", // Optional: for pretty-printing
	}
	return marshaler.Marshal(pb)
}

# ee-inv-cli
CLI to read inventory data stored in board EEPROMs.

The EEPROM must have been written using the ee-inv-program tool. 
The format of the EEPROM data is defined in `pkg/eeprom/content.go`, the payload is defined in `proto/eeinv.proto`.

The tool verifies the digital signature of the EEPROM data using the public key stored in in the tool. 
It refuses to read the EEPROM if the signature is invalid.

# Usage
```bash
$ ee-inv /sys/bus/i2c/devices/3-0050/eeprom@0:256
{
  "vendor": "Ci4Rail",
  "model": "S101-CPU01",
  "variant": 1,
  "majorVersion": 1,
  "serial": "8507436c-b82e-4b53-ab6f-91e9afb11472"
}
```

# Developer notes
## Generate go from protobuf
```bash
cd proto
protoc -I. eeinv.proto --go_out=../pkg
```
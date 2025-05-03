# ee-inv-cli
CLIs for inventory data stored in board EEPROMs


# Generate go from protobuf
```bash
cd proto
protoc -I. eeinv.proto --go_out=../pkg
```
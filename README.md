# Ethereum Transaction Extractor

This is a tool that extracts transactions from a synchronized Ethereum database using Geth. The extracted transactions are stored in the folder `txs`.

## Requirements
- build-essential
- golang

## Build
```bash
make
```

## Run
```bash
./ethex
```

## Troubleshooting
1. If the build fails, check the Go version. It should align with the version specified in `go.mod` of go-ethereum. For example, you can see the required Go version in go-ethereum v1.13.15 at https://github.com/ethereum/go-ethereum/blob/v1.13.15/go.mod.

2. If extraction fails, ensure that you have synchronized the Ethereum database using Geth, and check the version of go-ethereum in `go.mod`, which should match the version of Geth used for synchronization.

## Notes
You can extract specific block heights by modifying the `start` and `end` parameters in main.go.

## License
This project is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.html).

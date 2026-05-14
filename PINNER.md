# Running DAS-Pinner

## Prerequisites

- A Pinata account with a paid plan that allows CAR file uploads
  (see https://pinata.cloud/pricing). The free tier does not include CAR support.
- A Pinata JWT (Dashboard → API Keys → Create New Key → enable file uploads,
  groups read, optionally testAuthentication).
- Trusted setup for the service.

### Installing the Trusted Setup

To install the trusted setup, run the following command:

```sh
./install-trusted-setup.sh
```

## Configuration

Set the following environment variables before running the service:

- `PINATA_JWT` (required) — your Pinata JWT
- `PINATA_GROUP_ID` (optional) — Pinata group to organize uploads under
- `PINATA_NETWORK` (optional) — `public` (default) or `private`

## Running the service

```sh
export PINATA_JWT=<your-jwt>
./bin/pinner --addr :5080
```

Output:

```sh

░▒▓███████▓▒░ ░▒▓██████▓▒░ ░▒▓███████▓▒░      ░▒▓███████▓▒░░▒▓█▓▒░▒▓███████▓▒░░▒▓███████▓▒░░▒▓████████▓▒░▒▓███████▓▒░  
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░             ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░             ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░░▒▓██████▓▒░       ░▒▓███████▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓██████▓▒░ ░▒▓███████▓▒░  
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓███████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓███████▓▒░       ░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░ 
                                                                                                                       

Version: v0.1.0, commit: 00000000
Initializing root command...
2024-08-19T12:08:39.730-0700	INFO	das-pinner	pinner/main.go:47	Initializing trusted setup...
2024-08-19T12:08:41.533-0700	INFO	das-pinner	ipfs-node/w3storage.go:68	Initialized W3Storage with agent DID: did:key:z6MkfvChtMB5d5WJRGinGBWV1uuVdD6VmefLKPRU8Yog79YS
2024-08-19T12:08:41.915-0700	INFO	das-pinner	ipfs-node/w3storage.go:75	Added space with DID: did:key:z6MkiAxv94CHcwEmFxCRzrkCGq4MJDc1VC8PCCrkgA8wyAHz
generating 2048-bit RSA keypair...done
peer identity: QmY4FqTtiWZykV5D1c4vceYabSNqeh6TsFtqAndMAdcRk6
2024-08-19T12:08:43.046-0700	INFO	das-pinner	api/server.go:58	Starting server on 127.0.0.1:3001...
```

For more options, use the `--help` flag:

```sh
./bin/pinner --help
```

## Running the CLI Tool

To interact with the service, use the CLI tool:

```sh
./bin/pinner-cli upload --data <path-to-data> --addr <pinner-address>
```

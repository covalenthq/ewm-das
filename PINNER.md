# Running DAS-Pinner

## Prerequisites

Before running the service, you need to install the following dependencies:

- [web3.storage](https://web3.storage/docs/go-w3up/#install-w3-cli)
- Trusted setup for the service

### Setting up web3.storage

Create an account on [web3.storage](https://web3.storage/docs/how-to/create-account/#using-the-cli) and generate a private key:

```sh
w3 key create > private.key
```

The contents of the private key should look like this:

```sh
# did:key:z6MkhtbMWQq7dTrZXGuNMWQcFs3Wdr3E4esVbHFMX7GkiHmf
MgCbd3MtiwMFne6Fx7ta22YhWzI+lXEa4KwBQrN1WE/9V9+0BMxBp5XL6JTyn3r1P+IpZTTWBfp+800KqlpkAtCykk1Y=
```

Create permissions to add storage space (store/add) and to upload (upload/add) data:

```sh
w3 delegation create -c 'store/add' -c 'upload/add' -k <did-from-private.key> -o delegation.proof
```

### Installing the Trusted Setup

To install the trusted setup, run the following command:

```sh
./install-trusted-setup.sh
```

## Running the service

To start the service, use the following command:

```sh
./bin/pinner --w3-agent-key <web3.storage-agent-key:MgCbd3M...> --w3-delegation-proof-path delegation.proof
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

# Running DAS-Pinner

## Prerequisites

You need **one** of the following remote pin backends:

- **Self-hosted IPFS (Kubo)**: a running [Kubo](https://github.com/ipfs/kubo)
  daemon reachable over its HTTP RPC (default `http://127.0.0.1:5001/api/v0`).
  Vanilla Kubo accepts unauthenticated requests; for shared deployments
  put it behind a reverse proxy that enforces a bearer token.
- **Filebase**: A Filebase account ([filebase.com](https://filebase.com)). The free tier
  (5 GB storage, 5 GB bandwidth, 1 dedicated gateway, no credit card) is
  sufficient for local testing. An IPFS bucket in your Filebase account and an **IPFS RPC API token**
  scoped to that bucket. Generated from the Filebase console under the
  bucket's settings (separate from the S3 access key — we use the RPC token,
  not the S3 credentials).
- Trusted setup for the service.

### Installing the Trusted Setup

```sh
./install-trusted-setup.sh
```

## Configuration

The pinner picks its backend from environment variables at startup. **If
`IPFS_RPC_URL` is set, the self-hosted backend is used and Filebase is
ignored even if `FILEBASE_RPC_TOKEN` is also present.** The pinner refuses
to start if neither is set.

### Self-hosted IPFS (preferred)

- `IPFS_RPC_URL` (required) — full Kubo RPC base URL, e.g.
  `http://127.0.0.1:5001/api/v0`.
- `IPFS_RPC_TOKEN` (optional) — bearer token, only sent if non-empty.
  Leave unset for vanilla Kubo; set for reverse-proxied / auth-fronted
  Kubo.
- `DEDICATED_GATEWAY` (optional) — a public-ish gateway URL prepended to
  the read-side gateway pool. For a self-hosted setup point it at your
  Kubo gateway (e.g. `http://127.0.0.1:8080/`).

The pinner verifies connectivity at startup by POSTing
`${IPFS_RPC_URL}/version` and refuses to start if the endpoint is
unreachable or returns non-200.

### Filebase

- `FILEBASE_RPC_TOKEN` (required) — your Filebase IPFS RPC API token.
- `DEDICATED_GATEWAY` (optional) — your dedicated Filebase gateway URL
  (e.g. `https://<your-bucket>.myfilebase.com/`).

The pinner verifies the token at startup against
`https://rpc.filebase.io/api/v0/version` and refuses to start if it is
missing or invalid.

## Running the service

Self-hosted (Kubo):

```sh
ipfs daemon &        # or: docker run -d -p 5001:5001 -p 8080:8080 ipfs/kubo:v0.32.1
export IPFS_RPC_URL=http://127.0.0.1:5001/api/v0
export DEDICATED_GATEWAY=http://127.0.0.1:8080/
./bin/pinner --addr :5080
```

Filebase:

```sh
export FILEBASE_RPC_TOKEN=<your-token>
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

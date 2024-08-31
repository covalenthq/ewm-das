# Light Client

## Running Light Client in Docker

### Prerequisites

- Docker

### Usage

Clone the repository:

```bash
git clone https://github.com/covalenthq/das-ipfs-pinner
cd das-ipfs-pinner
```

Build the Docker image:

```bash
docker build -t covalent/light-client -f Dockerfile.lc .
```

Run the Docker container:

```bash
docker run -d --restart always --name light-client -e CLIENT_ID="YOUR UNIQUE ID" covalent/light-client
```

Check the status of the Light Client:

```bash
docker logs -f light-client
```

## Running Light Client Locally

### Prerequisites

- [ipfs](https://docs.ipfs.io/install/command-line/)

### Building from Source

- [Guide](../README.md#building-from-source)

### Running Light Client

To run the light-client, use the following command:

```sh
./bin/light-client --rpc-url <rpc-url> \
    --contract <contract-address> \
    --topic-id <topic-id> \
    --gcp-creds-file <gcp-creds-file> \
    --client-id <client-id> 
```

Note: Client ID is the unique identifier for the client. It can be any string, just make sure it is unique.

```sh
./bin/light-client --rpc-url wss://moonbeam.blastapi.io/618fd77b-a090-457b-b08a-373398006a5e --contract 0x4932bDc983e5146224b9C2e68cfFBFEb004A2824 --topic-id DAS-TO-BQ --gcp-creds-file gcp-creds.json --client-id ${CLIENT_ID}
```

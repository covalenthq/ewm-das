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
docker run -d --restart always --name light-client -e PRIVATE_KEY="YOUR HEX PRIV KEY" covalent/light-client
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
    --collect-url <collect-url> \
    --private-key <private-key> 
```

Note: Private key is the identity of your client. To generate a private key, use can use following resources:

- [Visual-key](https://visualkey.link/)
- [Vanity-eth](https://vanity-eth.tk/)
- [Eth-vanity](https://eth-vanity.io/#calc)

```sh
./bin/light-client --rpc-url wss://coordinator.das.test.covalentnetwork.org/rpc --collect-url https://us-central1-covalent-network-team-sandbox.cloudfunctions.net/ewm-das-collector --private-key ${PRIVATE_KEY}
```

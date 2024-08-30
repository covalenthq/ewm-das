# Running Light Client in Docker

## Prerequisites

- Docker

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/covalenthq/das-ipfs-pinner
    cd das-ipfs-pinner
    ```

2. Build the Docker image:

    ```bash
    docker build -t covalent/light-client -f Dockerfile.lc .
    ```

3. Run the Docker container:

    ```bash
    docker run -d --name light-client -e CLIENT_ID="YOUR UNIQUE ID" covalent/light-client
    ```

4. Check the status of the Light Client:

    ```bash
    docker logs -f light-client
    ```

# Light Client installation guide for Mac OS X

## Prerequisites

- [Homebrew](https://brew.sh/)

## IPFS Kubo client installation

Run the following command to install the IPFS Kubo client:

```bash
brew install ipfs
```

## IPFS Kubo client configuration

Run the following command to initialize the IPFS Kubo client:

```bash
ipfs init
```

## Install & Launch Light Client

Go to the [Light Client Releases](https://github.com/covalenthq/das-ipfs-pinner/releases) page and download the latest release.

Unzip the downloaded archive and navigate to the `bin` directory:

```bash
cd bin
```

Copy GCP credentials to the `bin` directory:

```bash
cp /path/to/your/credentials.json ./gcp-credentials.json
```

**Note:** The `gcp-credentials.json` file is required to authenticate with the Google Cloud Platform. You can obtain the credentials by downloading from [Covalent Slack Channel](https://covalent-hq.slack.com/archives/C071MF2RG76/p1724728668807929).

Open the `install.sh` file and update the `{YOUR_UNIQUE_ID}` with your unique ID at line 89:

```bash
...
    --client-id myuniqueid@covalenthq.com
...
```

Run the following command to install the Light Client:

```bash
./install.sh
```

The script will install all files in `$HOME/.covalenthq` directory.

To uninstall the Light Client, run the following command:

```bash
./uninstall.sh
```

## Status

To check the status of the Light Client, run the following command:

```bash
tail -n 1000 -f $HOME/.covalenthq/light-client.log
```

Result:

```bash

░▒▓█▓▒░      ░▒▓█▓▒░░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░       ░▒▓██████▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓████████▓▒░▒▓███████▓▒░▒▓████████▓▒░ 
░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░          ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░     
░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░          ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░     
░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒▒▓███▓▒░▒▓████████▓▒░  ░▒▓█▓▒░          ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░▒▓██████▓▒░ ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░     
░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░          ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░     
░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░          ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░     
░▒▓████████▓▒░▒▓█▓▒░░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░           ░▒▓██████▓▒░░▒▓████████▓▒░▒▓█▓▒░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░     
                                                                                                                                         
                                                                                                                                         

Version: e8c74c7, commit: e8c74c79e77cf5c65ada0cf9a3c74390022d11ae
2024-08-26T20:25:04.845-0700	INFO	light-client	light-client/main.go:96	Starting client...
2024-08-26T20:25:05.600-0700	INFO	light-client	event-listener/listener.go:82	Subscribed to logs for contract: 0x916B54696A70588a716F899bE1e8f2A5fFd5f135
...
```

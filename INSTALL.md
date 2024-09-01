# Light Client installation guide for Mac OS X

## Prerequisites

- [Homebrew](https://brew.sh/)

## IPFS Kubo client installation

NOTE: If you have already installed the client, you can skip this step.

Run the following command to install the client:

```bash
brew install ipfs
```

## IPFS Kubo client configuration

NOTE: If you have already configured the client, you can skip this step.

Run the following command to initialize the client:

```bash
ipfs init
```

## Install & Launch Light Client

Go to the [Light Client Releases](https://github.com/covalenthq/das-ipfs-pinner/releases) page and download the latest release.

Unzip the downloaded archive and navigate to the directory:

```bash
cd das-macos-latest
```

Run the following command to install the Light Client:

```bash
./install.sh <YOUR_UNIQUE_ID>
```

The script will install all files in `$HOME/.covalent` directory.

To uninstall the Light Client, run the following command:

```bash
$HOME/.covalent/uninstall.sh
```

## Status

To check the status of the Light Client, run the following command:

```bash
tail -n 1000 -f $HOME/.covalent/light-client.log
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
                                                                                                                                         
                                                                                                                                         

Version: v0.3.0, commit: e8c74c79e77cf5c65ada0cf9a3c74390022d11ae
2024-08-26T20:25:04.845-0700	INFO	light-client	light-client/main.go:96	Starting client...
2024-08-26T20:25:05.600-0700	INFO	light-client	event-listener/listener.go:82	Subscribed to logs for contract: 0x916B54696A70588a716F899bE1e8f2A5fFd5f135
...
```

# Light Client installation guide for Mac OS X

> **RETIRED (2026-11-23).** The single-command macOS installer below is no longer published: it downloaded the installer and the light client from a storage location that no longer exists, so it does not work. Run the light client with Docker or build it from source instead (see the light client documentation). Never pipe a downloaded script into `bash` together with your private key.

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

<!-- RETIRED 2026-11-23 (installer no longer published; see the note at the top)
## Install & Launch Light Client

Single command to install and launch the Light Client:

```bash
curl -sSf https://storage.googleapis.com/ewm-release-artefacts/{{VERSION}}/macos/install.sh | bash -s <YOUR PRIVATE KEY>
```

To generate a private key, use can use following resources:

- [Visual-key](https://visualkey.link/)
- [Vanity-eth](https://vanity-eth.tk/)
- [Eth-vanity](https://eth-vanity.io/#calc)

The script will install all files in `$HOME/.covalent` directory.

To uninstall the Light Client, run the following command:

```bash
$HOME/.covalent/uninstall.sh
```
-->

## Status

To check the status of the Light Client, run the following command:

```bash
tail -n 1000 -f $HOME/.covalent/light-client.log
```

Result:

```bash

███████ ██     ██ ███    ███      ██████ ██      ██ ███████ ███    ██ ████████ 
██      ██     ██ ████  ████     ██      ██      ██ ██      ████   ██    ██    
█████   ██  █  ██ ██ ████ ██     ██      ██      ██ █████   ██ ██  ██    ██    
██      ██ ███ ██ ██  ██  ██     ██      ██      ██ ██      ██  ██ ██    ██    
███████  ███ ███  ██      ██      ██████ ███████ ██ ███████ ██   ████    ██    
                                                                               

Version: v0.9.0, commit: 8d6709bee9e79d3c4e6ece35fed65da02f3850f4
2024-09-18T15:45:01.238-0700	INFO	light-client	light-client/main.go:91	Starting client...
2024-09-18T15:45:01.238-0700	INFO	light-client	light-client/main.go:97	Client idenity: 0x51b6D674514849aF97FB77BCac51bcdD7799842C
...
```

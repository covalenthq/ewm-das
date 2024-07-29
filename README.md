# Pinner

Pinner is a Go-based project consisting of a daemon service and a CLI tool. The daemon service handles backend operations, while the CLI tool provides an interface for interacting with the daemon.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Building from Source](#building-from-source)
- [Configuration](#configuration)
- [Development](#development)
- [License](#license)
- [Contributing](#contributing)

## Installation

To install Pinner, download the pre-built binaries from the [Releases](#) page or build from source as described below.

## Usage

### Running the Daemon

To start the daemon, use the following command:

```sh
./bin/pinner
./bin/pinner -debug
```

### Running the CLI Tool

To interact with the daemon, use the CLI tool:

```sh
./bin/pinner-cli -mode=store -data="Your data here"
./bin/pinner-cli -mode=extract
```

## Building from Source

### Prerequisites

- Go 1.22 or later

### Build Commands

1. Clone the repository:

```sh
git clone https://github.com/covalenthq/das-ipfs-pinner
cd das-ipfs-pinner
```

2. Install dependencies:

```sh
make deps
```

3. Build binaries:

```sh
make
```

This will compile the daemon and CLI tool into the bin directory.

### Clean Up

To clean up the build artifacts, run:

```sh
make clean
```

## Configuration

The project uses environment variables and flags for configuration. For example, you can set `DAEMON_ADDR` to change the address the daemon listens on.

## Development

### Formatting and Linting

To format the code, run:

```sh
make fmt
```

To lint the code, run:

```sh
make lint
```

### Testing

To run tests, use the following command:

```sh
make test
```

To run tests with a coverage report:

```sh
make test-cover
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.


## Contributing

TODO: Add contribution guidelines
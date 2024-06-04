# Qase CLI

`qli` is Qase on the command line. It allows you to work with test runs and import test results.

Qase CLI is available for Qase and Qase Enterprise users.

# Installation

## Build from source

1. Clone the repository

```bash
git clone https://github.com/qase-tms/qasectl.git && cd qasectl
```

2. Build the binary

```bash
make build
```

You will find the binary in the `build` directory.

## Docker image

1. Pull the Docker image

```bash
docker pull ghcr.io/qase-tms/qase:latest
```

2. Run the Docker container

```bash
docker run -it ghcr.io/qase-tms/qase:latest 
```

# Usage

The tool is designed to be used in a terminal.
You can run `qli` with the `--help` flag to see the available commands and options.

You can find more information about the commands and options in the [documentation](docs/command.md).

# Qase CLI

`qasectl` is the command-line interface (CLI) for Qase, enabling users to interact with test runs and import test results
directly from the terminal.

Qase CLI is available for both **Qase** and **Qase Enterprise** users.

## Installation

### Install via Homebrew (macOS/Linux)

```bash
brew tap qase-tms/tap
brew install qasectl
```

### Install via `go install`

The easiest way to install Qase CLI is using `go install`:

```bash
go install github.com/qase-tms/qasectl@latest
```

Make sure to add `$GOPATH/bin` to your `$PATH` environment variable to be able to run the `qasectl` command.

### Build from Source

If you'd like to build the CLI from source, follow these steps:

1. Clone the repository:

    ```bash
    git clone https://github.com/qase-tms/qasectl.git && cd qasectl
    ```

2. Build the binary using `make`:

    ```bash
    make build
    ```

   The compiled binary will be located in the `build` directory.

### Docker Image

You can also use Qase CLI via Docker. Follow these steps:

1. Pull the latest Docker image:

    ```bash
    docker pull ghcr.io/qase-tms/qase-cli:latest
    ```

2. Run the Docker container:

    ```bash
    docker run --rm ghcr.io/qase-tms/qase-cli:latest version
    ```

## Usage

`qasectl` is designed to be used directly in your terminal. You can run the following command to view available options:

```bash
qasectl --help
```

This will show all available commands and their descriptions.

### Qase Enterprise

By default the CLI talks to `api.qase.io`. Qase Enterprise users can point it at their own instance with the
`--api-host` option, or with the `QASE_TESTOPS_API_HOST` environment variable:

```bash
qasectl testops run create --api-host api.qase.example.com --project PRJ --token <token> --title "Test Run 1"

# or
export QASE_TESTOPS_API_HOST=api.qase.example.com
```

For more detailed information about each command and option, refer to the [full documentation](docs/command.md).

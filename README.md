# littlesocks

A SOCKS5 proxy server written in Go, supporting multiple platforms and architectures.

## Features

- SOCKS5 proxy implementation
- Configurable network address binding
- Cross-platform support (Linux, macOS, Windows)
- Multi-architecture (AMD64, ARM64)
- No authentication required
- Static binary with no runtime dependencies

## Supported Platforms

| OS | Architecture |
|---|---|
| Linux | AMD64, ARM64 |
| macOS | AMD64, ARM64 |
| Windows 11 | AMD64, ARM64 |

## Installation

### Homebrew (macOS and Linux)

Install via Homebrew:

```bash
brew install cloudygreybeard/tap/littlesocks
```

Or tap the repository first:

```bash
brew tap cloudygreybeard/tap
brew install littlesocks
```

### Container Images

Pull and run the latest image from GitHub Container Registry:

```bash
podman pull ghcr.io/cloudygreybeard/littlesocks:latest
podman run -p 1080:1080 ghcr.io/cloudygreybeard/littlesocks:latest
```

Images are available for `linux/amd64` and `linux/arm64`.

Note: `docker` commands work identically if you prefer Docker over Podman.

### Download Pre-built Binaries

Download the latest release for your platform from the [Releases](https://github.com/cloudygreybeard/littlesocks/releases) page.

### Build from Source

Requirements:
- Go 1.21 or higher

```bash
# Clone the repository
git clone https://github.com/cloudygreybeard/littlesocks.git
cd littlesocks

# Build
go build -o littlesocks

# Or install directly
go install github.com/cloudygreybeard/littlesocks@latest
```

## Usage

### Container Usage

Run with Podman:

```bash
# Default configuration (listens on 0.0.0.0:1080 inside container)
podman run -p 1080:1080 ghcr.io/cloudygreybeard/littlesocks:latest

# Custom port mapping
podman run -p 9050:1080 ghcr.io/cloudygreybeard/littlesocks:latest

# Custom address inside container
podman run -p 1080:1080 ghcr.io/cloudygreybeard/littlesocks:latest -addr 0.0.0.0:1080

# Run in background
podman run -d -p 1080:1080 --name littlesocks ghcr.io/cloudygreybeard/littlesocks:latest
```

### Binary Usage

Run the SOCKS5 proxy server on the default address (127.0.0.1:1080):

```bash
./littlesocks
```

### Custom Address

Specify a custom address and port:

```bash
# Listen on all interfaces, port 1080
./littlesocks -addr :1080

# Listen on specific IP and port
./littlesocks -addr 192.168.1.100:8080

# Listen only on localhost, custom port
./littlesocks -addr 127.0.0.1:9050
```

### Command Line Options

```
-addr string
    Address to bind the SOCKS5 server (host:port) (default "127.0.0.1:1080")
-version
    Print version information
```

### Examples

#### Local Development

```bash
./littlesocks -addr 127.0.0.1:1080
```

#### Network-wide Access

```bash
./littlesocks -addr 0.0.0.0:1080
```

⚠️ **Security Warning**: Binding to `0.0.0.0` makes the proxy accessible from any network interface. Use appropriate firewall rules in production environments.

## Testing the Proxy

### Using curl

```bash
# Test with curl
curl -x socks5h://127.0.0.1:1080 https://api.ipify.org
```

### Using Firefox

1. Open Firefox Settings
2. Navigate to Network Settings
3. Select "Manual proxy configuration"
4. Set SOCKS Host: `127.0.0.1`, Port: `1080`
5. Select "SOCKS v5"
6. Click OK

### Using ssh

```bash
ssh -o ProxyCommand="nc -X 5 -x 127.0.0.1:1080 %h %p" user@remote-host
```

## Development

### Prerequisites

- Go 1.21+
- [goreleaser](https://goreleaser.com/) (for building releases)

### Local Development

```bash
# Install dependencies
go mod download

# Run the application
go run main.go -addr :1080

# Run tests
go test ./...

# Build for current platform
go build -o littlesocks
```

### Building for Multiple Platforms

Using goreleaser for local snapshot builds:

```bash
# Build for all platforms
goreleaser build --snapshot --clean

# Binaries will be in ./dist/
```

### Cross-compilation with Go

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o littlesocks-linux-amd64

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o littlesocks-linux-arm64

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o littlesocks-windows-amd64.exe

# Windows ARM64
GOOS=windows GOARCH=arm64 go build -o littlesocks-windows-arm64.exe

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -o littlesocks-darwin-amd64

# macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o littlesocks-darwin-arm64
```

## Releasing

This project uses GitHub Actions with goreleaser for automated releases.

### Creating a Release

1. Tag your commit with a semantic version:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. GitHub Actions will automatically:
   - Build binaries for all platforms
   - Create archives (tar.gz for Unix, zip for Windows)
   - Build and push container images for linux/amd64 and linux/arm64
   - Update Homebrew formula in cloudygreybeard/homebrew-tap
   - Generate checksums
   - Create a GitHub release with all artifacts

### Release Workflow

The release workflow (`.github/workflows/release.yml`) triggers on version tags (`v*`) and:
- Uses Ubuntu latest runner
- Sets up Go 1.21
- Sets up Docker Buildx for multi-arch builds
- Runs goreleaser to build and publish binaries and container images

## Architecture

The application is built with:
- **[armon/go-socks5](https://github.com/armon/go-socks5)**: SOCKS5 server implementation
- **Standard Go libraries**: For networking and signal handling
- **No runtime dependencies**: Static binary compilation

Container images use `FROM scratch` for minimal size (binary only, no OS layer).

## Configuration

Currently, the proxy operates without authentication. For production use with authentication, you can extend the `socks5.Config` in `main.go`:

```go
conf := &socks5.Config{
    AuthMethods: []socks5.Authenticator{
        socks5.UserPassAuthenticator{
            Credentials: socks5.StaticCredentials{
                "user": "password",
            },
        },
    },
}
```

## Graceful Shutdown

The server handles `SIGINT` and `SIGTERM` signals for graceful shutdown:

```bash
# Send interrupt signal (Ctrl+C)
^C

# Or send SIGTERM
kill -TERM <pid>
```

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome. Please submit a Pull Request.

## Troubleshooting

### Port Already in Use

If you get an error about the port being in use:

```bash
# Check what's using the port (macOS/Linux)
lsof -i :1080

# Kill the process if needed
kill -9 <pid>
```

### Permission Denied (Ports < 1024)

On Unix systems, ports below 1024 require root privileges:

```bash
# Use sudo for privileged ports
sudo ./littlesocks -addr :80
```

### Firewall Issues

Ensure your firewall allows the port:

```bash
# macOS
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --add littlesocks

# Linux (ufw)
sudo ufw allow 1080/tcp

# Windows
# Use Windows Defender Firewall settings
```

## Acknowledgments

- [armon/go-socks5](https://github.com/armon/go-socks5) - SOCKS5 server library
- [goreleaser](https://goreleaser.com/) - Release automation

## Support

For bugs, questions, or feature requests, please open an issue on [GitHub](https://github.com/cloudygreybeard/littlesocks/issues).


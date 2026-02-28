# Onboarding Guide: gns3util

Welcome to the `gns3util` project! This tool is designed to automate the GNS3 v3 API, with a primary focus on educational environments.

## What is gns3util?

`gns3util` is a CLI tool that helps educators manage GNS3 labs for classes. It automates tasks like:
- Creating classes and student groups.
- Deploying exercises from templates (server-based or file-based).
- Managing remote GNS3 servers (installation, HTTPS setup).
- Orchestrating clusters of GNS3 servers.

## Codebase Map

Here is a quick overview of where things are:

- `cmd/gns3util/`: The entry point of the application.
- `internal/cli/`: Contains all the CLI command definitions using [Cobra](https://github.com/spf13/cobra).
    - `root.go`: The root command and global flag definitions.
    - `cmds/`: Subcommands organized by category (e.g., `class`, `exercise`, `auth`).
    - `cli_pkg/`: Shared logic and packages used specifically by the CLI (config, utils, fuzzy finding).
- `pkg/api/`: The GNS3 v3 API client.
    - `client.go`: Low-level HTTP client and request handling.
    - `endpoints/`: High-level wrappers for API endpoints (GET, POST, etc.).
    - `schemas/`: Go structs representing API request/response bodies.
- `pkg/ssh/`: SSH utilities for remote server management.
- `scripts/`: Example scripts and templates for testing and common workflows.

## Tech Stack

- **Language**: Go
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **Shell Completions**: [Carapace](https://github.com/rsteube/carapace)
- **Database**: SQLite (managed via [sqlc](https://sqlc.dev/) for cluster management)
- **API**: GNS3 v3 REST API

## Getting Started for Developers

1. **Clone the repo**
2. **Build the project**:
   ```bash
   go build -o gns3util ./cmd/gns3util
   ```
3. **Run tests**:
   ```bash
   go test ./...
   ```

## Contributing Standards

Before submitting a Pull Request, you **must** ensure the code follows our quality standards. No lint errors are allowed.

Run the following commands using [mise](https://mise.jdx.dev/):
```bash
# Run the linter
mise run lint

# Format the code
mise run format
```

All contributions must pass linting and formatting checks to be considered for merge.

## Next Steps

- Check out [ARCHITECTURE.md](./ARCHITECTURE.md) for a deeper dive into how the project is structured.
- See [PACKAGES.md](./PACKAGES.md) for a detailed breakdown of every package in the codebase.

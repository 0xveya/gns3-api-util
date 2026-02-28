# gns3util Package Documentation

This document provides a detailed breakdown of each package in the `gns3util` codebase, explaining its intent, responsibilities, and how it relates to other parts of the system.

## Root Directories

### `cmd/gns3util/`
- **Purpose**: Application entry point.
- **Contents**: Contains `main.go`, which simply calls `cli.Execute()`.
- **Relates to**: `internal/cli`.

---

## CLI Layer (`internal/cli/`)

This is where the user interface (CLI) is defined. We use [Cobra](https://github.com/spf13/cobra) for command routing.

### `internal/cli/`
- **Purpose**: Root CLI command definitions.
- **Key Files**: 
    - `root.go`: Sets up the base command, global flags (`--server`, `--key-file`), and initializes the global configuration context.
- **Logic**: Most files here (like `project.go`, `node.go`, etc.) define the "command groups" (e.g., `gns3util project ...`).

### `internal/cli/cmds/`
- **Purpose**: Concrete implementations of subcommands.
- **Structure**: Organized by resource type (e.g., `auth`, `class`, `exercise`, `project`).
- **Logic**: Each file typically contains a `New...Cmd()` function and a `run...()` function that handles the command's execution logic, flag parsing, and calls to the `pkg/api` layer.

### `internal/cli/cli_pkg/`
Internal support packages for the CLI.

#### `cli_pkg/authentication/`
- **Purpose**: Handles authentication logic, token storage, and keyfile management.
- **Relates to**: Used by `cmds/auth` and the API client.

#### `cli_pkg/config/`
- **Purpose**: Manages the CLI's internal state and configuration.
- **Logic**: Uses Go context to pass global options (server URL, insecure flag, etc.) down through command execution.

#### `cli_pkg/cluster/`
- **Purpose**: Core logic for GNS3 cluster orchestration.
- **Subpackages**:
    - `db/`: Handles the SQLite database for storing cluster metadata, using `sqlc` for type-safe queries.
- **Logic**: Manages how multiple GNS3 nodes are grouped together and how projects are distributed among them.

#### `cli_pkg/fuzzy/`
- **Purpose**: Interactive UI for the CLI.
- **Logic**: Implements a fuzzy finder (using `go-fuzzyfinder`) to allow users to interactively select projects, templates, or nodes.

#### `cli_pkg/utils/`
A collection of specialized utilities used across the CLI.
- `utils/class/`: Logic for creating and managing classes and student groups.
- `utils/messageUtils/`: The standard way to print output. It handles colors, bold text, and status prefixes (`Success:`, `Error:`).
- `utils/colorUtils/`: Low-level ANSI color handling.
- `utils/server/`: An interactive web server used specifically for the "interactive class creator" UI.
- `utils/pathUtils/`: Cross-platform path handling (e.g., finding the config directory).

---

## API Layer (`pkg/api/`)

A standalone package for interacting with the GNS3 v3 REST API. This could theoretically be used as a library by other projects.

### `pkg/api/`
- **`client.go`**: The core `GNS3ApiClient`. It handles:
    - Base URL and API versioning (`/v3`).
    - Authentication headers.
    - Timeout and retry logic.
    - Decoding JSON responses and handling common error codes (403, 422).

### `pkg/api/endpoints/`
- **Purpose**: A centralized "map" of all GNS3 API routes.
- **Logic**: Instead of hardcoding URLs in the CLI code, we call functions like `endpoints.Get.Project(id)`. This makes it easy to update the tool if the GNS3 API changes.
- **Categories**: Divided into `get.go`, `post.go`, `put.go`, and `delete.go`.

### `pkg/api/schemas/`
- **Purpose**: Data models for the GNS3 API.
- **Contents**: Structs that match the JSON request and response bodies of the GNS3 API. Includes `Project`, `Node`, `User`, `ACL`, etc.

---

## Remote Management (`pkg/ssh/`)

- **Purpose**: Low-level infrastructure automation.
- **Logic**: Uses SSH to connect to remote servers to:
    - Install/Uninstall the GNS3 server.
    - Set up Caddy for SSL/HTTPS.
    - Manage systemd services.
    - This package is decoupled from the GNS3 API as it operates at the OS level.

---

## Summary of Relationships

1.  **User** runs a command (e.g., `gns3util project ls`).
2.  **`internal/cli/root.go`** parses global flags.
3.  **`internal/cli/cmds/project/`** (logic) is called.
4.  **`pkg/api/`** (client) is used to make a request.
5.  **`pkg/api/endpoints/`** provides the URL.
6.  **`pkg/api/schemas/`** provides the data structure for the response.
7.  **`internal/cli/cli_pkg/utils/messageUtils/`** formats the final output for the user.

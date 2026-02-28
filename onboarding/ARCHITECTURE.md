# gns3util Architecture

`gns3util` is a CLI tool built with [Cobra](https://github.com/spf13/cobra). It follows a standard Go project layout.

## Project Structure Deep Dive

### CLI Layer (`internal/cli`)

The CLI layer is where the commands are defined. Each command group is typically located in `internal/cli/cmds/`.

- `root.go`: Entry point for Cobra. Global flags (`--server`, `--key-file`, etc.) are defined here.
- `PersistentPreRunE`: Used in the root command to validate global flags and initialize context with `GlobalOptions`.

### Shared CLI Logic (`internal/cli/cli_pkg`)

- `config`: Handles global CLI configuration and context.
- `utils`: Contains helper functions for common tasks:
    - `class/`: Class and group creation logic.
    - `server/`: Interactive server for class creation.
    - `messageUtils/`: Pretty-printing messages with colors and icons.
- `fuzzy/`: Logic for the fuzzy finder UI.

### API Layer (`pkg/api`)

The API layer is responsible for all communication with the GNS3 v3 server.

- `GNS3ApiClient`: Low-level HTTP client handling headers, authentication, and request execution.
- `endpoints/`: High-level wrappers for GNS3 API endpoints.
- `schemas/`: Go structs mapping to GNS3 API JSON objects.

### Remote Server Management (`pkg/ssh`)

Utilities for managing remote servers via SSH, including:
- Installing/uninstalling GNS3.
- Configuring Caddy for HTTPS.
- Firewall management.

## Development Workflow

### Adding a New Command

1.  **Define the command**: Create a new file in `internal/cli/cmds/` or a new subcommand in an existing file.
2.  **Add flags**: Use Cobra to define flags for your command.
3.  **Implement logic**: Add the core logic in a separate package or within the command's `RunE` function.
4.  **Register the command**: Add your new command to the root command in `internal/cli/root.go`.

### Working with the GNS3 API

When adding support for a new GNS3 API endpoint:

1.  **Add the schema**: Define the request/response structs in `pkg/api/schemas/`.
2.  **Add the endpoint**: Implement the high-level wrapper in `pkg/api/endpoints/`.
3.  **Use the client**: Call the endpoint from your CLI command logic.

### Cluster Management and Database

For cluster-related features, the project uses SQLite and [sqlc](https://sqlc.dev/).

- Database migrations are handled in `internal/cli/cli_pkg/cluster/db/`.
- `sqlc` generates Go code from SQL queries defined in `.sql` files.

## Tools and Utilities

- `messageUtils`: Use this for all CLI output to ensure consistency.
    - `SuccessMsg`, `ErrorMsg`, `WarningMsg`, `InfoMsg`, `Bold`.
- `fuzzy`: Use this to provide interactive selection for resources like projects, nodes, and templates.

## Error Handling

- Return meaningful errors from lower-level packages.
- Use `messageUtils.ErrorMsg` in the `Execute` function to display errors to the user.
- Prefer `fmt.Errorf("...: %w", err)` to wrap errors and maintain context.

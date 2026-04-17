# NPC Descriptor

A simple tool and API to generate random descriptors for NPCs (Non-Player Characters) in both English and French.

## Features

- **Multi-language support**: English (en) and French (fr).
- **CLI Tool**: Quick generation of descriptors from the command line.
- **REST API**: JSON endpoint for integration into other applications.
- **Embedded Data**: Descriptors are embedded into the Go binary for easy distribution.

## Installation

Ensure you have Go installed, then clone the repository and build:

```bash
make build
```

## CLI Usage

Run the CLI tool to get three random descriptors for a character.

```bash
# Default (English)
go run cmd/cli/main.go

# French
go run cmd/cli/main.go --lang fr
```

### Output Example (en)
`This character can be described as Suspicious, Wary and Bitter`

### Output Example (fr)
`Ce personnage peut être décrit comme étant Dépendant, Sévère et Cupide`

## API Usage

Start the API server:

```bash
make run
```

The server runs on port `8080` by default.

### Endpoint: `GET /api/descriptors`

Returns a list of random descriptors.

**Query Parameters:**
- `n`: Number of descriptors to return (1-10, default: 3).
- `lang`: Locale to use (`en` or `fr`, default: `en`).

**Example Request:**

```bash
curl "http://localhost:8080/api/descriptors?n=2&lang=fr"
```

**Example Response:**

```json
{
  "descriptors": ["Stoïque", "Attrayant"]
}
```

## Project Structure

- `data/`: JSON files containing descriptors for each locale.
- `cmd/`: Entry points for the CLI and API applications.
- `internal/`: Internal logic for the API.
- `descriptor.go`: Core library logic for selecting descriptors.

## License

MIT

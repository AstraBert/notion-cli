# notion-cli

CLI to quickly create and read Notion pages from the terminal

## Installation

In order to install `notion-cli` there are three ways:

1. Using `go`: if you already have `go` 1.22+ installed in your environment, installing `notion-cli` is effortless

```bash
go install github.com/AstraBert/notion-cli@latest
```

2. Using `npm` (recommended):

```bash
npm install @cle-does-things/notion-cli@latest
```

3. Install with `brew` (**coming soon!**)

In this last case, be careful to specify your OS (supported: linux, windows, macos) and your architecture type (supported: amd, arm).

## Usage

`notion-cli` has two commands: `read` (aliased also to `r`) and `write` (aliased also to `w`).

In order to use the commands, you need to have `NOTION_API_KEY` available within your environment:

```bash
export NOTION_API_KEY="..."
```

Follow [this guide](https://developers.notion.com/guides/get-started/create-a-notion-integration) to get your API key and provide it with the necessary permissions to read and write pages.

### `read`

**Aliases:** `r`

Fetches the content of a Notion page by its ID and prints it to stdout.

**Usage**

```bash
notion-cli read <page-id>
```

**Arguments**

| Argument | Description | Required |
|----------|-------------|----------|
| `page-id` | The ID of the Notion page to read | Yes |

**Examples**

```bash
# Read a page by its ID
notion-cli read abb4215a-8f8f-47fb-81e5-353a0aec683f

# Save the output to a file
notion-cli read abb4215a-8f8f-47fb-81e5-353a0aec683f > page.md
```

### `write`

**Aliases:** `w`

Creates a new Notion page under a given parent (either a page or a database) and prints the ID of the newly created page to stdout.

### Usage

```bash
notion-cli write --parent-id <id> --content <markdown> [--parent-type <type>] [--title <title>]
```

**Flags**

| Flag | Alias | Description | Required | Default |
|------|-------|-------------|----------|---------|
| `--parent-id` | `-i` | ID of the parent page or database | Yes | — |
| `--content` | `-c` | Markdown content for the page body | Yes | — |
| `--parent-type` | `-p` | Type of the parent: `page` or `database` | No | `page` |
| `--title` | `-t` | Title for the new page | No | `""` |

**Examples**

```bash
# Create a page under another page
notion-cli write --parent-id abb4215a-8f8f-47fb-81e5-353a0aec683f \
  --content "# Hello\nThis is my new page." \
  --title "My New Page"

# Create a page under a database
notion-cli write --parent-id abb4215a-8f8f-47fb-81e5-353a0aec683f \
  --parent-type database \
  --content "# Meeting Notes\nDiscussed Q1 roadmap." \
  --title "Meeting Notes"

# Using short flags
notion-cli write -i abb4215a-8f8f-47fb-81e5-353a0aec683f -c "Hello world" -t "My Page"
```

## Contributing

We welcome contributions! Please read our [Contributing Guide](./CONTRIBUTING.md) to get started.

## License

This project is licensed under the [MIT License](./LICENSE)

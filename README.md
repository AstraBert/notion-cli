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

3. Install with `brew` (only for Linux and Mac users):

```bash
brew tap AstraBert/notion-cli
brew install notion-cli
```

## Usage

`notion-cli` has three commands: `read` (aliased also to `r`), `write` (aliased also to `w`) and `append` (aliased also to `a`).

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

**Flags**

| Flag | Alias | Description | Required | Default |
|------|-------|-------------|----------|---------|
| `--max-retries` | `-m` | Maximum number of retries for failed API calls | No | `3` |
| `--retry-interval` | `-r` | Retry interval (in seconds) for failed API calls | No | `1` |

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
| `--max-retries` | `-m` | Maximum number of retries for failed API calls | No | `3` |
| `--retry-interval` | `-r` | Retry interval (in seconds) for failed API calls | No | `1` |

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

### `append` 

**Aliases:** `w`

Appends markdown content to the end of an existing Notion page and prints the ID of the modified page to stdout.

**Arguments**

| Argument | Description | Required |
|----------|-------------|----------|
| `page-id` | The ID of the Notion page to modify | Yes |

**Flags**

| Flag | Alias | Description | Required | Default |
|------|-------|-------------|----------|---------|
| `--content` | `-c` | Markdown content for the page body | Yes | — |
| `--max-retries` | `-m` | Maximum number of retries for failed API calls | No | `3` |
| `--retry-interval` | `-r` | Retry interval (in seconds) for failed API calls | No | `1` |

**Examples**

```bash
# Append content to a page
notion-cli append abb4215a-8f8f-47fb-81e5-353a0aec683f \
  --content "# Hello\nThis is a new block."

# Using short flags
notion-cli append abb4215a-8f8f-47fb-81e5-353a0aec683f -c "Hello world"
```

## Use as an Agent Skill

You can use `notion-cli` as an agent skill, downloading it with the `skills` CLI tool:

```bash
npx skills add AstraBert/notion-cli
```

Or copy-pasting the [`SKILL.md`](./skills/notion-cli/SKILL.md) file to your own skills setup.

## Contributing

We welcome contributions! Please read our [Contributing Guide](./CONTRIBUTING.md) to get started.

## License

This project is licensed under the [MIT License](./LICENSE)

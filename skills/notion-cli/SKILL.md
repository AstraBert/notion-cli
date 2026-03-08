---
name: notion-cli
description: Use this skill when it is necessary to quickly read, write or modify a Notion page.
compatibility: Requires brew, npm or go 1.22+ to be installed. Requires a NOTION_API_KEY available within the environment to be executed.
license: MIT
metadata:
  author: Clelia Astra Bertelli
  version: "0.1.0"
---

# `notion-cli` Skill

Quickly read, write and modify Notion pages from the terminal.

## Initial Setup

When this skill is invoked, respond with:

```
I'm ready to use notion-cli to read, write or modify Notion pages. Before we begin, please confirm that:

- `notion-cli` is installed globally (`go install AstraBert/notion-cli@latest` or `npm install -g @cle-does-things/notion-cli@latest` or `brew tap AstraBert/notion-cli && brew install notion-cli`)
- A `NOTION_API_KEY` is available as an environment variable in the current shell

If both are set, please provide:

1. The ID of a page to read
2. The content (or the prompt to produce the content) of a page to write, and the ID of the parent element (specifying whether the parent element is a database or a page)
3. The ID of a page to modify, with the content to append (or a prompt to generate the content)

I will produce the appropriate `notion-cli` command, and once execution is approved, report the results.
```

Then wait for the user's input.

---

## Step 0: Install notion-cli (if needed)

If `notion-cli` is not yet installed, install it globally:

- With `npm` (recommended)

```bash
npm i -g @cle-does-things/notion-cli
```

- With `go` (1.22+ supported):

```bash
go install AstraBert/notion-cli@latest
```

- With `brew` (recommended for MacOs and Linux users):

```bash
brew tap AstraBert/notion-cli
brew install notion-cli
```

Verify installation:

```bash
notion-cli --help
```

---

## Step 1: Produce the CLI Command

### `read`: Fetch a Page
```bash
# Read a page by its ID
notion-cli read abb4215a-8f8f-47fb-81e5-353a0aec683f
# Save output to a file
notion-cli read abb4215a-8f8f-47fb-81e5-353a0aec683f > page.md
# Using the short alias
notion-cli r abb4215a-8f8f-47fb-81e5-353a0aec683f
```

### `write`: Create a New Page

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

### `append`: Append to an Existing Page

```bash
# Append content to a page
notion-cli append abb4215a-8f8f-47fb-81e5-353a0aec683f \
  --content "# Hello\nThis is a new block."
# Using short flags
notion-cli append abb4215a-8f8f-47fb-81e5-353a0aec683f -c "Hello world"
```

### Key Options Reference

**`read`**

| Flag | Alias | Description | Default |
|------|-------|-------------|---------|
| `--max-retries` | `-m` | Max retries for failed API calls | `3` |
| `--retry-interval` | `-r` | Retry interval in seconds | `1` |

**`write`**

| Flag | Alias | Description | Required | Default |
|------|-------|-------------|----------|---------|
| `--parent-id` | `-i` | ID of the parent page or database | Yes | — |
| `--content` | `-c` | Markdown content for the page body | Yes | — |
| `--parent-type` | `-p` | `page` or `database` | No | `page` |
| `--title` | `-t` | Title for the new page | No | `""` |
| `--max-retries` | `-m` | Max retries for failed API calls | No | `3` |
| `--retry-interval` | `-r` | Retry interval in seconds | No | `1` |

**`append`**

| Flag | Alias | Description | Required | Default |
|------|-------|-------------|----------|---------|
| `--content` | `-c` | Markdown content to append | Yes | — |
| `--max-retries` | `-m` | Max retries for failed API calls | No | `3` |
| `--retry-interval` | `-r` | Retry interval in seconds | No | `1` |

---

## Step 2 - Execute and Report

Once the CLI command has been produced, ask for permission to execute and, if the permission is granted, run the command and report all what you did to the user.

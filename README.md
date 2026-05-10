# skpm CLI

> The command-line tool for publishing Skript packages to the skpm registry.

Install scripts. Share scripts. Keep servers in sync.

---

## Installation

```sh
curl -fsSL https://raw.githubusercontent.com/skpm-dev/cli/main/install.sh | sh
```

Downloads the latest release and places `skpm` in `/usr/local/bin`.

**Windows:** grab the binary from the [releases page](https://github.com/skpm-dev/cli/releases).

## Setup

Create a GitHub personal access token with `read:user` scope and export it:

```bash
export SKPM_GITHUB_TOKEN=ghp_your_token_here
```

Add that line to your `~/.zshrc` or `~/.bashrc` to persist it across sessions.

---

## Publishing a package

### 1. Create `skpm.json`

```json
{
  "name": "my-economy",
  "description": "A simple economy system for Skript servers",
  "author": "yourgithubusername",
  "version": "1.0.0",
  "repo": "https://github.com/you/my-economy",
  "skript": ">=2.8.0",
  "minecraft": ">=1.20",
  "addons": {},
  "files": [
    "economy.sk"
  ]
}
```

| Field | Required | Description |
|---|---|---|
| `name` | Yes | Unique package name — used in `/skpm install <name>` |
| `description` | Yes | Short description of what the package does |
| `author` | Yes | Your GitHub username |
| `version` | Yes | Semantic version (`major.minor.patch`) |
| `repo` | Yes | Link to the source repository |
| `skript` | No | Minimum Skript version required (e.g. `>=2.8.0`) |
| `minecraft` | No | Minimum Minecraft version required (e.g. `>=1.20`) |
| `addons` | No | Required Skript addons and their minimum versions |
| `files` | Yes | `.sk` files to include in the package |

### 2. Publish

```bash
skpm publish
```

First publish uses the version in `skpm.json`. On subsequent publishes, skpm prompts you to choose a bump:

```
Found existing package at version 1.0.0
What type of release is this?
  [1] patch — bug fixes        (1.0.0 → 1.0.1)
  [2] minor — new features     (1.0.0 → 1.1.0)
  [3] major — breaking changes (1.0.0 → 2.0.0)

Enter choice:
```

skpm updates `skpm.json` with the new version and opens a pull request on the registry. Once a maintainer merges it, the package is live.

### 3. Validate without publishing

```bash
skpm validate
```

Checks that `skpm.json` is valid and all required fields are present.

---

## Ownership

The first account to publish a package name owns it. All future publishes verify your GitHub identity — if your username doesn't match the stored author, the publish is rejected.

---

## Related

- [skpm-dev/plugin](https://github.com/skpm-dev/plugin) — Bukkit plugin that installs packages in-game
- [skpm-dev/registry](https://github.com/skpm-dev/registry) — package registry and API

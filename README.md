# skpm CLI

**[skpm.org](https://skpm.org)** — the package manager for Skript.

> Scaffold, validate, and publish Skript packages to the skpm registry from your terminal.

---

## Installation

```sh
curl -fsSL https://raw.githubusercontent.com/skpm-dev/cli/main/install.sh | sh
```

Installs the `skpm` binary to `/usr/local/bin`. **Windows:** download from the [releases page](https://github.com/skpm-dev/cli/releases).

---

## Setup

Create a GitHub personal access token with **`read:user` scope**, then export it:

```bash
export SKPM_GITHUB_TOKEN=ghp_your_token_here
```

Add that line to `~/.zshrc` or `~/.bashrc` to persist it across sessions.

---

## Commands

| Command | Description |
|---|---|
| `skpm init` | Scaffold a `skpm.json` in the current directory |
| `skpm validate` | Validate `skpm.json` and confirm all files exist on disk |
| `skpm publish` | Publish the current package to the registry |
| `skpm info <package>` | Show metadata and version history for a package |
| `skpm search <query>` | Search the registry by name or description |
| `skpm completion <shell>` | Generate shell completion (bash, zsh, fish, powershell) |

---

## Publishing a package

### 1. Create `skpm.json`

Run `skpm init` or write it manually:

```json
{
  "name": "my-economy",
  "description": "A simple economy system for Skript servers",
  "author": "yourgithubusername",
  "version": "1.0.0",
  "skript": ">=2.8.0",
  "minecraft": ">=1.20",
  "addons": {
    "skript-reflect": ">=2.4.0"
  },
  "files": [
    "economy.sk"
  ]
}
```

| Field | Required | Notes |
|---|---|---|
| `name` | **Yes** | Lowercase, letters/digits/hyphens, 3–39 chars. Used in `/skpm install <name>` |
| `description` | **Yes** | Short summary shown in search results |
| `author` | **Yes** | Your GitHub username — tied to ownership |
| `version` | **Yes** | Semver (`major.minor.patch`) |
| `skript` | No | Semver constraint on required Skript version |
| `minecraft` | No | Semver constraint on required Minecraft version |
| `addons` | No | Map of required Skript addon → semver constraint |
| `dependencies` | No | Map of required skpm package → semver constraint |
| `files` | **Yes** | `.sk` files to include — must exist on disk |

### 2. Validate locally

```bash
skpm validate
```

Runs the same name, version, and constraint checks the registry enforces — catches errors before the network round-trip.

### 3. Publish

```bash
skpm publish
```

**First publish** uses the version in `skpm.json` as-is.

**Subsequent publishes** prompt for a version bump:

```
Found existing package at version 1.0.0.
What type of release is this?
  [1] patch — bug fixes        (1.0.0 → 1.0.1)
  [2] minor — new features     (1.0.0 → 1.1.0)
  [3] major — breaking changes (1.0.0 → 2.0.0)
```

skpm opens a pull request on the registry repo. Once a maintainer merges it, the package is live. `skpm.json` is updated with the new version automatically on success.

Pass `--bump patch|minor|major` to skip the prompt in CI.

---

## Ownership

- The **first account** to publish a name owns it.
- All subsequent publishes verify your GitHub identity via `SKPM_GITHUB_TOKEN`.
- Publishing as a different user than the stored author returns **403 Forbidden**.

---

## Environment variables

| Variable | Required | Description |
|---|---|---|
| `SKPM_GITHUB_TOKEN` | **Yes** | GitHub PAT with `read:user` scope — used for publish auth |

---

## Related

- **[skpm-dev/plugin](https://github.com/skpm-dev/plugin)** — Bukkit plugin that installs packages in-game
- **[skpm-dev/registry](https://github.com/skpm-dev/registry)** — Registry API and data store

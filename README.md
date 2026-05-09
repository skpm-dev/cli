# skpm CLI

The command-line tool for publishing Skript packages to the [skpm registry](https://github.com/skpm-dev/registry).

## Prerequisites

- A GitHub account
- A GitHub personal access token with `read:user` scope

Set your token as an environment variable:

```bash
export SKPM_GITHUB_TOKEN=ghp_your_token_here
```

Add this to your `~/.zshrc` or `~/.bashrc` to make it permanent.

## Installation

```sh
curl -fsSL https://raw.githubusercontent.com/skpm-dev/cli/main/install.sh | sh
```

This downloads the latest release for your OS and installs `skpm` to `/usr/local/bin`. You may be prompted for your password.

Windows users: download the binary from the [releases page](https://github.com/skpm-dev/cli/releases).

## Publishing a package

**1. Create a `skpm.json` in your project directory:**

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
| `name` | Yes | Unique package name (used in `/skpm install <name>`) |
| `description` | Yes | Short description of what the package does |
| `author` | Yes | Your GitHub username |
| `version` | Yes | Semantic version (`major.minor.patch`) |
| `repo` | Yes | Link to the source repository |
| `skript` | No | Minimum Skript version required |
| `minecraft` | No | Minimum Minecraft version required |
| `addons` | No | Map of required Skript addons and their versions |
| `files` | Yes | List of `.sk` files to include |

**2. Run publish:**

```bash
skpm publish
```

If this is your first publish, the version from `skpm.json` is used as-is.

If the package already exists, skpm will prompt you to choose a version bump:

```
Found existing package at version 1.0.0
What type of release is this?
  [1] patch — bug fixes        (1.0.0 → 1.0.1)
  [2] minor — new features     (1.0.0 → 1.1.0)
  [3] major — breaking changes (1.0.0 → 2.0.0)

Enter choice:
```

skpm updates `skpm.json` with the new version, then opens a pull request on the registry repo. A maintainer reviews and merges it. Once merged, the package is live.

**3. Validate without publishing:**

```bash
skpm validate
```

Checks that `skpm.json` is present and all required fields are filled in.

## Ownership

The first person to publish a package name owns it. Subsequent publishes verify your GitHub identity — if your GitHub username does not match the stored author, the publish is rejected.

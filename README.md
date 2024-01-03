<div align="center">

  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/loghinalexandru/loghinalexandru.github.io/master/static/img/anchor_banner_dark.png">
    <img src="https://raw.githubusercontent.com/loghinalexandru/loghinalexandru.github.io/master/static/img/anchor_banner_light.png" width="25%">
  </picture>

</div>

<div align="center">

  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/loghinalexandru/loghinalexandru.github.io/master/static/img/anchor_title_dark.png">
    <img src="https://raw.githubusercontent.com/loghinalexandru/loghinalexandru.github.io/master/static/img/anchor_title_light.png">
  </picture>

</div>

</br>

<div align="center">

  ![Go Report Card](https://goreportcard.com/badge/github.com/loghinalexandru/anchor)
  ![CI](https://github.com/loghinalexandru/anchor/actions/workflows/ci.yaml/badge.svg)
  ![GitHub release (with filter)](https://img.shields.io/github/v/release/loghinalexandru/anchor)

</div>

---
**Anchor** is a simple CLI tool for storing bookmarks neatly organized in a flattened label hierarchy. It features a custom and intuitive TUI via [bubbletea](https://github.com/charmbracelet/bubbletea) and automatic syncing with git via [go-git](https://github.com/go-git/go-git) as a backing storage.

<img src="https://raw.githubusercontent.com/loghinalexandru/loghinalexandru.github.io/master/static/img/anchor_demo.gif">

# Installation

Download one of the [pre-built binaries](https://github.com/loghinalexandru/anchor/releases/latest) and make it available for the shell you are using. Make sure to grab the right one for the operating system/architecture you intend to use it on.

# Build from source

Prerequisites to build anchor from source:

- Go 1.21 or later

Build and place it under ```$GOBIN```:

```text
 go install github.com/loghinalexandru/anchor@latest
```

# Storage

For now by default it uses the local file system as storage. You can specify what kind of backing storage you want via a config file explained in the next section.

Valid options:

- local
- git via ssh auth

# Usage

In order to use **anchor** you first need to create a home for all your bookmarks. Before any operation you need to initialize the storage. For this run:

```text
anchor init
```

In order to switch to **git** or any other storage, create a file under ```~/.anchor/config/anchor.yaml``` with the following config:

```yaml
storage: git
```

For this to work you need to have a repository already created and a **ssh** key already setup. The authentication is done via the **ssh-agent** as mentioned in the [go-git](https://github.com/go-git/go-git) documentation.

This is a one-time-only configuration since the file will be persisted on the preffered backing storage.

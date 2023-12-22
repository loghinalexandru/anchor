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
**Anchor** is a simple CLI tool for managing bookmarks for all platforms. It supports a custom and intuitive TUI via [bubbletea](https://github.com/charmbracelet/bubbletea) and automatic syncing with local or [git](https://git-scm.com/) as a back-end storage.

# Installation

Download one of the [pre-built binaries](https://github.com/loghinalexandru/anchor/releases/latest) and make it available for the shell you are using. Make sure to grab the right one for the operating system/architecture you intend to use it on.

# Build from source

Prerequisites to build anchor from source:

- Go 1.21 or later

Build and place it under ```$GOBIN```:

``` go
 go install github.com/loghinalexandru/anchor@latest
```

# Storage

For now by default it uses the local file system as storage. You can specify what kind of back-end storage you want via a config file explained in the next section.

Valid options:

- local
- git via ssh auth

# Usage

In order to use **Anchor** you first need to create a home for all your bookmarks. Before any operation you need to initialize the storage. For this run:

``` bash
anchor init
```


In order to switch to **git** or any other storage, create a file under ```~/.anchor/config/anchor.yaml``` with the following config:

```yaml
storage: git
```

This is a one-time-only since the file will be persisted on the preffered back-end storage.
<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/loghinalexandru/loghinalexandru.github.io/master/static/img/anchor_banner_dark.png">
    <img src="https://raw.githubusercontent.com/loghinalexandru/loghinalexandru.github.io/master/static/img/anchor_banner_light.png" width="30%">
  </picture>
</div>
<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/loghinalexandru/loghinalexandru.github.io/master/static/img/anchor_title_dark.png">
    <img src="https://raw.githubusercontent.com/loghinalexandru/loghinalexandru.github.io/master/static/img/anchor_title_light.png" width="85%">
  </picture>
</div>

# Known Limitations
- Always run ```anchor sync``` before any mutating operations. The CLI uses [go-git](https://github.com/go-git/go-git), and it does not support merge conflicts resolution.

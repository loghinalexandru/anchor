<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="images/banner_dark.png">
    <img src="images/banner_light.png" width="30%">
  </picture>
  <h1>Anchor</h1>
</div>

# Known Limitations
- Always run ```anchor sync``` before any mutating operations. The CLI uses [go-git](https://github.com/go-git/go-git), and it does not support merge conflicts resolution.
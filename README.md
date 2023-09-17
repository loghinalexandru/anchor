# Known Limitations

- Always run ```anchor sync``` before any mutating operations. The CLI uses [go-git](https://github.com/go-git/go-git), and it does not support merge conflicts resolution.
# Embedded server payload

`deploy_windows.go` embeds `homecloud-linux-amd64` from this directory with
`//go:embed payload/homecloud-linux-amd64`.

The binary is a build artifact and is **not committed** (see `.gitignore`).
A fresh clone therefore cannot compile the client until the payload exists.

Regenerate it before the first Windows build:

```powershell
.\build.ps1
```

`build.ps1` cross-compiles `./cmd` for `linux/amd64`, writes the result here,
then builds the Windows client. Run it whenever the server code changes.

# HomeCloud desktop UX and binding contract

Assets are dependency-free HTML/CSS/JS, embedded by Wails (`frontend/*`). The native frame stays visible; the backend opens a centered 900×600 resizable window (minimum 640×460). No fixed overlays or custom title bar. Tabs support arrows/Home/End, responsive scrolling and visible keyboard focus. Opening index.html directly is a labelled static preview, never simulated live data.

## UX
- **Language:** a global Русский/English selector in the header. Default is Russian (owner language). Choice is persisted in `localStorage` only (`homecloud.language`), never sent to the server. All owned static labels and dynamic notices/errors/confirmations re-render immediately via `setLanguage`; raw SSH/Docker/HTTP errors pass through `normalize` (redacted) without translation so diagnostics stay accurate. Passwords and fingerprints are never included in translation keys. Tests cover default RU, immediate switch, and persistence/restore.
- **SSH:** auto-scan mDNS with manual host/ports override. Discovery never grants trust. Probe omits the password; host/SSH-port edits invalidate fingerprint approval. A separate checkbox requires independent fingerprint comparison. Connect and deployment are separate forms; deployment confirms installation and uses a separate instance port. Password clears after each authentication attempt; no credentials or preferences are persisted.
- **Dashboard:** refreshable CPU/memory snapshot and top 20 processes. Values are server data, not demo data.
- **Containers:** list, start, stop, restart; stop/restart require confirmation. No invented deletion endpoint.
- **Files:** browse real absolute paths over the SSH/SFTP connection (not the sandboxed WebDAV share), with breadcrumbs, a Home shortcut, recursive in-folder search (bounded depth/results), and native upload/download dialogs. Upload confirms overwrite risk. Drive mapping explains the actual Windows Bearer-auth limitation instead of opening an unsafe proxy.
- **Terminal:** an interactive SSH shell inside the tab (xterm.js). Keystrokes, including Enter, arrows, Tab and Ctrl+C, go straight to the PTY. The view keeps ANSI colour, cursor movement and a 5000-line scrollback, and resizes the remote PTY to the pane. A destructive-looking line (`rm`, `mkfs`, `dd`, shutdown and similar) still asks for confirmation before Enter is sent.
- **Remote Desktop:** launches a real RDP session with `mstsc.exe`, writing a temporary `.rdp` and seeding the login through Windows Credential Manager (`cmdkey /generic:TERMSRV/<host>`). A TCP port check and an SSH-based xrdp probe report reachability before connecting. Saved devices (host, ports, user, optional DPAPI-encrypted password) can be listed, loaded, saved and deleted; `ListDevices` never returns secrets. A separate view-only screenshot panel remains available over SSH.
- Errors normalize to readable bounded text. Known Bearer/token patterns are redacted. Token status is boolean only; token values never enter frontend state. Unavailable operations become in-memory reminders, never automatically replayed. No raw network HTTP requests are made by the frontend.

## Wails binding namespace
`window.go.client.App` (package `client`, exported bound `App` in client.go). Wails injects the namespace; no generated wrapper import required. Events use `window.runtime.EventsOn`.

| UI | Method / arguments | Result / remote route |
|---|---|---|
| Discovery | `Discover()` | `[{name,host,port}]`; port is HomeCloud HTTP port, not SSH |
| Probe | `ProbeHostKey(SSHConfig)` | untrusted SHA256 fingerprint |
| Connect | `Connect(SSHConfig)` | ConnectionStatus |
| Deploy | `Deploy(SSHConfig)` | human-readable startup result |
| Status | `Status()` | ConnectionStatus |
| Disconnect | `Disconnect()` | void |
| Stats | `Request('/api/v1/stats','')` | GET JSON `{cpu_percent:number[],memory:{used,total,usedPercent,...}}` |
| Processes | `Request('/api/v1/processes','')` | GET JSON `[{pid,name,memory_bytes}]` |
| Containers | `Request('/api/v1/containers','')` | GET Docker container array `{Id,Names,Image,State,Status,...}` |
| Container action | `Request('/api/v1/containers/control', JSON.stringify({container_id,action}))` | POST; action start/stop/restart; `{status:'ok'}` |
| Browse | `ListFiles(path)` | SFTP `ReadDir`; `[{name,path,isDir,size,modified}]`, absolute paths |
| Home | `HomeDir()` | SFTP `Getwd`; absolute remote home directory |
| Search | `SearchFiles(dir, query)` | recursive SFTP walk (depth 8, 500 results, 50k visited); `[{name,path,isDir,size,modified}]` |
| Upload | `UploadFile(directory)`, `UploadFolder(directory)` | native dialog then SFTP copy into the folder shown in the confirm dialog; the filesystem root is redirected to the login home because it is not writable; a folder is recreated with its subfolders and only regular files are copied (symlinks are skipped); existing names are overwritten; returns `[remotePath...]` |
| Download | `DownloadFile(filePath)` | native folder dialog + SFTP read; returns the saved local path |
| Devices | `ListDevices()`, `GetDevice(id)`, `SaveDevice(device)`, `DeleteDevice(id)` | local SQLite `%LOCALAPPDATA%\HomeCloud\client.db`; passwords DPAPI-encrypted; only `GetDevice` returns the decrypted password. A successful `Connect` remembers the verified profile (host, SSH port, panel port, user, fingerprint, password) and marks it most recently used. On startup the most recently used device is loaded into both the SSH and RDP forms with its stored fingerprint re-checked. Saving from the RDP tab keeps the stored SSH ports and fingerprint unless the host changed. |
| RDP check | `RDPCheck(RDPRequest)` | local TCP probe of the RDP port |
| RDP probe | `RDPProbe(port)` | SSH check for xrdp/sesman and a listening port |
| RDP connect | `RDPConnect(RDPRequest)` | writes `.rdp`, seeds `cmdkey`, launches `mstsc.exe` |
| RDP forget | `RDPForget()` | `cmdkey /delete:TERMSRV/<host>` for the last RDP host |
| Map | `MountDrive(letter)` | currently rejects: Windows WebClient cannot use Bearer auth |
| Terminal | `OpenTerminal()`, `TerminalInput(data)`, `ResizeTerminal(cols,rows)`, `CloseTerminal()` | SSH PTY |

`SSHConfig = {host,port,serverPort,user,password,fingerprint}`.
`ConnectionStatus = {connected,host,serverPort,tokenReady,terminalOpen,message}`.
`RDPRequest = {host,port,user,password,width,height,fullscreen,remember}`.
`Device = {id,name,host,port,serverPort,user,fingerprint,rdpPort,hasPassword,password,clearPassword}`.
Request returns JSON **text**, parsed by the frontend; Go errors reject the promise. Health (`/api/v1/health`) is also whitelisted but not separately polled. No SSE HTTP endpoint exists: terminal uses native events `terminal:output`, `terminal:closed`, and `connection:closed`, each with a string payload.

## Verification
Run `node --check internal/client/frontend/app.js`, `node --test internal/client/frontend/frontend.test.cjs`, `go test ./internal/client`, and a Windows production build of `./cmd` with Wails production tags. Runtime GUI interaction still requires WebView2 and a real SSH host; DOM harness tests are not a substitute for that integration test.
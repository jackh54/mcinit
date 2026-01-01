# mcinit

A cross-platform Minecraft developer CLI tool that creates and manages local dev servers quickly and reproducibly.

## Features

- 🚀 **Fast Setup**: Create a dev server with a single command
- 🔄 **Reproducible**: Team-wide consistency via `mcinit.json` config
- 💻 **Cross-Platform**: Works identically on Windows, macOS, and Linux
- 🎯 **Developer-Focused**: Plugin linking, auto-restart, and workspace integration
- 📦 **Multiple Server Types**: Vanilla, Paper, Purpur, Folia, Velocity, Waterfall, BungeeCord

## Installation

### Homebrew (macOS/Linux)

```bash
brew install mcinit
```

### Scoop (Windows)

```bash
scoop install mcinit
```

### Winget (Windows)

```bash
winget install mcinit
```

### Manual Installation

Download the latest release from [GitHub Releases](https://github.com/jackh54/mcinit/releases).

## Quick Start

Create a new Paper server:

```bash
mcinit init --type paper --mc 1.21.4 --accept-eula --ram 4G
cd server
mcinit start
```

## Usage

### Initialize a Server

```bash
mcinit init --type paper --mc 1.21.4 --accept-eula --ram 4G --path ./my-server
```

### Start/Stop Server

```bash
mcinit start
mcinit start --background
mcinit stop
mcinit restart
```

### View Logs

```bash
mcinit logs
mcinit logs --follow
mcinit logs --grep "error"
```

### Check Status

```bash
mcinit status
```

## Configuration

After running `init`, a `mcinit.json` file is created in your server directory. Commit this file to version control for reproducible setups.

Example `mcinit.json`:

```json
{
  "version": "1.0.0",
  "server": {
    "type": "paper",
    "minecraftVersion": "1.21.4",
    "build": "latest",
    "jarPath": "server.jar",
    "name": "dev-server"
  },
  "java": {
    "version": "auto",
    "path": "auto"
  },
  "jvm": {
    "xms": "2G",
    "xmx": "4G",
    "flags": "aikar"
  },
  "serverConfig": {
    "port": 25565,
    "nogui": true,
    "maxPlayers": 20,
    "onlineMode": false,
    "difficulty": "easy"
  }
}
```

## Supported Server Types

- **vanilla**: Official Mojang server
- **paper**: PaperMC (most common for plugins)
- **purpur**: Performance-focused Paper fork
- **folia**: Multi-threaded Paper fork
- **velocity**: Modern proxy server
- **waterfall**: BungeeCord fork
- **bungee**: Classic proxy server

## Development

### Building from Source

```bash
git clone https://github.com/jackh54/mcinit.git
cd mcinit
make build
```

### Running Tests

```bash
make test
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please open an issue or PR.


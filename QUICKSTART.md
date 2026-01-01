# Quick Start Guide

## Installation

### Via Homebrew (macOS/Linux)
```bash
brew install mcinit
```

### Via Scoop (Windows)
```bash
scoop install mcinit
```

### Via Winget (Windows)
```bash
winget install mcinit
```

### Manual Installation

**Unix (macOS/Linux)**:
```bash
curl -sSL https://raw.githubusercontent.com/jackh54/mcinit/main/scripts/install.sh | bash
```

**Windows (PowerShell)**:
```powershell
irm https://raw.githubusercontent.com/jackh54/mcinit/main/scripts/install.ps1 | iex
```

### Build from Source
```bash
git clone https://github.com/jackh54/mcinit.git
cd mcinit
make build
sudo cp bin/mcinit /usr/local/bin/
```

## Quick Start: Create Your First Server

### 1. Initialize a Paper Server (Most Common)
```bash
mcinit init --type paper --mc 1.21.4 --accept-eula --ram 4G
```

This will:
- ✅ Download the Paper server jar for Minecraft 1.21.4
- ✅ Create `mcinit.json` configuration
- ✅ Generate startup scripts (shell/PowerShell/CMD)
- ✅ Accept the EULA automatically
- ✅ Configure 4GB RAM with Aikar's optimized JVM flags
- ✅ Detect and validate your Java installation

### 2. Start the Server
```bash
cd server
mcinit start
```

Or run in background:
```bash
mcinit start --background
```

### 3. Monitor Server
```bash
# Check status
mcinit status

# View logs
mcinit logs

# Follow logs in real-time
mcinit logs --follow

# Filter logs
mcinit logs --grep "error" --lines 100
```

### 4. Manage Server
```bash
# Stop server
mcinit stop

# Restart server
mcinit restart

# Force stop if needed
mcinit stop --force
```

## Common Use Cases

### Vanilla Server
```bash
mcinit init --type vanilla --mc 1.21.4 --accept-eula --ram 2G
```

### Purpur Server (Performance Fork)
```bash
mcinit init --type purpur --mc 1.21.4 --accept-eula --ram 6G --nogui
```

### Test Server with Specific Java
```bash
mcinit init --type paper --mc 1.21.4 --java 21 --accept-eula --ram 4G
```

### Development Server with Minimal Flags
```bash
mcinit init --type paper --mc 1.21.4 --flags minimal --accept-eula --ram 4G --gitignore
```

### Custom Port
```bash
mcinit init --type paper --mc 1.21.4 --port 25566 --accept-eula --ram 4G
```

### Custom RAM Configuration
```bash
mcinit init --type paper --mc 1.21.4 --xms 2G --xmx 8G --accept-eula
```

## Reproducible Team Setup

1. **First teammate** initializes server:
```bash
mcinit init --type paper --mc 1.21.4 --accept-eula --ram 4G
git add server/mcinit.json
git commit -m "Add Minecraft server config"
git push
```

2. **Other teammates** clone and set up:
```bash
git clone <repo-url>
cd <repo>/server
# mcinit reads mcinit.json and recreates the exact setup
mcinit start
```

## Configuration File

The `mcinit.json` file stores all server settings:
```json
{
  "version": "1.0.0",
  "server": {
    "type": "paper",
    "minecraftVersion": "1.21.4",
    "jarPath": "server.jar"
  },
  "java": {
    "version": "auto",
    "path": "auto"
  },
  "jvm": {
    "xms": "4G",
    "xmx": "4G",
    "flags": "aikar"
  },
  "serverConfig": {
    "port": 25565,
    "nogui": true
  }
}
```

Commit this file to version control for reproducible setups!

## Supported Server Types

| Type | Description | Source |
|------|-------------|--------|
| `vanilla` | Official Mojang server | Mojang API |
| `paper` | Most popular plugin server | PaperMC API v2 |
| `purpur` | Performance-focused fork | Purpur API |
| `folia` | Multi-threaded Paper fork | PaperMC API v2 |
| `velocity` | Modern proxy server | PaperMC API v2 |
| `waterfall` | BungeeCord fork | PaperMC API v2 |

## JVM Flags Presets

### Aikar's Flags (Default)
Optimized garbage collection settings recommended by PaperMC:
```bash
mcinit init --type paper --mc 1.21.4 --flags aikar --accept-eula --ram 4G
```

### Minimal Flags
Basic G1GC configuration:
```bash
mcinit init --type paper --mc 1.21.4 --flags minimal --accept-eula --ram 4G
```

### Custom Flags
Provide your own JVM arguments:
```bash
mcinit init --type paper --mc 1.21.4 --flags custom --accept-eula --ram 4G
```

## Troubleshooting

### Java Not Found
```bash
# Specify Java path explicitly
mcinit init --type paper --mc 1.21.4 --java /path/to/java --accept-eula --ram 4G

# Or install Java 17/21
```

### Server Won't Start
```bash
# Check logs
mcinit logs --grep "error"

# Verify EULA is accepted
cat server/eula.txt

# Check Java version
mcinit status
```

### Port Already in Use
```bash
# Use a different port
mcinit init --type paper --mc 1.21.4 --port 25566 --accept-eula --ram 4G
```

## Next Steps

- Read the full [README](README.md) for detailed documentation
- Check out [CONTRIBUTING.md](CONTRIBUTING.md) to contribute
- Report issues on [GitHub Issues](https://github.com/jackh54/mcinit/issues)
- Join the community discussions

## Getting Help

```bash
# General help
mcinit --help

# Command-specific help
mcinit init --help
mcinit start --help
mcinit logs --help
```

---

**Happy Minecrafting! 🎮**


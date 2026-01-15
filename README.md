# Loggar.dev

**AI-powered log triage for Mac Terminal**

Loggar.dev helps developers quickly analyze and triage server logs, identifying root causes, secondary effects, and recommended actions—all from a terminal interface.

## 🚀 Features

- **AI-Powered Analysis**: Uses Google AI Studio to analyze logs and identify issues
- **Terminal-First**: Beautiful, color-coded output designed for the command line
- **Fast Authentication**: Simple JWT-based auth with local token storage
- **Flexible Input**: Analyze logs from files or stdin
- **Easy Installation**: Install via Homebrew or curl script

## 📦 Installation

### Homebrew (Coming Soon)

```bash
brew install loggar-dev/tap/loggar
```

### Curl Script (Coming Soon)

```bash
curl -fsSL https://loggar.dev/install.sh | bash
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/ayomide/loggar.git
cd loggar

# Install dependencies
go mod download

# Build the CLI
go build -o triage ./cmd/loggar

# Move to PATH (optional)
sudo mv triage /usr/local/bin/
```

## 🎯 Quick Start

### 1. Authenticate

```bash
triage auth
```

This opens your browser for login and saves your JWT token locally.

### 2. Analyze Logs

```bash
# Analyze a log file
triage analyze server.log

# Or pipe logs via stdin
cat error.log | triage analyze

# Get raw JSON output
triage analyze server.log --json
```

### Example Output

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PRIMARY ISSUE
→ Database connection pool exhaustion
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
SECONDARY EFFECTS
• Auth service timeouts
• Payment retries failing
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
LIKELY CAUSES
1. Unreleased DB connections (63%)
2. Traffic spike exceeded pool size (27%)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
RECOMMENDED ACTIONS
• Check connection release in auth middleware
• Inspect pool max size vs current RPS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## 🛠️ Development Setup

### Prerequisites

- Go 1.20+
- PostgreSQL
- Google AI Studio API key

### Setup

1. **Clone and install dependencies**

```bash
git clone https://github.com/ayomide/loggar.git
cd loggar
go mod download
```

2. **Set up PostgreSQL**

```bash
# Create database
createdb loggar

# Run schema
psql loggar < scripts/setup-db.sql
```

3. **Configure environment variables**

```bash
cp .env.example .env
# Edit .env with your configuration
```

4. **Run the backend server**

```bash
go run api/server.go
```

5. **Build and test the CLI**

```bash
go build -o triage ./cmd/loggar
./triage --help
```

## 📚 Commands

| Command | Description |
|---------|-------------|
| `triage auth` | Authenticate and save token |
| `triage analyze <file>` | Analyze log file |
| `triage analyze --stdin` | Analyze from stdin |
| `triage version` | Show CLI version |
| `triage help` | Show help message |

## 🏗️ Architecture

```
┌────────────┐       ┌────────────┐       ┌─────────────┐
│  CLI (Go)  │ <---> │ Backend API│ <---> │ Google AI   │
│   triage   │       │ Gin/Fiber  │       │ Studio API  │
└────────────┘       └────────────┘       └─────────────┘
       │
       ▼
 ~/.loggar/config.json
```

## 🎯 Roadmap

### v1 (Free Tier) ✅
- [x] CLI with auth and analyze commands
- [x] Backend API with JWT authentication
- [x] Google AI Studio integration
- [x] Pretty terminal output
- [ ] Homebrew installation
- [ ] Landing page

### v2 (Paid Tier)
- [ ] Usage tracking and analytics
- [ ] Postmortem generation
- [ ] Team incident sharing
- [ ] Web dashboard
- [ ] History search

## 📄 License

MIT

## 🤝 Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## 📧 Contact

For questions or support, visit [loggar.dev](https://loggar.dev)

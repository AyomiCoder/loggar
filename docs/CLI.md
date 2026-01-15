# logger

ai-powered log triage. terminal-first.

## Installation

### Instant Install
Install the Loggar CLI with a single command:
```bash
curl -fsSL https://raw.githubusercontent.com/AyomiCoder/loggar/main/scripts/install.sh | sh
```

---

## Configuration

The CLI stores authentication tokens in `~/.loggar/config.json`. This file is created automatically when you authenticate.

---

## Commands

### 1. Version
Display the CLI version.
```bash
loggar version
```
**Sample response:**
```
Loggar CLI v0.1.0
```

### 2. Help
Show help information for any command.
```bash
loggar --help
```
**Sample response:**
```
Loggar.dev - Analyze server logs and identify root causes with AI

Usage:
  loggar [command]

Available Commands:
  analyze     Analyze log files with AI
  auth        Authenticate with Loggar.dev
  help        Help about any command
  version     Print the version number
```

### 3. Authentication
Authenticate with the Loggar backend. If you don't have an account, create one first.

#### Sign up
```bash
loggar signup
```
**Sample response:**
```
Create a new Loggar.dev account
Email: new@logger.dev
Password (min 6 chars): ********
Confirm Password: ********

✓ Account created successfully! You can now run 'loggar auth' to login.
```

#### Login
```bash
loggar auth
```
**Sample response:**
```
Login to Loggar.dev
Email: dev@logger.dev
Password: ********

✓ Successfully authenticated as dev@logger.dev
Token saved to /Users/ayomide/.loggar/config.json
```

#### Reset token
```bash
loggar auth --reset
```
**Sample response:**
```
✓ Token cleared successfully
```

### 5. Analyze Logs
Analyze log files or stdin input to identify issues and get AI-powered recommendations.

#### From file:
```bash
loggar analyze server.log
```
**Sample response:**
```bash
💡 Application failed to start because it couldn't reach the database on port 5432.

WHATS BROKEN
• connection refused on 127.0.0.1:5432
• postgresql health check failed

ROOT CAUSE
1. postgresql service is likely stopped or crashed (95%)
2. network config error or local firewall blocking port (15%)

HOW TO FIX
• check postgres status with systemctl or brew services
• verify port 5432 is listening
```

#### From stdin:
```bash
cat error.log | loggar analyze
```

**Sample response:**
```bash
💡 Application failed to start because it couldn't reach the database on port 5432.

WHATS BROKEN
• connection refused on 127.0.0.1:5432
• postgresql health check failed

ROOT CAUSE
1. postgresql service is likely stopped or crashed (95%)
2. network config error or local firewall blocking port (15%)

HOW TO FIX
• check postgres status with systemctl or brew services
• verify port 5432 is listening
```

#### JSON output:
```bash
loggar analyze server.log --json
```
**Sample response:**
```json
{
  "summary": "Service crash due to Java OutOfMemoryError (Heap Space)",
  "sections": [
    {
      "title": "WHATS BROKEN",
      "content": [
        "Service failed to auto-restart",
        "Container resource limit reached (memory)"
      ]
    },
    {
      "title": "ROOT CAUSE",
      "content": [
        "Insufficient Java heap space configuration for current workload (100%)",
        "Container memory limits are set too low (90%)"
      ]
    },
    {
      "title": "HOW TO FIX",
      "content": [
        "Increase Java heap size configuration (-Xmx)",
        "Review and increase container memory limits",
        "Analyze memory usage patterns for potential leaks"
      ]
    }
  ]
}
```

#### Verbose mode:
```bash
loggar analyze server.log --verbose
```
**Sample response:**
```bash
→ Analyzing 323 bytes of logs...

💡 Application failed to start because it couldn't reach the database on port 5432.

WHATS BROKEN
• connection refused on 127.0.0.1:5432
• postgresql health check failed

ROOT CAUSE
1. postgresql service is likely stopped or crashed (95%)
2. network config error or local firewall blocking port (15%)

HOW TO FIX
• check postgres status with systemctl or brew services
• verify port 5432 is listening
```

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
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PRIMARY ISSUE
→ database connection failure (connection refused)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
SECONDARY EFFECTS
• auth service timeouts
• payment processing failures
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
LIKELY CAUSES
1. postgresql service is stopped or crashed (90%)
2. resource exhaustion on database host (70%)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
RECOMMENDED ACTIONS
• verify status of database service on port 5432
• check host cpu and memory utilization
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

#### From stdin:
```bash
cat error.log | loggar analyze
```

**Sample response:**
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PRIMARY ISSUE
→ database connection failure (connection refused)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
SECONDARY EFFECTS
• auth service timeouts
• payment processing failures
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
LIKELY CAUSES
1. postgresql service is stopped or crashed (90%)
2. resource exhaustion on database host (70%)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
RECOMMENDED ACTIONS
• verify status of database service on port 5432
• check host cpu and memory utilization
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

#### JSON output:
```bash
loggar analyze server.log --json
```
**Sample response:**
```json
{
  "primary_issue": "Service crash due to Java OutOfMemoryError (Heap Space)",
  "secondary_effects": [
    "Service failed to auto-restart",
    "Container resource limit reached (memory)"
  ],
  "first_seen": "2026-01-15T17:05:44Z",
  "likely_causes": [
    {
      "cause": "Insufficient Java heap space configuration for current workload",
      "confidence": 1.0
    },
    {
      "cause": "Container memory limits are set too low",
      "confidence": 0.9
    }
  ],
  "recommended_actions": [
    "Increase Java heap size configuration (-Xmx)",
    "Review and increase container memory limits",
    "Analyze memory usage patterns for potential leaks"
  ],
  "similar_past_incidents": []
}
```

#### Verbose mode:
```bash
loggar analyze server.log --verbose
```
**Sample response:**
```bash
→ Analyzing 323 bytes of logs...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PRIMARY ISSUE
→ database connection failure (connection refused)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
SECONDARY EFFECTS
• auth service timeouts
• payment processing failures
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
LIKELY CAUSES
1. postgresql service is stopped or crashed (90%)
2. resource exhaustion on database host (70%)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
RECOMMENDED ACTIONS
• verify status of database service on port 5432
• check host cpu and memory utilization
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

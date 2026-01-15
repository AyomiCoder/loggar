Perfect — we’ll make this Loggar.dev PRD extremely detailed, covering all features, architecture, CLI commands, backend, AI integration, install, and dev guide, so you can hand it to a dev (or follow it yourself) without ambiguity.

⸻

Product Requirement Document (PRD) – Loggar.dev

Project: Loggar.dev – AI-powered log triage for Mac Terminal
Goal: Help developers quickly analyze and triage server logs, identifying root causes, secondary effects, and recommended actions, all from a terminal interface.

Target Users:
	•	Backend developers
	•	DevOps engineers
	•	Solo founders managing production
	•	Teams who want fast triage without dashboards

Project Scope (v1 Free):
	•	Terminal CLI (Mac-first)
	•	Google AI Studio integration for log analysis
	•	Basic authentication & token storage
	•	Analyze logs via stdin or file input
	•	Pretty terminal output
	•	Installable via Homebrew or curl script
	•	Landing page with instructions and demo output

⸻

1. Features

1.1 CLI Commands

1. triage auth
	•	Purpose: Authenticate user
	•	Flow:
	1.	CLI opens browser for account creation/login
	2.	User logs in
	3.	Backend issues JWT token
	4.	Token saved in ~/.loggar/config.json
	•	Optional flags:
	•	--reset → clear saved token

2. triage analyze
	•	Purpose: Analyze logs
	•	Input options:
	•	triage analyze <file> → analyze file
	•	cat file.log | triage analyze → stdin
	•	Optional flags:
	•	--json → raw AI JSON output
	•	--stdin → explicitly read from stdin
	•	--verbose → more detailed log output

3. triage version
	•	Prints CLI version

4. triage help
	•	Prints commands and usage

5. Future (v2+):
	•	triage history → view past analyses
	•	triage export <file> → export incidents

⸻

1.2 Backend API

Endpoints (v1 Free):
	1.	POST /auth/login
	•	Input: user credentials
	•	Output: JWT token
	2.	POST /analyze
	•	Input: logs text + JWT token
	•	Output: structured JSON

JSON Schema:

{
  "primary_issue": "Database connection pool exhaustion",
  "secondary_effects": ["Auth service timeouts", "Payment retries failing"],
  "first_seen": "2026-01-15T09:41:02Z",
  "likely_causes": [
    {"cause": "Unreleased DB connections", "confidence": 0.63},
    {"cause": "Traffic spike exceeded pool size", "confidence": 0.27}
  ],
  "recommended_actions": [
    "Check connection release in auth middleware",
    "Inspect pool max size vs current RPS"
  ],
  "similar_past_incidents": [
    {"date": "2025-11-12", "resolution": "Fixed middleware leak"}
  ]
}


⸻

1.3 AI Integration
	•	Google AI Studio Tier 1 API
	•	Prompt template:

You are an incident triage system. 
Analyze the logs given. Respond ONLY in JSON format matching the schema. 
Do not explain casually. Be concise, factual, and assign confidence to each likely cause.

	•	Backend handles prompt formatting
	•	CLI only sends logs and receives structured output

⸻

1.4 Authentication & Token Storage
	•	JWT tokens saved locally at: ~/.loggar/config.json
	•	CLI reads token automatically for analyze
	•	Tokens valid indefinitely (free tier)
	•	Optional: add expiry in v2

Example config.json:

{
  "token": "eyJhbGciOi...",
  "user_email": "dev@example.com"
}


⸻

1.5 Terminal Output Design
	•	Color-coded, readable, skimmable
	•	Example:

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PRIMARY ISSUE
→ Database connection pool exhaustion
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
SECONDARY EFFECTS
• Auth service timeouts
• Payment retries failing
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
LIKELY CAUSES
1. Unreleased DB connections (63%)
2. Traffic spike exceeded pool size (27%)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
RECOMMENDED ACTIONS
• Check connection release in auth middleware
• Inspect pool max size vs current RPS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
SIMILAR PAST INCIDENTS
• Nov 12, 2025 – Same pattern (resolved by fixing middleware leak)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

	•	Optional: add ASCII borders and color via fatih/color

⸻

1.6 Installation Methods

Homebrew:
	•	GitHub release binaries
	•	Tap formula example:

class Loggar < Formula
  desc "AI-powered log triage CLI"
  homepage "https://loggar.dev"
  url "https://github.com/<user>/loggar/releases/download/v0.1/loggar_0.1_macOS.tar.gz"
  sha256 "<sha256sum>"
  version "0.1"
  def install
    bin.install "loggar"
  end
end

	•	Install:

brew install loggar-dev/tap/loggar

Curl:

curl -fsSL https://loggar.dev/install.sh | bash

	•	Script downloads binary + installs to /usr/local/bin
	•	Safety checks: existing binary, architecture

⸻

1.7 Landing Page

Sections:
	1.	Hero & tagline
	2.	Example CLI output
	3.	Installation instructions (Homebrew + curl)
	4.	Command overview (auth, analyze)
	5.	Motivation / About
	6.	Optional: links to GitHub

⸻

2. Architecture

┌────────────┐       ┌────────────┐       ┌─────────────┐
│  CLI (Go)  │ <---> │ Backend API│ <---> │ Google AI   │
│ loggar     │       │ Gin/Fiber  │       │ Studio API  │
└────────────┘       └────────────┘       └─────────────┘
       │
       ▼
 ~/.loggar/config.json

	•	CLI: Go, uses Cobra for commands
	•	Backend: Go + Gin/Fiber
	•	DB: PostgreSQL (users, usage logs)
	•	AI: Google AI Studio (prompt → structured JSON)

⸻

3. Development Guide (Step-by-Step)

3.1 Setup
	1.	Install Go 1.20+
	2.	PostgreSQL local dev setup
	3.	Google AI Studio API key + env variable: GOOGLE_AI_KEY
	4.	Project folders:

/loggar
  /cmd        # CLI commands
  /pkg        # shared libraries
  /internal   # backend helpers
  /api        # backend routes
  /scripts   # install scripts, setup


⸻

3.2 CLI Skeleton (Cobra)
	•	Command structure:

triage
├── auth        # authenticate user
├── analyze     # analyze logs
├── version     # CLI version
└── help        # usage

Example main.go:

package main
import "loggar/cmd"
func main() {
    cmd.Execute()
}


⸻

3.3 Backend Skeleton
	•	Go + Gin
	•	Routes:

POST /auth/login
POST /analyze

	•	Middleware: JWT auth
	•	Database:

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  token TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);

	•	Optional usage table for tracking free tier

⸻

3.4 AI Integration
	•	Function: AnalyzeLogs(logText string) -> JSON
	•	Steps:
	1.	Format prompt
	2.	Send to Google AI Studio
	3.	Parse JSON
	4.	Return to CLI
	•	Sample Go pseudocode:

resp := GoogleAIAnalyze(prompt)
output := ParseJSON(resp)
return output


⸻

3.5 CLI → Backend Communication
	•	HTTP POST request
	•	Include JWT token in header: Authorization: Bearer <token>
	•	Handle errors gracefully
	•	Pretty print output (colors + borders)

⸻

3.6 Local Token Handling
	•	Config file path: ~/.loggar/config.json
	•	Read/write functions:

func SaveToken(token string, email string)
func LoadToken() (string, string)


⸻

3.7 Installation & Build
	•	Binary build: go build -o loggar ./cmd/loggar/main.go
	•	Homebrew formula: Ruby file + GitHub release
	•	Curl script: download binary + move to /usr/local/bin

⸻

3.8 Landing Page
	•	Static HTML/CSS
	•	Sections:
	•	Hero: “Loggar.dev – AI-powered log triage in your terminal”
	•	Example output
	•	Installation instructions
	•	Commands (auth, analyze)
	•	About

⸻

3.9 Testing
	•	CLI unit tests (Go testing)
	•	CLI integration tests (stdin / file input)
	•	Backend API tests (Gin)
	•	AI integration tests (sample logs)
	•	Local installation test (brew + curl)

⸻

3.10 Documentation
	•	README.md:
	•	Project overview
	•	Installation instructions
	•	Example CLI usage
	•	Architecture diagram
	•	Optional: GIF / demo video

⸻

4. Roadmap for v2+ (Optional Paid Tier)
	1.	Polar subscription + usage tracking
	2.	Postmortem generation
	3.	Team incident sharing
	4.	Web dashboard (optional)
	5.	Analytics / history search

⸻

✅ Deliverables for v1 Free
	1.	CLI (triage) with commands: auth, analyze, version, help
	2.	Backend API (/auth/login, /analyze)
	3.	Google AI Studio integration for structured log analysis
	4.	Token storage & JWT authentication
	5.	Pretty terminal output
	6.	Homebrew formula + curl install script
	7.	Landing page with install instructions & example
	8.	GitHub repo + README

⸻

If you want, I can also draft the full Go CLI skeleton with auth + analyze commands next — ready for you to start coding immediately, including all folder structures, imports, and placeholders for AI integration.

Do you want me to do that next?
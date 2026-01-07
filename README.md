# wt

A CLI for managing isolated git worktrees for parallel Claude Code sessions.

## Installation

### Homebrew

```bash
brew install emilrex/tap/wt
```

### Go

```bash
go install github.com/emilrex/wt@latest
```

Make sure `~/go/bin` is in your PATH:

```bash
export PATH="$HOME/go/bin:$PATH"
```

## Usage

```bash
wt new [name] [-b branch]  # Create session and launch Claude Code
wt fg <session>            # Resume session
wt ls                      # List sessions
wt rm <session>            # Remove session
wt rm --all                # Remove all sessions
wt cd <session>            # Open shell in session directory
```

## How it works

Each session creates:
- A git worktree in `~/.wt/{repo}-{session}`
- A branch named `wt-{session}`

Sessions are isolated from each other, so Claude can work on multiple tasks in parallel without conflicts.

Session names support partial matching - `wt fg auth` will match `auth-feature` if it's the only match.

## Inspiration

This project was inspired by [claude-wt](https://github.com/jlowin/claude-wt).

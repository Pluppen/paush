# ABOUTME: Sourced by the vhs tapes to stage a clean demo environment.
# ABOUTME: Sandboxes HOME, puts the paush build on PATH, and creates a tidy demo repo.

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="$ROOT:$PATH"
export HOME=/tmp/demohome
export SHELL=/bin/zsh
rm -rf /tmp/demohome /tmp/demoproj
mkdir -p /tmp/demohome /tmp/demoproj
cd /tmp/demoproj
git init -q -b main
cat > README.md <<'EOF'
# zen-garden
A small tool for tending digital gardens.
EOF
printf 'module zen-garden\n\ngo 1.23\n' > go.mod
cat > main.go <<'EOF'
package main

func main() {}
EOF
git add .
git -c user.email=you@example.com -c user.name=you commit -qm "plant the garden"
clear

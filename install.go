// ABOUTME: Installs typo aliases (sl, gti, cta, ...) that trigger a breathing break.
// ABOUTME: Manages a marked alias block in the user's shell rc file, idempotently.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	blockStart = "# >>> paush typo aliases >>>"
	blockEnd   = "# <<< paush typo aliases <<<"
)

// typos maps common command misspellings to the command that was meant.
var typos = []struct{ typo, meant string }{
	{"sl", "ls"},
	{"dri", "dir"},
	{"cta", "cat"},
	{"gti", "git"},
	{"grpe", "grep"},
	{"pdw", "pwd"},
	{"mkdri", "mkdir"},
	{"celar", "clear"},
	{"exti", "exit"},
	{"suod", "sudo"},
	{"vmi", "vim"},
	{"pyhton", "python"},
	{"tial", "tail"},
	{"ehco", "echo"},
	{"tuoch", "touch"},
	{"whcih", "which"},
	{"maek", "make"},
	{"dokcer", "docker"},
	{"nmp", "npm"},
}

// aliasBlock renders the marked alias block pointing every typo at exe.
func aliasBlock(exe string) string {
	var b strings.Builder
	b.WriteString(blockStart + "\n")
	b.WriteString("# a mindful pause instead of a typo; remove with `paush uninstall`\n")
	// The typo subcommand ignores trailing arguments, so aliases survive
	// invocations like `sl -la` or `gti status`.
	for _, t := range typos {
		fmt.Fprintf(&b, "alias %s='%s typo' # %s\n", t.typo, exe, t.meant)
	}
	b.WriteString(blockEnd + "\n")
	return b.String()
}

// upsertBlock returns content with the marked block replaced, or appended
// after a blank line when no block exists yet.
func upsertBlock(content, block string) string {
	content = stripBlock(content)
	if content == "" {
		return block
	}
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return content + "\n" + block
}

// stripBlock returns content with the marked block and its surrounding blank
// lines removed; content without a block is returned unchanged.
func stripBlock(content string) string {
	start := strings.Index(content, blockStart)
	end := strings.Index(content, blockEnd)
	if start == -1 || end == -1 || end < start {
		return content
	}
	end += len(blockEnd)
	if end < len(content) && content[end] == '\n' {
		end++
	}
	before := strings.TrimRight(content[:start], "\n")
	after := strings.TrimLeft(content[end:], "\n")
	switch {
	case before == "":
		return after
	case after == "":
		return before + "\n"
	default:
		return before + "\n\n" + after
	}
}

// rcFile returns the shell rc file the alias block belongs in.
func rcFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	shell := filepath.Base(os.Getenv("SHELL"))
	switch shell {
	case "zsh":
		return filepath.Join(home, ".zshrc"), nil
	case "bash":
		return filepath.Join(home, ".bashrc"), nil
	default:
		return "", fmt.Errorf("unsupported shell %q: add the aliases to your shell config yourself:\n\n%s", shell, aliasBlock(executablePath()))
	}
}

// tildify abbreviates home to ~ in path for display.
func tildify(path, home string) string {
	if path == home {
		return "~"
	}
	if strings.HasPrefix(path, home+string(filepath.Separator)) {
		return "~" + path[len(home):]
	}
	return path
}

func displayPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return tildify(path, home)
}

func executablePath() string {
	exe, err := os.Executable()
	if err != nil {
		return "paush"
	}
	return exe
}

func install() error {
	rc, err := rcFile()
	if err != nil {
		return err
	}
	content, err := os.ReadFile(rc)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	updated := upsertBlock(string(content), aliasBlock(executablePath()))
	if err := os.WriteFile(rc, []byte(updated), 0o644); err != nil {
		return err
	}
	fmt.Printf("installed %d typo aliases in %s\n", len(typos), displayPath(rc))
	fmt.Println("restart your shell (or `source` the file) to activate them")
	return nil
}

func uninstall() error {
	rc, err := rcFile()
	if err != nil {
		return err
	}
	content, err := os.ReadFile(rc)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("no aliases installed: %s does not exist\n", displayPath(rc))
			return nil
		}
		return err
	}
	stripped := stripBlock(string(content))
	if stripped == string(content) {
		fmt.Printf("no paush alias block found in %s\n", displayPath(rc))
		return nil
	}
	if err := os.WriteFile(rc, []byte(stripped), 0o644); err != nil {
		return err
	}
	fmt.Printf("removed typo aliases from %s\n", displayPath(rc))
	return nil
}

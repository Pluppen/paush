// ABOUTME: Tests for typo-alias installation into shell rc files.
// ABOUTME: Verifies alias block generation and idempotent insert/removal of the marked block.
package main

import (
	"strings"
	"testing"
)

func TestAliasBlockContainsTypos(t *testing.T) {
	block := aliasBlock("/usr/local/bin/paush")
	for _, want := range []string{
		"alias sl='/usr/local/bin/paush typo'",
		"alias dri='/usr/local/bin/paush typo'",
		"alias cta='/usr/local/bin/paush typo'",
		"alias gti='/usr/local/bin/paush typo'",
	} {
		if !strings.Contains(block, want) {
			t.Errorf("alias block missing %q", want)
		}
	}
	if !strings.HasPrefix(block, blockStart+"\n") {
		t.Error("alias block does not start with the start marker")
	}
	if !strings.HasSuffix(block, blockEnd+"\n") {
		t.Error("alias block does not end with the end marker")
	}
}

func TestUpsertBlockIntoEmptyFile(t *testing.T) {
	block := blockStart + "\nalias sl='/bin/paush' # ls\n" + blockEnd + "\n"
	if got := upsertBlock("", block); got != block {
		t.Errorf("upsert into empty content = %q, want %q", got, block)
	}
}

func TestUpsertBlockAppendsAfterExistingContent(t *testing.T) {
	block := blockStart + "\nalias sl='/bin/paush' # ls\n" + blockEnd + "\n"
	got := upsertBlock("export PATH=$PATH:/opt\n", block)
	want := "export PATH=$PATH:/opt\n\n" + block
	if got != want {
		t.Errorf("upsert after content = %q, want %q", got, want)
	}
}

func TestUpsertBlockReplacesExistingBlock(t *testing.T) {
	oldBlock := blockStart + "\nalias sl='/old/paush' # ls\n" + blockEnd + "\n"
	newBlock := blockStart + "\nalias sl='/new/paush' # ls\n" + blockEnd + "\n"
	content := "export A=1\n\n" + oldBlock + "\nalias g=git\n"
	got := upsertBlock(content, newBlock)
	if strings.Contains(got, "/old/paush") {
		t.Error("old block content survived upsert")
	}
	if !strings.Contains(got, "/new/paush") {
		t.Error("new block content missing after upsert")
	}
	if !strings.Contains(got, "export A=1") || !strings.Contains(got, "alias g=git") {
		t.Error("surrounding content lost during upsert")
	}
	if strings.Count(got, blockStart) != 1 {
		t.Errorf("expected exactly one block after upsert, found %d markers", strings.Count(got, blockStart))
	}
}

func TestStripBlockRemovesBlock(t *testing.T) {
	block := blockStart + "\nalias sl='/bin/paush' # ls\n" + blockEnd + "\n"
	content := "export A=1\n\n" + block + "\nalias g=git\n"
	got := stripBlock(content)
	if strings.Contains(got, blockStart) || strings.Contains(got, "/bin/paush") {
		t.Errorf("block survived strip: %q", got)
	}
	if !strings.Contains(got, "export A=1") || !strings.Contains(got, "alias g=git") {
		t.Errorf("surrounding content lost during strip: %q", got)
	}
}

func TestTildify(t *testing.T) {
	cases := []struct{ path, home, want string }{
		{"/Users/elias/.zshrc", "/Users/elias", "~/.zshrc"},
		{"/tmp/other/.zshrc", "/Users/elias", "/tmp/other/.zshrc"},
		{"/Users/elias", "/Users/elias", "~"},
		{"/Users/eliasson/.zshrc", "/Users/elias", "/Users/eliasson/.zshrc"},
	}
	for _, c := range cases {
		if got := tildify(c.path, c.home); got != c.want {
			t.Errorf("tildify(%q, %q) = %q, want %q", c.path, c.home, got, c.want)
		}
	}
}

func TestStripBlockWithoutBlockIsUnchanged(t *testing.T) {
	content := "export A=1\n"
	if got := stripBlock(content); got != content {
		t.Errorf("strip without block = %q, want %q", got, content)
	}
}

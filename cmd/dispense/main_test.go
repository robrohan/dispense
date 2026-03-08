package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// =========================================================================
// parseFrontMatter

func TestParseFrontMatter_Valid(t *testing.T) {
	input := []byte("---\ntitle: Hello\nauthor: rob\n---\nContent here")
	node, _, _ := parseFrontMatter(input)
	if node == nil {
		t.Fatal("expected node, got nil")
	}
	fm := node.(*FrontMatter).Data
	if fm["title"] != "Hello" {
		t.Errorf("expected title 'Hello', got %v", fm["title"])
	}
	if fm["author"] != "rob" {
		t.Errorf("expected author 'rob', got %v", fm["author"])
	}
}

func TestParseFrontMatter_NoPrefix(t *testing.T) {
	input := []byte("No front matter here")
	node, _, advance := parseFrontMatter(input)
	if node != nil {
		t.Errorf("expected nil node, got %v", node)
	}
	if advance != 0 {
		t.Errorf("expected advance 0, got %d", advance)
	}
}

func TestParseFrontMatter_DoubleNewlinePrefix(t *testing.T) {
	// "---\n\n" prefix should be treated as a horizontal rule, not front matter
	input := []byte("---\n\nSome content")
	node, _, advance := parseFrontMatter(input)
	if node != nil {
		t.Errorf("expected nil node for ---\\n\\n prefix, got %v", node)
	}
	if advance != 0 {
		t.Errorf("expected advance 0, got %d", advance)
	}
}

func TestParseFrontMatter_NoClosingDelimiter(t *testing.T) {
	input := []byte("---\ntitle: Hello\nno closing delimiter")
	node, remaining, advance := parseFrontMatter(input)
	if node != nil {
		t.Errorf("expected nil node for unclosed front matter, got %v", node)
	}
	if remaining == nil {
		t.Error("expected remaining data to be returned for unclosed front matter")
	}
	if advance != 0 {
		t.Errorf("expected advance 0, got %d", advance)
	}
}

func TestParseFrontMatter_BadYAML(t *testing.T) {
	// A tab at the start of a YAML line is invalid
	input := []byte("---\n\tkey: bad yaml\n---\nContent")
	node, remaining, advance := parseFrontMatter(input)
	if node != nil {
		t.Errorf("expected nil node for bad YAML, got %v", node)
	}
	if remaining == nil {
		t.Error("expected original data returned when YAML is invalid")
	}
	if advance != 0 {
		t.Errorf("expected advance 0, got %d", advance)
	}
}

func TestParseFrontMatter_Empty(t *testing.T) {
	// Empty front matter block is valid - just produces an empty map
	input := []byte("---\n---\nContent")
	node, _, _ := parseFrontMatter(input)
	if node == nil {
		t.Fatal("expected node for empty front matter, got nil")
	}
	fm := node.(*FrontMatter).Data
	if len(fm) != 0 {
		t.Errorf("expected empty map, got %v", fm)
	}
}

func TestParseFrontMatter_AdvanceSkipsPastDelimiters(t *testing.T) {
	// Verify the returned advance consumes the full front matter block
	// so the markdown parser continues from the right position
	body := "Content after front matter"
	input := []byte("---\ntitle: Test\n---\n" + body)
	_, _, advance := parseFrontMatter(input)
	// advance stops after the closing "---"; the trailing newline is not consumed
	remaining := strings.TrimPrefix(string(input[advance:]), "\n")
	if !strings.HasPrefix(remaining, body) {
		t.Errorf("expected remaining input to start with %q, got %q", body, remaining)
	}
}

// =========================================================================
// mdToHTML

func TestMdToHTML_DefaultsToPostTemplate(t *testing.T) {
	input := []byte("# Hello\n\nSome content")
	_, tmpl, _ := mdToHTML(input)
	if tmpl != "post" {
		t.Errorf("expected default template 'post', got %q", tmpl)
	}
}

func TestMdToHTML_NoFrontMatterProducesEmptyMap(t *testing.T) {
	input := []byte("# Hello\n\nSome content")
	_, _, fm := mdToHTML(input)
	if len(fm) != 0 {
		t.Errorf("expected empty fm for content without front matter, got %v", fm)
	}
}

func TestMdToHTML_FrontMatterFieldsPassedThrough(t *testing.T) {
	input := []byte("---\ntitle: My Title\nauthor: rob\n---\n\n# Hello\n")
	_, _, fm := mdToHTML(input)
	if fm["title"] != "My Title" {
		t.Errorf("expected title 'My Title', got %v", fm["title"])
	}
	if fm["author"] != "rob" {
		t.Errorf("expected author 'rob', got %v", fm["author"])
	}
}

func TestMdToHTML_CustomTemplate(t *testing.T) {
	input := []byte("---\ntitle: My Page\ntemplate: page\n---\n\nContent\n")
	_, tmpl, _ := mdToHTML(input)
	if tmpl != "page" {
		t.Errorf("expected template 'page', got %q", tmpl)
	}
}

func TestMdToHTML_NoTemplateFallsBackToPost(t *testing.T) {
	input := []byte("---\ntitle: My Page\n---\n\nContent\n")
	_, tmpl, _ := mdToHTML(input)
	if tmpl != "post" {
		t.Errorf("expected fallback template 'post', got %q", tmpl)
	}
}

func TestMdToHTML_RendersBold(t *testing.T) {
	input := []byte("**bold text**")
	out, _, _ := mdToHTML(input)
	if !strings.Contains(string(out), "<strong>bold text</strong>") {
		t.Errorf("expected <strong>bold text</strong>, got: %s", string(out))
	}
}

func TestMdToHTML_RendersItalic(t *testing.T) {
	input := []byte("*italic text*")
	out, _, _ := mdToHTML(input)
	if !strings.Contains(string(out), "<em>italic text</em>") {
		t.Errorf("expected <em>italic text</em>, got: %s", string(out))
	}
}

func TestMdToHTML_RendersHeading(t *testing.T) {
	input := []byte("# Top Level Heading")
	out, _, _ := mdToHTML(input)
	if !strings.Contains(string(out), "<h1") {
		t.Errorf("expected <h1 in output, got: %s", string(out))
	}
	if !strings.Contains(string(out), "Top Level Heading") {
		t.Errorf("expected heading text in output, got: %s", string(out))
	}
}

func TestMdToHTML_RendersCodeBlock(t *testing.T) {
	input := []byte("```\nsome code\n```\n")
	out, _, _ := mdToHTML(input)
	if !strings.Contains(string(out), "<code>") {
		t.Errorf("expected <code> in output, got: %s", string(out))
	}
}

func TestMdToHTML_RendersLink(t *testing.T) {
	input := []byte("[click here](https://example.com)")
	out, _, _ := mdToHTML(input)
	if !strings.Contains(string(out), `href="https://example.com"`) {
		t.Errorf("expected href in output, got: %s", string(out))
	}
}

func TestMdToHTML_EmptyInput(t *testing.T) {
	input := []byte("")
	out, tmpl, fm := mdToHTML(input)
	if tmpl != "post" {
		t.Errorf("expected 'post' for empty input, got %q", tmpl)
	}
	if len(fm) != 0 {
		t.Errorf("expected empty fm for empty input, got %v", fm)
	}
	_ = out // no crash is the requirement
}

// =========================================================================
// FilePathWalkDir

func TestFilePathWalkDir_FindsMdFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.md"), []byte("# A"), 0644)
	os.WriteFile(filepath.Join(dir, "b.md"), []byte("# B"), 0644)
	os.WriteFile(filepath.Join(dir, "c.txt"), []byte("not markdown"), 0644)

	files, err := FilePathWalkDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 .md files, got %d: %v", len(files), files)
	}
	for _, f := range files {
		if filepath.Ext(f) != ".md" {
			t.Errorf("expected only .md files, got %q", f)
		}
	}
}

func TestFilePathWalkDir_IgnoresNonMdFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "page.html"), []byte("<html>"), 0644)
	os.WriteFile(filepath.Join(dir, "style.css"), []byte("body{}"), 0644)
	os.WriteFile(filepath.Join(dir, "data.json"), []byte("{}"), 0644)

	files, err := FilePathWalkDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d: %v", len(files), files)
	}
}

func TestFilePathWalkDir_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	files, err := FilePathWalkDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files in empty dir, got %d", len(files))
	}
}

func TestFilePathWalkDir_NonExistentDir(t *testing.T) {
	_, err := FilePathWalkDir("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Error("expected error for non-existent directory, got nil")
	}
}

func TestFilePathWalkDir_IncludesSubdirFiles(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	os.Mkdir(sub, 0755)
	os.WriteFile(filepath.Join(dir, "top.md"), []byte("# Top"), 0644)
	os.WriteFile(filepath.Join(sub, "nested.md"), []byte("# Nested"), 0644)

	files, err := FilePathWalkDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Current behaviour: Walk is recursive, so both files are found
	if len(files) != 2 {
		t.Errorf("expected 2 files (top-level + subdirectory), got %d: %v", len(files), files)
	}
}

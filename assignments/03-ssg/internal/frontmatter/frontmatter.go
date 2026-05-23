// Package frontmatter parses YAML front matter from a Markdown document.
package frontmatter

// Split separates `---`-delimited YAML front matter from the body.
// Returns (yamlBytes, bodyBytes). yamlBytes is nil when no front matter present.
func Split(src []byte) (yaml []byte, body []byte) {
	// TODO
	return nil, src
}

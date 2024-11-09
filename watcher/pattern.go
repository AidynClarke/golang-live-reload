package watcher

import (
	"log"

	"github.com/gobwas/glob"
)

type GlobMatcher struct {
	globs []glob.Glob
}
func NewGlobMatcher(patterns []string) *GlobMatcher {
	globs := make([]glob.Glob, len(patterns))

	for i, pattern := range patterns {
		g, err := glob.Compile(pattern)

		if err != nil {
			log.Fatal(err)
		}

		globs[i] = g
	}

	return &GlobMatcher{
		globs: globs,
	}
}

func (m *GlobMatcher) Match(path string) bool {
	for _, g := range m.globs {
		if g.Match(path) {
			return true
		}
	}
	return false
}


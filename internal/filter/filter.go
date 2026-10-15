// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package filter

import (
	"strings"

	"github.com/gobwas/glob"
)

// Names returns the names matching the glob patterns, with negative patterns support.
// Without positive patterns all names match.
func Names(names, patterns []string) (_ []string, err error) {
	var (
		exclude = make([]glob.Glob, 0, len(patterns))
		include = make([]glob.Glob, 0, len(patterns))
	)

	// compile patterns, with negative support
	for _, pattern := range patterns {
		var index = 0
		if strings.HasPrefix(pattern, "!") {
			index = 1
		}

		var g glob.Glob
		if g, err = glob.Compile(pattern[index:]); err != nil {
			return nil, err
		}

		if index == 0 {
			include = append(include, g)
			continue
		}

		exclude = append(exclude, g)
	}

	// all included by default
	if len(include) == 0 {
		include = append(include, glob.MustCompile("*"))
	}

	var filtered = make([]string, 0, len(names))
	for _, name := range names {
		if match(name, include, exclude) {
			filtered = append(filtered, name)
		}
	}

	return filtered, nil
}

func match(name string, include, exclude []glob.Glob) bool {
	for _, pattern := range exclude {
		if pattern.Match(name) {
			return false
		}
	}

	for _, pattern := range include {
		if pattern.Match(name) {
			return true
		}
	}

	return false
}

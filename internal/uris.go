package main

import (
	"bytes"
	"strings"
)

// URLJoin joins the given paths so that there is only ever one '/' character between the paths
func URLJoin(paths ...string) string {
	/* #nosec */
	var buffer bytes.Buffer
	last := len(paths) - 1
	for i, path := range paths {
		p := path
		if i > 0 {
			buffer.WriteString("/")
			p = strings.TrimPrefix(p, "/")
		}
		if i < last {
			p = strings.TrimSuffix(p, "/")
		}
		buffer.WriteString(p)
	}
	return buffer.String()
}

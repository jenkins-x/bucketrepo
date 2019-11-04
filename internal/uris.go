package main

import (
	"bytes"
	"strings"

	"github.com/sirupsen/logrus"
)

// URLJoin joins the given paths so that there is only ever one '/' character between the paths
func URLJoin(paths ...string) string {
	/* #nosec */
	var buffer bytes.Buffer
	last := len(paths) - 1
	for i, path := range paths {
		p := path
		if i > 0 {
			_, err := buffer.WriteString("/")
			if err != nil {
				logrus.Errorf("failed to write to memory buffer: %s", err.Error())
			}
			p = strings.TrimPrefix(p, "/")
		}
		if i < last {
			p = strings.TrimSuffix(p, "/")
		}
		_, err := buffer.WriteString(p)
		if err != nil {
			logrus.Errorf("failed to write to memory buffer: %s", err.Error())
		}
	}
	return buffer.String()
}

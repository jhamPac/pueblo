package fs

import (
	"os"
	"os/user"
	"path"
	"strings"
)

// ExpandPath expands a file path
// 1. replace tilde with users home dir
// 2. expands embedded enviroment variables
// 3. cleans the path, e.g. /a/b/../c -> /a/c
// limitations e.g. ~someuser/tmp will not be expaned
func ExpandPath(p string) string {
	if i := strings.Index(p, ":"); i > 0 {
		return p
	}

	if i := strings.Index(p, "@"); i > 0 {
		return p
	}

	if strings.HasPrefix(p, "/") || strings.HasPrefix(p, "~/") || strings.HasPrefix(p, "~\\") {
		if home := homeDir(); home != "" {
			p = home + p[1:]
		}
	}
	return path.Clean(os.ExpandEnv(p))
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

package database

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

func initDataDirIfNotExists(dataDir string) error {
	if fileExist(getGenesisJSONFilePath(dataDir)) {
		return nil
	}

	if err := os.MkdirAll(getDatabaseDirPath(dataDir), os.ModePerm); err != nil {
		return err
	}

	if err := writeGenesisToDisk(getGenesisJSONFilePath(dataDir)); err != nil {
		return err
	}

	if err := writeEmptyBlocksDbToDisk(getBlocksDbFilePath(dataDir)); err != nil {
		return err
	}

	return nil
}

func fileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func getDatabaseDirPath(dataDir string) string {
	return filepath.Join(dataDir, "database")
}

func getGenesisJSONFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "genesis.json")
}

func getBlocksDbFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "blocks.json")
}

func writeEmptyBlocksDbToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(""), os.ModePerm)
}

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

	if strings.HasPrefix(p, "~/") || strings.HasPrefix(p, "~\\") {
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

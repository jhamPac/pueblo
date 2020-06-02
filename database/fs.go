package database

import (
	"os"
	"path/filepath"
	"io/ioutil"
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
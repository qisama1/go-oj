package utils

import "os"

func CodeSave(code []byte) (string, error) {
	dirName := "code/" + GetUUID()
	path := dirName + "main.go"
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	f.Write(code)
	defer f.Close()
	return path, err
}

package share

import (
	"log"
	"os"
)

func SaveTo(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()
	_, err = file.Write(content)
	return err
}

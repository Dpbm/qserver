package files

import (
	"bufio"
	"errors"
	"os"

	logger "github.com/Dpbm/shared/log"
)

type Qasm struct {
	Path     string
	Filename string
	File     *os.File
}

func (file *Qasm) RemoveFile() {
	if file.File == nil {
		logger.LogError(errors.New("file does not exists"))
		return
	}

	err := os.Remove(file.Path)

	if err != nil {
		logger.LogError(err)
		return
	}

	file.File = nil
	file.Filename = ""
	file.Path = ""
}

func (file *Qasm) CreateFile() error {
	qasmFile, err := os.Create(file.Path)
	if err != nil {
		return err
	}

	file.File = qasmFile
	return nil
}

func (file *Qasm) Close() {
	if file.File == nil {
		return
	}

	file.File.Close()
}

func (file *Qasm) AddChunckToFile(chunck string) error {
	if file.File == nil {
		return errors.New("file not defined")
	}

	writer := bufio.NewWriter(file.File)
	defer writer.Flush()

	qasmWritting, err := writer.WriteString(chunck)

	if err != nil {
		return err
	}

	if qasmWritting <= 0 {
		return errors.New("nothing was written")
	}

	return nil
}

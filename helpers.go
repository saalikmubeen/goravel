package goravel

import "os"

// CreateDirIfNotExists creates a new directory if it does not exist
func (g *Goravel) CreateDirIfNotExists(path string) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0755) // default permissions
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateFileIfNotExists creates a new file at path if it does not exist
func (g *Goravel) CreateFileIfNotExists(path string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}

		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}
	return nil
}

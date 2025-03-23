package rds

import "os"

func (rds *RDS) Create(username string, subfolder string) error {
	path, err := rds.rdsPath(username, subfolder)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return os.Chown(path, 1000, 100)
}

func (rds *RDS) List() []string {
	files, err := os.ReadDir(rds.BasePath)
	if err != nil {
		return nil
	}

	var names []string
	for _, file := range files {
		if file.IsDir() {
			names = append(names, file.Name())
		}
	}
	return names
}

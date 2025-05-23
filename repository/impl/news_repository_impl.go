package impl

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/MrWhok/IMK-FP-BACKEND/model"
)

type FileNewsRepo struct {
	FilePath string
}

func NewFileNewsRepo(path string) *FileNewsRepo {
	return &FileNewsRepo{FilePath: path}
}
func (r *FileNewsRepo) Save(data model.CachedData) error {
	dir := filepath.Dir(r.FilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	f, err := os.Create(r.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(data)
}

func (r *FileNewsRepo) Load() (model.CachedData, error) {
	var data model.CachedData
	f, err := os.Open(r.FilePath)
	if err != nil {
		return data, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&data)
	return data, err
}

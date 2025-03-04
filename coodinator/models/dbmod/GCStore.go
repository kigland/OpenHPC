package dbmod

import (
	"github.com/kigland/HPC-Scheduler/lib/store"
	"gorm.io/gorm"
)

type GCStore struct {
	db *gorm.DB
}

func NewGCStore(db *gorm.DB) *GCStore {
	return &GCStore{db: db}
}

func (s *GCStore) Get(key string) (string, error) {
	var gc GC
	if err := s.db.Where("id = ?", key).First(&gc).Error; err != nil {
		return "", err
	}
	return gc.Value, nil
}

func (s *GCStore) Set(key string, value string) error {
	return s.db.Create(&GC{ID: key, Value: value}).Error
}

func (s *GCStore) Delete(key string) error {
	return s.db.Delete(&GC{}, "id = ?", key).Error
}

var _ store.IKVStore = &GCStore{}

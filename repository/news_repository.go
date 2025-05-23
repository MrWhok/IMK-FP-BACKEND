package repository

import "github.com/MrWhok/IMK-FP-BACKEND/model"

type NewsRepository interface {
	Save(data model.CachedData) error
	Load() (model.CachedData, error)
}

package service

import "github.com/MrWhok/IMK-FP-BACKEND/model"

type NewsService interface {
	FetchAndUpdate()
	GetNews() model.CachedData
}

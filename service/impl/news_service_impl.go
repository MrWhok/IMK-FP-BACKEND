package impl

import (
	"log"
	"sync"
	"time"

	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/MrWhok/IMK-FP-BACKEND/repository"
	"github.com/MrWhok/IMK-FP-BACKEND/scraper"
	"github.com/MrWhok/IMK-FP-BACKEND/service"
)

type NewsServiceImpl struct {
	repo  repository.NewsRepository
	cache model.CachedData
	mu    sync.Mutex
}

func NewNewsServiceImpl(repo repository.NewsRepository) service.NewsService {
	data, _ := repo.Load()
	return &NewsServiceImpl{repo: repo, cache: data}
}

func (s *NewsServiceImpl) FetchAndUpdate() {
	news, err := scraper.ScrapNews()
	if err != nil {
		log.Println("Gagal scraping:", err)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache = model.CachedData{
		LastUpdated: time.Now(),
		Data:        news,
	}
	s.repo.Save(s.cache)
}

func (s *NewsServiceImpl) GetNews() model.CachedData {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cache
}

package services

import (
	"song-lib/internal/database/postgres"
	"song-lib/internal/models"
)

type ServiceSonger interface {
	GetSongs(group, name string, page, limit int) ([]models.Song, error)
	AddSong(song *models.Song) (int64, error)
	DeleteSong(id int64) (int64, error)
	UpdateSong(song *models.Song) (int64, error)
	GetSongText(id int64) (*models.Song, error)
}

type Service struct {
	db postgres.DBSonger
}

func New(db postgres.DBSonger) ServiceSonger {
	return &Service{db: db}
}

func (s *Service) GetSongs(group, name string, page, limit int) ([]models.Song, error) {
	return s.db.GetSongs(group, name, page, limit)
}

func (s *Service) AddSong(song *models.Song) (int64, error) {
	return s.db.AddSong(song)
}

func (s *Service) DeleteSong(id int64) (int64, error) {
	return s.db.DeleteSong(id)
}

func (s *Service) UpdateSong(song *models.Song) (int64, error) {
	return s.db.UpdateSong(song)
}

func (s *Service) GetSongText(id int64) (*models.Song, error) {
	return s.db.GetSongText(id)
}

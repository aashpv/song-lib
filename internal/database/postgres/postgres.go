package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq" // init postgres driver
	"github.com/pressly/goose/v3"
	"song-lib/internal/config"
	"song-lib/internal/models"
)

type DBSonger interface {
	GetSongs(group, name string, page, limit int) ([]models.Song, error)
	AddSong(song *models.Song) (int64, error)
	DeleteSong(id int64) (int64, error)
	UpdateSong(song *models.Song) (int64, error)
	GetSongText(id int64) (*models.Song, error)
}

type Database struct {
	Db *sql.DB
}

func CreateDatabaseIfNotExists(cfg config.Database) error {
	const op = "internal.database.postgres.CreateDatabaseIfNotExists"

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=disable dbname=postgres",
		cfg.Host, cfg.User, cfg.Password, cfg.Port)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("%s: open: %w", op, err)
	}
	defer db.Close()

	// Проверка, существует ли база данных
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT datname FROM pg_database WHERE datname = $1)", cfg.Name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: query row: %w", op, err)
	}

	// Если базы данных нет, создаем её
	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.Name))
		if err != nil {
			return fmt.Errorf("%s: exec: %w", op, err)
		}
		fmt.Printf("Database %s created\n", cfg.Name)
	} else {
		fmt.Printf("Database %s already exists\n", cfg.Name)
	}

	return nil
}

func New(cfg config.Database) (DBSonger, error) {
	const op = "internal.database.postgres.New"
	const sub = "Migrate"

	if err := CreateDatabaseIfNotExists(cfg); err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: open: %w", op, err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("%s.%s: set dialect: %w", op, sub, err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return nil, fmt.Errorf("%s.%s: up: %w", op, sub, err)
	}

	return &Database{Db: db}, nil
}

func (d *Database) GetSongs(group, name string, page, limit int) ([]models.Song, error) {
	const op = "internal.database.postgres.GetSongs"
	query := "SELECT id, group_name, name, release_date, text, link FROM songs WHERE 1=1"

	// Filtration
	var args []interface{}
	if group != "" {
		query += " AND group_name = $1"
		args = append(args, group)
	}
	if name != "" {
		query += " AND name = $2"
		args = append(args, name)
	}

	// Pagination
	offset := (page - 1) * limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := d.Db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query %w", op, err)
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		err = rows.Scan(&song.ID, &song.Group, &song.Name, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			return nil, fmt.Errorf("%s: row scan: %w", op, err)
		}
		songs = append(songs, song)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: err %w", op, err)
	}
	return songs, nil
}

func (d *Database) AddSong(song *models.Song) (int64, error) {
	const op = "internal.database.postgres.AddSong"
	query := "INSERT INTO songs (group_name, name, release_date, text, link) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var songId int64

	err := d.Db.QueryRow(query, song.Group, song.Name, song.ReleaseDate, song.Text, song.Link).Scan(&songId)
	if err != nil {
		return 0, fmt.Errorf("%s: query row: %w", op, err)
	}

	return songId, nil
}

func (d *Database) DeleteSong(id int64) (int64, error) {
	const op = "internal.database.postgres.DeleteSong"
	query := "DELETE FROM songs WHERE id = $1"

	result, err := d.Db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("%s: exec %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s: rows affected %w", op, err)
	}

	return rowsAffected, nil
}

func (d *Database) UpdateSong(song *models.Song) (int64, error) {
	const op = "internal.database.postgres.UpdateSong"
	query := "UPDATE songs SET group_name = $1, name = $2, release_date = $3, text = $4, link = $5 WHERE id = $6"

	result, err := d.Db.Exec(query,
		song.Group,
		song.Name,
		song.ReleaseDate,
		song.Text,
		song.Link,
		song.ID,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: exec %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s: rows affected %w", op, err)
	}

	return rowsAffected, nil
}

func (d *Database) GetSongText(id int64) (*models.Song, error) {
	const op = "internal.database.postgres.GetSongText"
	query := "SELECT id, group_name, name, release_date, text, link FROM songs WHERE id = $1"
	var song models.Song

	err := d.Db.QueryRow(query, id).Scan(&song.ID, &song.Group, &song.Name, &song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: no song found", op)
		}
		return nil, fmt.Errorf("%s: query row scan %w", op, err)
	}
	return &song, nil
}

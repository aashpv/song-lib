package postgres_test

import (
	"database/sql"
	"log"
	"os"
	"song-lib/internal/database/postgres"
	"song-lib/internal/models"
	"testing"

	_ "github.com/lib/pq"
)

// TestMain - инициализация перед всеми тестами
func TestMain(m *testing.M) {
	dsn := "host=localhost user=postgres password=postgres dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	// Создаем таблицы для тестов
	err = createTestTables(db)
	if err != nil {
		log.Fatalf("failed to create test tables: %v", err)
	}

	// Запускаем тесты
	code := m.Run()

	// Удаляем таблицы после тестов
	err = dropTestTables(db)
	if err != nil {
		log.Fatalf("failed to drop test tables: %v", err)
	}

	// Закрываем базу данных
	db.Close()

	// Завершаем выполнение тестов
	os.Exit(code)
}

// createTestTables - создание тестовых таблиц
func createTestTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS songs (
			id SERIAL PRIMARY KEY,
			group_name VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			release_date DATE,
			text TEXT,
			link VARCHAR(255)
		);
	`)
	return err
}

// dropTestTables - удаление тестовых таблиц
func dropTestTables(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS songs;")
	return err
}

// TestAddSong - интеграционный тест для метода AddSong
func TestAddSong(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	repo := postgres.Database{Db: db}

	song := &models.Song{
		Group:       "Muse",
		Name:        "Supermassive Black Hole",
		ReleaseDate: "2006-07-16",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	// Тестируем добавление песни
	id, err := repo.AddSong(song)
	if err != nil {
		t.Fatalf("failed to add song: %v", err)
	}

	// Проверяем, что песня была добавлена
	var addedSong models.Song
	err = db.QueryRow("SELECT group_name, name FROM songs WHERE id = $1", id).Scan(&addedSong.Group, &addedSong.Name)
	if err != nil {
		t.Fatalf("failed to retrieve added song: %v", err)
	}

	if addedSong.Group != song.Group || addedSong.Name != song.Name {
		t.Errorf("expected song group: %s and name: %s, got group: %s and name: %s", song.Group, song.Name, addedSong.Group, addedSong.Name)
	}

	// Удаляем данные после теста
	_, err = db.Exec("DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		t.Fatalf("failed to delete added song: %v", err)
	}

	// Закрываем базу данных
	db.Close()
}

// TestGetSong - интеграционный тест для метода GetSongs
func TestGetSongs(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	repo := postgres.Database{Db: db}

	songs := []*models.Song{
		{Group: "Muse", Name: "Supermassive Black Hole", ReleaseDate: "2006-07-16"},
		{Group: "Radiohead", Name: "Creep", ReleaseDate: "1993-09-21"},
	}
	for _, song := range songs {
		_, err = repo.AddSong(song)
		if err != nil {
			t.Fatalf("failed to add song: %v", err)
		}
	}

	// Тестируем получение песен
	result, err := repo.GetSongs("Muse", "", 1, 10)
	if err != nil {
		t.Fatalf("failed to get songs: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 song, got %d", len(result))
	}

	// Чистим данные после теста
	_, _ = db.Exec("DELETE FROM songs WHERE group_name = 'Muse' OR group_name = 'Radiohead'")

	// Закрываем базу данных
	db.Close()
}

// TestDeleteSong - интеграционный тест для метода DeleteSong
func TestDeleteSong(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	repo := postgres.Database{Db: db}

	song := &models.Song{
		Group:       "Muse",
		Name:        "Supermassive Black Hole",
		ReleaseDate: "2006-07-16",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	id, err := repo.AddSong(song)
	if err != nil {
		t.Fatalf("failed to add song: %v", err)
	}

	// Тестируем удаление песни
	err = repo.DeleteSong(id)
	if err != nil {
		t.Fatalf("failed to delete song: %v", err)
	}

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM songs WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		t.Fatalf("failed to check if song exists: %v", err)
	}
	if exists {
		t.Errorf("expected song to be deleted, but it still exists")
	}

	// Закрываем базу данных
	db.Close()
}

// TestUpdateSong - интеграционный тест для метода UpdateSong
func TestUpdateSong(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	repo := postgres.Database{Db: db}

	song := &models.Song{
		Group:       "Muse",
		Name:        "Supermassive Black Hole",
		ReleaseDate: "2006-07-16",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	id, err := repo.AddSong(song)
	if err != nil {
		t.Fatalf("failed to add song: %v", err)
	}

	// Тестируем изменение песни
	song.ID = id
	song.Name = "Updated Song Name"
	err = repo.UpdateSong(song)
	if err != nil {
		t.Fatalf("failed to update song: %v", err)
	}

	// Проверяем, что песня была обновлена
	var updatedSong models.Song
	err = db.QueryRow("SELECT name FROM songs WHERE id = $1", song.ID).Scan(&updatedSong.Name)
	if err != nil {
		t.Fatalf("failed to retrieve updated song: %v", err)
	}

	if updatedSong.Name != "Updated Song Name" {
		t.Errorf("expected updated song name to be 'Updated Song Name', but got '%s'", updatedSong.Name)
	}

	// Удаляем данные после теста
	_, err = db.Exec("DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		t.Fatalf("failed to delete song after test: %v", err)
	}

	// Закрываем базу данных
	db.Close()
}

// TestGetSongText - интеграционный тест для метода GetSongText
func TestGetSongText(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	repo := postgres.Database{Db: db}

	song := &models.Song{
		Group:       "Muse",
		Name:        "Supermassive Black Hole",
		ReleaseDate: "2006-07-16",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	id, err := repo.AddSong(song)
	if err != nil {
		t.Fatalf("failed to add song: %v", err)
	}

	// Тестируем получение текста песни
	retrievedSong, err := repo.GetSongText(id)
	if err != nil {
		t.Fatalf("failed to get song text: %v", err)
	}

	if retrievedSong.Text != song.Text {
		t.Errorf("expected song text: %s, but got: %s", song.Text, retrievedSong.Text)
	}

	// Удаляем данные после теста
	_, err = db.Exec("DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		t.Fatalf("failed to delete song after test: %v", err)
	}

	// Закрываем базу данных
	db.Close()
}

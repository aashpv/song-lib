package text

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"song-lib/internal/lib/resp"
	"song-lib/internal/models"
	"strconv"
	"strings"
)

type Response struct {
	Group  string   `json:"group"`
	Song   string   `json:"song"`
	Verses []string `json:"verses"`
	Total  int      `json:"total"`
}

type SongTextGetter interface {
	GetSongText(id int64) (*models.Song, error)
}

// New gets the lyrics of the song paginated
// @Summary Get song lyrics with pagination
// @Description Get song text divided by verses with pagination support
// @Tags Songs
// @Param id path int true "Song ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of verses per page" default(3)
// @Produce  json
// @Success 200 {object} text.Response
// @Failure 400 {object} resp.Response "Invalid request"
// @Failure 404 {object} resp.Response "Song not found"
// @Failure 500 {object} resp.Response "Failed to get song text"
// @Router /songs/{id}/text [get]
func New(log *slog.Logger, getter SongTextGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.transport.handlers.text.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			log.Error("invalid song id", "error", err)
			render.JSON(w, r, resp.Error("invalid song id"))
			return
		}

		song, err := getter.GetSongText(id)
		if err != nil {
			log.Error("failed to get song text", "error", err)
			render.JSON(w, r, resp.Error("failed to get song text"))
			return
		}

		verses := strings.Split(song.Text, "\n\n")

		page, limit := parsePagination(r)

		start := (page - 1) * limit
		if start >= len(verses) {
			render.JSON(w, r, resp.Error("no verses found for this page"))
			return
		}

		end := start + limit
		if end > len(verses) {
			end = len(verses)
		}

		log.Info("verses retrieved successfully", slog.Int64("song_id", id), slog.Int("page", page))

		render.JSON(w, r, Response{
			Group:  song.Group,
			Song:   song.Name,
			Verses: verses[start:end],
			Total:  len(verses),
		})
	}
}

func parsePagination(r *http.Request) (int, int) {
	page := 1
	limit := 3 // По умолчанию 3 куплета на страницу

	if p := r.URL.Query().Get("page"); p != "" {
		if pInt, err := strconv.Atoi(p); err == nil && pInt > 0 {
			page = pInt
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if lInt, err := strconv.Atoi(l); err == nil && lInt > 0 {
			limit = lInt
		}
	}

	return page, limit
}

package get

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"song-lib/internal/lib/resp"
	"song-lib/internal/models"
	"strconv"
)

type SongGetter interface {
	GetSongs(group, name string, page, limit int) ([]models.Song, error)
}

// New gets songs with pagination
// @Summary Get song with pagination
// @Description Get songs with pagination support
// @Tags Songs
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of songs per page" default(10)
// @Produce  json
// @Success 200 {object} []models.Song
// @Failure 400 {object} resp.Response "Invalid request"
// @Failure 500 {object} resp.Response "Failed to get songs"
// @Router /songs [get]
func New(log *slog.Logger, getter SongGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.transport.handlers.get.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		group := r.URL.Query().Get("group")
		name := r.URL.Query().Get("name")
		if group == "" || name == "" {
			log.Error("missing group or name in query parameters")
			render.JSON(w, r, resp.Error("group and name parameters are required"))
			return
		}
		log.Info("request parameters decoded", slog.String("group", group), slog.String("name", name))

		page, limit := parsePagination(r)

		songs, err := getter.GetSongs(group, name, page, limit)
		if err != nil {
			log.Error("failed to get songs", "error", err)
			render.JSON(w, r, resp.Error("failed to get songs"))
			return
		}

		if len(songs) == 0 {
			log.Info("no songs found for request", slog.String("group", group), slog.String("name", name))
			render.JSON(w, r, resp.Error("no songs found"))
			return
		}

		log.Info("songs retrieved successfully", slog.Int("count", len(songs)))

		render.JSON(w, r, songs)
	}
}

func parsePagination(r *http.Request) (int, int) {
	page := 1
	limit := 10

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

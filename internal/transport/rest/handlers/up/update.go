package up

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"song-lib/internal/lib/resp"
	"song-lib/internal/models"
	"strconv"
)

type Request struct {
	Group       string `json:"group" validate:"required"`
	Name        string `json:"song" validate:"required"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type Response struct {
	resp.Response
	Msg string `json:"msg"`
}

type SongUpdater interface {
	UpdateSong(song *models.Song) (int64, error)
}

// New changes the song in the library
// @Summary Update a song
// @Description Update a song in the library by its ID
// @Tags Songs
// @Param id path int true "Song ID"
// @Param song body up.Request true "Updated song details"
// @Produce  json
// @Success 200 {object} resp.Response
// @Failure 404 {object} resp.Response "Song not found"
// @Failure 500 {object} resp.Response "Failed to update song"
// @Router /songs/{id} [put]
func New(log *slog.Logger, updater SongUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.transport.handlers.update.New"

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

		var req Request
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", "error", err)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", "error", err)
			render.JSON(w, r, resp.Error("invalid request: missing or invalid group and song"))
			return
		}

		updatedSong := &models.Song{
			ID:          id,
			Group:       req.Group,
			Name:        req.Name,
			ReleaseDate: req.ReleaseDate,
			Text:        req.Text,
			Link:        req.Link,
		}

		rowsAffected, err := updater.UpdateSong(updatedSong)
		if err != nil {
			log.Error("failed to update song", "error", err)
			render.JSON(w, r, resp.Error("failed to update song"))
			return
		}

		if rowsAffected == 0 {
			log.Error("song not found", slog.Int64("song_id", id))
			render.JSON(w, r, resp.Error("song not found"))
			return
		}

		log.Info("song updated successfully", slog.Int64("song_id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Msg:      "success",
		})
	}
}

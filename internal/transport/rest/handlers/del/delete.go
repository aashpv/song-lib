package del

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"song-lib/internal/lib/resp"
	"strconv"
)

type Response struct {
	resp.Response
	Msg string `json:"msg"`
}

type SongDeleter interface {
	DeleteSong(id int64) (int64, error)
}

// New deletes a song from the library
// @Summary Delete a song
// @Description Delete a song from the library by its ID
// @Tags Songs
// @Param id path int true "Song ID"
// @Produce  json
// @Success 200 {object} del.Response
// @Failure 404 {object} resp.Response "Song not found"
// @Failure 500 {object} resp.Response "Failed to delete song"
// @Router /songs/{id} [delete]
func New(log *slog.Logger, deleter SongDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.transport.handlers.del.New"

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

		rowsAffected, err := deleter.DeleteSong(id)
		if err != nil {
			log.Error("failed to delete song", "error", err)
			render.JSON(w, r, resp.Error("failed to delete song"))
			return
		}

		if rowsAffected == 0 {
			log.Error("song not found", slog.Int64("song_id", id))
			render.JSON(w, r, resp.Error("song not found"))
			return
		}

		log.Info("song deleted successfully", slog.Int64("song_id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Msg:      "success",
		})
	}
}

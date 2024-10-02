package add

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"song-lib/internal/lib/resp"
	"song-lib/internal/models"
	"strings"
)

type Request struct {
	Group string `json:"group" validate:"required"`
	Song  string `json:"song" validate:"required"`
}

type Response struct {
	resp.Response
	Msg string `json:"msg,omitempty"`
}

type SongDetails struct {
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongAdder interface {
	AddSong(song *models.Song) (int64, error)
}

// New adds a new song to the library
// @Summary Add a new song
// @Description Add a new song to the library and fetch additional details from an external API
// @Tags Songs
// @Accept  json
// @Produce  json
// @Param song body add.Request true "Song details"
// @Success 200 {object} add.Response "Song added successfully"
// @Failure 400 {object} resp.Response "Invalid request"
// @Failure 500 {object} resp.Response "Failed to add song"
// @Router /songs [post]
func New(log *slog.Logger, adder SongAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.transport.handlers.add.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", "error", err)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", "error", err)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		group := strings.Replace(req.Group, " ", "+", -1)
		song := strings.Replace(req.Song, " ", "+", -1)
		externalApiUrl := fmt.Sprintf("http://localhost:8081/info?group=%s&song=%s", group, song)
		fmt.Println(externalApiUrl)

		response, err := http.Get(externalApiUrl)
		if err != nil {
			log.Error("failed to send request", "error", err)
			render.JSON(w, r, resp.Error("failed to send request"))
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			log.Error("failed to send request", "status_code", response.StatusCode)
			render.JSON(w, r, resp.Error("failed to send request"))
			return
		}

		var songDetails SongDetails
		err = render.DecodeJSON(response.Body, &songDetails)
		if err != nil {
			log.Error("failed to decode response", "error", err)
			render.JSON(w, r, resp.Error("failed to decode response"))
			return
		}

		newSong := &models.Song{
			Group:       req.Group,
			Name:        req.Song,
			ReleaseDate: songDetails.ReleaseDate,
			Text:        songDetails.Text,
			Link:        songDetails.Link,
		}

		id, err := adder.AddSong(newSong)
		if err != nil {
			log.Error("failed to add song", "error", err)
			render.JSON(w, r, resp.Error("failed to add song"))
			return
		}

		log.Info("song added successfully", slog.Int64("song_id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Msg:      "success",
		})
	}
}

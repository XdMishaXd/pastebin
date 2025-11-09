package save

import (
	"log/slog"
	resp "main_service/internal/lib/api/response"
	sl "main_service/internal/lib/logger"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Text string `json:"text" validate:"required"`
	TTL  int    `json:"ttl,omitempty"`
}

type Response struct {
	resp.Response
	Hash string `json:"hash"`
}

type TextSaver interface {
	SaveText(text string, ttl int) (string, error)
}

func New(log *slog.Logger, textSaver TextSaver, defaultTTL int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.text.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed to decode request"))

			return
		}

		log.Info("Request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("Invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		timeToLive := req.TTL
		if timeToLive == 0 {
			timeToLive = defaultTTL
		}

		hash, err := textSaver.SaveText(req.Text, timeToLive)
		if err != nil {
			log.Error("failed to save text", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed to save text"))

			return
		}

		log.Info("Text added", slog.String("hash", hash))

		ResponseOK(w, r, hash)
	}
}

func ResponseOK(w http.ResponseWriter, r *http.Request, hash string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Hash:     hash,
	})
}

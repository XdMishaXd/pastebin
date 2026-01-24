package save

import (
	"context"
	"log/slog"
	"net/http"

	resp "main_service/internal/lib/api/response"
	sl "main_service/internal/lib/logger"
	"main_service/internal/models"

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

// New godoc
// @Summary      Сохранить текст
// @Description  Сохранить текст и получить уникальный хеш для доступа
// @Tags         texts
// @Accept       json
// @Produce      json
// @Param        request  body      Request  true  "Текст для сохранения"
// @Success      201      {object}  Response "Успешно сохранено"
// @Failure      400      {object}  ErrorResponse "Некорректный запрос"
// @Failure      500      {object}  ErrorResponse "Внутренняя ошибка сервера"
// @Router       /text/save [post]

func New(ctx context.Context, log *slog.Logger, textSaver models.TextOperator, defaultTTL int) http.HandlerFunc {
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

		log.Info("Request body decoded")

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("Invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		timeToLive := req.TTL
		if timeToLive == 0 {
			timeToLive = defaultTTL
		}

		hash, err := textSaver.SaveText(ctx, req.Text, timeToLive)
		if err != nil {
			log.Error("failed to save text", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Internal error"))

			return
		}

		log.Info("Text added", slog.String("hash", hash))

		render.Status(r, http.StatusCreated)
		ResponseOK(w, r, hash)
	}
}

func ResponseOK(w http.ResponseWriter, r *http.Request, hash string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Hash:     hash,
	})
}

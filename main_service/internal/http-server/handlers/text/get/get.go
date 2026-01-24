package get

import (
	"context"
	"log/slog"
	"net/http"

	resp "main_service/internal/lib/api/response"
	sl "main_service/internal/lib/logger"
	"main_service/internal/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	Text string `json:"text"`
}

// New godoc
// @Summary      Получить текст
// @Description  Получить сохраненный текст по его уникальному хешу
// @Tags         texts
// @Accept       json
// @Produce      json
// @Param        hash  string  true  "Уникальный хеш текста" example:"abc123def456"
// @Success      200   {object}  Response "Текст успешно получен"
// @Failure      400   {object}  ErrorResponse "Хеш не указан или некорректен"
// @Failure      500   {object}  ErrorResponse "Внутренняя ошибка сервера"
// @Router       /text/{hash} [get]

func New(ctx context.Context, log *slog.Logger, textGetter models.TextOperator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.text.get.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		hash := chi.URLParam(r, "hash")
		if hash == "" {
			log.Info("Hash is empty")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("Hash is empty"))

			return
		}

		text, err := textGetter.GetText(ctx, hash)
		if err != nil {
			log.Error("failed to get text", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to get text"))

			return
		}

		log.Info("Text returned", slog.String("hash", hash))

		ResponseOK(w, r, text)
	}
}

func ResponseOK(w http.ResponseWriter, r *http.Request, text string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Text:     text,
	})
}

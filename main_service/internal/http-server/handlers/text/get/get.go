package get

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	resp "main_service/internal/lib/api/response"
	sl "main_service/internal/lib/logger"
	"main_service/internal/models"
	"main_service/internal/storage"

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
// @Description  Получает сохраненный текст по его уникальному хешу. Популярные тексты кэшируются в Redis для быстрого доступа.
// @Tags         texts
// @Accept       json
// @Produce      json
// @Param        hash  path  string  true  "Уникальный хеш текста (буквенно-цифровая строка)"  minlength(6)  maxlength(64)  example(a1b2c3d4e5f6)
// @Success      200   {object}  object{status=string,text=string}  "Текст успешно получен"  example({"status": "ok", "text": "Hello, World!"})
// @Failure      400   {object}  object{status=string,error=string}  "Хеш не указан"  example({"status": "error", "error": "Hash is empty"})
// @Failure      404   {object}  object{status=string,error=string}  "Текст не найден"  example({"status": "error", "error": "Text not found"})
// @Failure      500   {object}  object{status=string,error=string}  "Ошибка при получении текста"  example({"status": "error", "error": "Failed to get text"})
// @Router       /text/{hash} [get]
// @Security     none
// @x-order      2
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

			if errors.Is(err, storage.ErrTextNotFound) || errors.Is(err, storage.ErrTTLIsExpired) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp.Error("Text not found"))

				return
			}

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Failed to get text"))

			return
		}

		w.Header().Set("Cache-Control", "private, max-age=60")

		log.Info("Text got successfully", slog.String("hash", hash))

		ResponseOK(w, r, text)
	}
}

func ResponseOK(w http.ResponseWriter, r *http.Request, text string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Text:     text,
	})
}

package save

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/storage"
	"url-shortener/lib/api/response"
	"url-shortener/lib/random"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 7

type URLSaver interface {
	SaveURL(alias string, urlToSave string) (int64, error)
}

func New(log *slog.Logger, saver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save"
		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to parse request", "error", err)
			render.JSON(w, r, response.Error("failed to parse request"))

			return
		}
		log.Info("loaded request", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("failed to validate request", "error", err)
			render.JSON(w, r, response.Error("failed to validate request"))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		// add similiar
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
		id, err := saver.SaveURL(req.Url, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.Url))

			render.JSON(w, r, response.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", "error", err)

			render.JSON(w, r, response.Error("failed to add url"))

			return
		}
		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, alias)

	}

}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Alias:    alias,
	})
}

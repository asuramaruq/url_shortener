package delete

import (
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/asuramaruq/url_shortener/internal/lib/api/response"
	"github.com/asuramaruq/url_shortener/internal/lib/logger/sl"
	"github.com/asuramaruq/url_shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) (bool, error)
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("not found"))
			return
		}

		deleted, err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, resp.Error("url not found"))
			return
		}

		if err != nil {
			log.Error("failed to delete URL", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("url deleted", slog.Bool("deleted", deleted), slog.String("alias", alias))
		render.JSON(w, r, resp.OK())
	}
}

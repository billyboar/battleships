package v1

import (
	"context"
	"net/http"

	"github.com/billyboar/battleships/helpers"

	"github.com/billyboar/battleships/models"
)

type contextKey string

const (
	SessionCtx contextKey = "session"
)

// LoadSessionToCtx embeds session into ctx
func (api *APIServer) LoadSessionToCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.URL.Query().Get("session_id")
		events, err := api.Store.GetEvents(sessionID)
		if err != nil {
			helpers.RenderError(w, "cannot get session events", err, http.StatusBadRequest)
			return
		}

		session, err := models.BuildSessionEvents(events, sessionID)
		if err != nil {
			helpers.RenderError(w, "cannot build session", err, http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), SessionCtx, session)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *APIServer) GlobalCORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			headers := w.Header()
			headers.Add("Access-Control-Allow-Origin", "*")
			headers.Add("Vary", "Origin")
			headers.Add("Vary", "Access-Control-Request-Method")
			headers.Add("Vary", "Access-Control-Request-Headers")
			headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
			headers.Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

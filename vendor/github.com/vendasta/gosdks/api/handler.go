package api

import (
	"context"
	"encoding/json"
	"net/http"
)

type Response struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type Handler interface {
	Authorize(r *http.Request) error
	ParseArgs(r *http.Request) (map[string]interface{}, error)
	ValidateArgs(args map[string]interface{}) error
	Process(ctx context.Context, args map[string]interface{}) (string, error)
}

func Handle(h Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.Authorize(r); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Response{Error: err.Error()})
			return
		}
		args, err := h.ParseArgs(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Error: err.Error()})
			return
		}
		err = h.ValidateArgs(args)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Error: err.Error()})
			return
		}
		resp, err := h.Process(r.Context(), args)
		if err != nil {
			http.Error(w, "Error encountered", 500)
			return
		}
		if err = json.NewEncoder(w).Encode(Response{Data: resp}); err != nil {
			http.Error(w, "Error encountered", 500)
		}
	}
}

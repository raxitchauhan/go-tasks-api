package utils

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)


type (
	ErrorResponse struct {
		Errors []ErrorDescription `json:"errors"`
	}

	ErrorDescription struct {
		ID      string      `json:"id"`
		Code    string      `json:"code"`
		Status  int         `json:"status"`
		Title   string      `json:"title"`
		Details string      `json:"detail"`
		Source  *FieldError `json:"source,omitempty"`
	}

	// FieldError describes an error for a specific field, usually provided upon the request
	FieldError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
)

// WriteJSON encodes the provided data into JSON format and writes it into the given reader with
// Content-Type set to "application/json; charset=utf-8"
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if status != http.StatusNoContent {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Error().Err(err).Msg("failed to write error response")
		}
	}
}

// WriteJSONError sets the fields and writes them into the given reader as JSON with
// Content-Type set to "application/json; charset=utf-8"
func WriteJSONError(
	w http.ResponseWriter,
	status int,
	errDesc ErrorDescription,
	sources ...FieldError,
) {
	if errDesc.Title == "" {
		errDesc.Title = "an error occurred"
	}

	var errResps []ErrorDescription
	if len(sources) > 0 {
		errResps = make([]ErrorDescription, 0, len(sources))

		for i := range sources {
			resp := configureErrorResponse(errDesc, &sources[i])
			errResps = append(errResps, resp)
		}
	} else {
		resp := configureErrorResponse(errDesc, nil)
		errResps = []ErrorDescription{resp}
	}

	WriteJSON(w, status, ErrorResponse{Errors: errResps})
}

func configureErrorResponse(resp ErrorDescription, source *FieldError) ErrorDescription {
	if resp.ID == "" {
		resp.ID = uuid.New().String()
	}

	resp.Source = source

	return resp
}

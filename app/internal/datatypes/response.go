package datatypes

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type APIResponse struct {
	Code    int32       `json:"code"`
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseBadRequest(logger *zap.SugaredLogger, msg string, w http.ResponseWriter) {
	resp := APIResponse{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Errorw("Couldn't encode resources", zap.String("err", err.Error()))
	}
}

func NotFound(logger *zap.SugaredLogger, msg string, w http.ResponseWriter) {
	resp := APIResponse{
		Code:    http.StatusNotFound,
		Message: msg,
	}
	w.WriteHeader(http.StatusNotFound)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Errorw("Couldn't encode resources", zap.String("err", err.Error()))
	}
}

func UnprocessableEntity(logger *zap.SugaredLogger, msg string, w http.ResponseWriter) {
	resp := APIResponse{
		Code:    http.StatusUnprocessableEntity,
		Message: msg,
	}
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Errorw("Couldn't encode resources", zap.String("err", err.Error()))
	}
}

func ResponseInternalError(logger *zap.SugaredLogger, msg string, w http.ResponseWriter) {
	resp := APIResponse{
		Code:    http.StatusInternalServerError,
		Message: msg,
	}
	w.WriteHeader(http.StatusInternalServerError)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Errorw("Couldn't encode resources", zap.String("err", err.Error()))
	}
}

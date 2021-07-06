package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	dt "github.com/emorikvendy/url-shortener/internal/datatypes"
	"github.com/emorikvendy/url-shortener/internal/resources"
	sdt "github.com/emorikvendy/url-shortener/internal/sources/datatypes"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

func getStats(logger *zap.SugaredLogger, adapter sdt.StatsAdapter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			dt.UnprocessableEntity(logger, "ID must be an integer", w)
			logger.Debugw("ID must be an integer", zap.Any("request", r.Body))
			return
		}
		stats, err := adapter.GetByURLID(id)
		if err != nil {
			dt.ResponseInternalError(logger, err.Error(), w)
			logger.Errorw("Couldn't get stats", zap.String("err", err.Error()))
			return
		}

		resp := dt.APIResponse{
			Code: http.StatusOK,
			Data: *stats,
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			logger.Errorw("Couldn't encode resources", zap.String("err", err.Error()), zap.Any("request", r.Body))
		}
	}
}

func AddStatsRoutes(r *mux.Router, logger *zap.SugaredLogger, adapters resources.Adapters) {
	r.HandleFunc("/url/{id:[0-9]+}/stats", getStats(logger, adapters.Stats)).Methods(http.MethodGet)
}

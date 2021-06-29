package controllers

import (
	"encoding/json"
	"fmt"
	dt "github.com/emorikvendy/url-shortener/internal/datatypes"
	sdt "github.com/emorikvendy/url-shortener/internal/sources/datatypes"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strconv"
)

func addURL(logger *zap.SugaredLogger, adapter sdt.URLAdapter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		urlItem := new(dt.URL)
		err := json.NewDecoder(r.Body).Decode(urlItem)
		if err != nil {
			dt.ResponseBadRequest(logger, "Couldn't parse request body", w)
			logger.Infow("Couldn't parse request body", zap.Any("request", r.Body), zap.String("err", err.Error()))
			return
		}

		if !isURL(*urlItem.Link) {
			dt.UnprocessableEntity(logger, "URL must be valid and contain schema and host", w)
			logger.Infow("URL must be valid and contain schema and host", zap.Any("request", r.Body))
			return
		}

		if ok, err2 := adapter.URLExists(*urlItem.Link); ok {
			dt.UnprocessableEntity(logger, "URL already exists", w)
			logger.Infow("URL already exists", zap.Any("request", r.Body))
			return
		} else if err2 != nil {

		}
		err = adapter.Store(urlItem)
		if err != nil {
			dt.ResponseInternalError(logger, err.Error(), w)
			logger.Errorw("Couldn't save URL", zap.Any("request", r.Body), zap.String("err", err.Error()))
			return
		}

		resp := dt.APIResponse{
			Code: http.StatusCreated,
			Data: *urlItem,
		}
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			logger.Errorw("Couldn't encode resources", zap.String("err", err.Error()), zap.Any("request", r.Body))
		}
	}
}

func deleteURL(logger *zap.SugaredLogger, adapter sdt.URLAdapter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			dt.UnprocessableEntity(logger, "ID must be an integer", w)
			logger.Debugw("ID must be an integer", zap.Any("request", r.Body))
			return
		}
		err = adapter.Delete(id)
		if err != nil {
			dt.ResponseInternalError(logger, err.Error(), w)
			logger.Errorw("Couldn't delete URL", zap.Any("request", r.Body), zap.String("err", err.Error()))
			return
		}

		resp := dt.APIResponse{
			Code: http.StatusOK,
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			logger.Errorw("Couldn't encode resources", zap.String("err", err.Error()), zap.Any("request", r.Body))
		}
	}
}

func getURL(logger *zap.SugaredLogger, adapter sdt.URLAdapter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			dt.UnprocessableEntity(logger, "ID must be an integer", w)
			logger.Debugw("ID must be an integer", zap.Any("request", r.Body))
			return
		}
		urlItem, err := adapter.GetByID(id)
		if err != nil {
			dt.ResponseInternalError(logger, err.Error(), w)
			logger.Errorw("Couldn't delete URL", zap.Any("request", r.Body), zap.String("err", err.Error()))
			return
		}

		resp := dt.APIResponse{
			Code: http.StatusOK,
			Data: *urlItem,
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			logger.Errorw("Couldn't encode resources", zap.String("err", err.Error()), zap.Any("request", r.Body))
		}
	}
}

func updateURL(logger *zap.SugaredLogger, adapter sdt.URLAdapter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			dt.UnprocessableEntity(logger, "ID must be an integer", w)
			logger.Debugw("ID must be an integer", zap.Any("request", r.Body))
			return
		}
		newURLItem := new(dt.URL)
		err = json.NewDecoder(r.Body).Decode(newURLItem)
		if err != nil {
			dt.ResponseBadRequest(logger, "Couldn't parse request body", w)
			logger.Debugw("Couldn't parse request body", zap.Any("request", r.Body), zap.String("err", err.Error()))
			return
		}
		urlItem, err := adapter.GetByID(id)
		if err != nil {
			dt.ResponseInternalError(logger, "Couldn't save URL", w)
			logger.Errorw("Couldn't delete URL", zap.Any("request", r.Body), zap.String("err", err.Error()))
			return
		} else if urlItem == nil {
			dt.NotFound(logger, fmt.Sprintf("url with id=%d was not found", id), w)
			return
		}
		newURLItem.ID = &id
		err = adapter.Store(newURLItem)
		if err != nil {
			dt.ResponseInternalError(logger, "Couldn't save URL", w)
			logger.Errorw("Couldn't save URL", zap.Any("request", r.Body), zap.String("err", err.Error()))
			return
		}
		resp := dt.APIResponse{
			Code: http.StatusOK,
			Data: *urlItem,
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			logger.Errorw("Couldn't encode resources", zap.String("err", err.Error()), zap.Any("request", r.Body))
		}
	}
}

func redirect(logger *zap.SugaredLogger, adapter sdt.URLAdapter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		urlItem, err := adapter.GetByHash(vars["hash"])
		if err != nil {
			dt.ResponseInternalError(logger, err.Error(), w)
			logger.Errorw("Couldn't delete URL", zap.Any("request", r.Body), zap.String("err", err.Error()))
			return
		}
		if urlItem == nil {
			dt.NotFound(logger, fmt.Sprintf("url with hash=%s was not found", vars["hash"]), w)
			return
		}
		http.Redirect(w, r, *urlItem.Link, http.StatusTemporaryRedirect)
	}
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func AddURLRoutes(r *mux.Router, logger *zap.SugaredLogger, adapter sdt.URLAdapter, hashLen int) {
	r.HandleFunc(fmt.Sprintf("/r/{hash:[0-9a-z]{%d}}", hashLen), redirect(logger, adapter)).Methods(http.MethodGet)
	r.HandleFunc("/url", addURL(logger, adapter)).Methods(http.MethodPost)
	r.HandleFunc("/url/{id:[0-9]+}", deleteURL(logger, adapter)).Methods(http.MethodDelete)
	r.HandleFunc("/url/{id:[0-9]+}", updateURL(logger, adapter)).Methods(http.MethodPatch)
	r.HandleFunc("/url/{id:[0-9]+}", getURL(logger, adapter)).Methods(http.MethodGet)
}

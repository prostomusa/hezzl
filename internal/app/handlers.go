package app

import (
	"Project/internal/dto"
	"Project/internal/manager"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	DefaultLimitGoods  = 10
	DefaultOffsetGoods = 0
)

func (s *HezzlWebService) goodList(w http.ResponseWriter, r *http.Request) {
	limit, errParse := strconv.Atoi(r.URL.Query().Get("limit"))
	if errParse != nil {
		limit = DefaultLimitGoods
	}
	offset, errParse := strconv.Atoi(r.URL.Query().Get("offset"))
	if errParse != nil {
		offset = DefaultOffsetGoods
	}
	result, err := s.manager.GoodManager.ListGood(limit, offset)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (s *HezzlWebService) createGood(w http.ResponseWriter, r *http.Request) {
	var requestBody dto.CreateGoodRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	projectId, errParse := strconv.Atoi(r.URL.Query().Get("projectId"))
	if err != nil || requestBody.Name == "" || errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := s.manager.GoodManager.CreateGood(projectId, requestBody.Name)
	if err != nil {
		if errors.Is(err, manager.AlreadyExist) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("good с projectId = %v уже существует", projectId)))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
		}
		return
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (s *HezzlWebService) updateGood(w http.ResponseWriter, r *http.Request) {
	var requestBody dto.UpdateGoodRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	projectId, errParse := strconv.Atoi(r.URL.Query().Get("projectId"))
	if err != nil || requestBody.Name == "" || errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, errParse := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || requestBody.Name == "" || errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := s.manager.GoodManager.UpdateGood(id, projectId, requestBody)
	if err != nil {
		if errors.Is(err, manager.NotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Code: 3, Message: "errors.good.notFound"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
		}
		return
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (s *HezzlWebService) reprioritiize(w http.ResponseWriter, r *http.Request) {
	var requestBody dto.ReprioritizeGoodRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if requestBody.NewPriority <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	projectId, errParse := strconv.Atoi(r.URL.Query().Get("projectId"))
	if errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, errParse := strconv.Atoi(r.URL.Query().Get("id"))
	if errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := s.manager.GoodManager.ReprioritiizeGood(id, projectId, requestBody.NewPriority)
	if err != nil {
		if errors.Is(err, manager.NotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Code: 3, Message: "errors.good.notFound"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
		}
		return
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (s *HezzlWebService) deleteGood(w http.ResponseWriter, r *http.Request) {
	projectId, errParse := strconv.Atoi(r.URL.Query().Get("projectId"))
	if errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, errParse := strconv.Atoi(r.URL.Query().Get("id"))
	if errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := s.manager.GoodManager.DeleteGood(id, projectId)
	if err != nil {
		if errors.Is(err, manager.NotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Code: 3, Message: "errors.good.notFound"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
		}
		return
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

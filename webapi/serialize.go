package webapi

import (
	"encoding/json"
	"net/http"
)

type Serializer interface {
	JSON(rw http.ResponseWriter, status int, obj interface{})
}

type WebAPISerializer struct{}

func (w WebAPISerializer) JSON(rw http.ResponseWriter, status int, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)

	if obj != nil && status != http.StatusNoContent {
		rw.Write(data)
	}
	return
}

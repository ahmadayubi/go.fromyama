package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status int32 `json:"status"`
}

type missingParamError struct {
	s string
}

func (e *missingParamError) Error() string{
	return e.s
}

func ErrorResponse(w http.ResponseWriter, err error){
	http.Error(w, err.Error(), http.StatusBadRequest)
	return
}

func ForbiddenResponse(w http.ResponseWriter){
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(&Response{Message: "Unauthorized", Status: http.StatusForbidden}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusForbidden)
	w.Write(buf.Bytes())
}

func ObjectAddedToDatabase (w http.ResponseWriter, m string){
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(&Response{Message: m, Status: http.StatusCreated}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

func JSONResponse(w http.ResponseWriter, status int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func ParseRequestBody (r *http.Request, body *map[string]string, needed []string) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err = json.Unmarshal(reqBody, body);err != nil {
		return err
	}
	for key, val := range *body {
		if val == ""{
			return &missingParamError{key+" Missing"}
		}
	}
	for i := range needed {
		if (*body)[needed[i]] == ""{
			return &missingParamError{needed[i]+" Missing"}
		}
	}
	return nil
}

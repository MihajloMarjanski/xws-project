package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"user-service/model"
)

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) {

	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func decodeUserBody(r io.Reader) (*model.User, error) {

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt model.User
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func decodeInterestBody(r io.Reader) (*model.Interest, error) {

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt model.Interest
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func decodeExperienceBody(r io.Reader) (*model.Experience, error) {

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt model.Experience
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

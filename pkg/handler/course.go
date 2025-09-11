package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/juanrojas09/go_lib_response/response"
	"github.com/juanrojas09/gocourse_course/internal/courses"
)

func NewHttpServer(ctx context.Context, endpoints courses.Endpoints) http.Handler {
	r := mux.NewRouter()

	//funcion por defecto para manejo de errores

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	//declaracion de endpoints.
	r.Handle("/courses", httptransport.NewServer((endpoint.Endpoint)(endpoints.Create), decodeCourseCreation, encodeCourse)).Methods(http.MethodPost)

	r.Handle("/courses", httptransport.NewServer((endpoint.Endpoint)(endpoints.Get),
		decodeCourseGetAll, encodeCourse, opts...)).Methods(http.MethodGet)

	r.Handle("/courses/{id}", httptransport.NewServer((endpoint.Endpoint)(endpoints.GetById),
		decodeCourseGetById, encodeCourse, opts...)).Methods(http.MethodGet)

	r.Handle("/courses/{id}", httptransport.NewServer((endpoint.Endpoint)(endpoints.Update),
		decodeCourseUpdate, encodeCourse, opts...)).Methods(http.MethodPatch)

	r.Handle("/courses/{id}", httptransport.NewServer((endpoint.Endpoint)(endpoints.Delete),
		decodeCourseDelete, encodeCourse, opts...)).Methods(http.MethodDelete)

	return r
}

func decodeCourseGetAll(ctx context.Context, r *http.Request) (interface{}, error) {

	queryParams := r.URL.Query()
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	name := queryParams.Get("name")
	var namePtr *string
	if name != "" {
		namePtr = &name
	}
	startDate := queryParams.Get("start_date")
	var startDatePtr *string
	if startDate != "" {
		startDatePtr = &startDate
	}
	endDate := queryParams.Get("end_date")
	var endDatePtr *string
	if endDate != "" {
		endDatePtr = &endDate
	}

	req := courses.Filters{
		Name:      namePtr,
		StartDate: startDatePtr,
		EndDate:   endDatePtr,
		Limit:     &limit,
		Page:      &page,
	}

	return req, nil

}

func decodeCourseGetById(ctx context.Context, r *http.Request) (interface{}, error) {
	path := mux.Vars(r)
	id := path["id"]
	req := courses.GetRequest{
		ID: id,
	}

	return req, nil

}

func decodeCourseUpdate(ctx context.Context, r *http.Request) (interface{}, error) {
	var req courses.UpdateRequest
	path := mux.Vars(r)
	id := path["id"]
	req.ID = &id
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return req, nil
}

func encodeCourse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusByCode())
	return json.NewEncoder(w).Encode(r)
}

func decodeCourseDelete(ctx context.Context, r *http.Request) (interface{}, error) {
	path := mux.Vars(r)
	id := path["id"]
	req := courses.DeleteRequest{
		ID: id,
	}

	return req, nil
}

func decodeCourseCreation(ctx context.Context, r *http.Request) (interface{}, error) {
	var req courses.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

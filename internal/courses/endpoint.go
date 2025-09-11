package courses

import (
	"context"
	"errors"
	"log"

	"github.com/juanrojas09/go_lib_response/response"
)

type (
	Controller func(ctx context.Context, req interface{}) (interface{}, error)

	Endpoints struct {
		Create  Controller
		Get     Controller
		GetById Controller
		Update  Controller
		Delete  Controller
	}

	//Requests dtos
	CreateRequest struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	UpdateRequest struct {
		ID        *string `json:"id"`
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
	}

	Filters struct {
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
		Limit     *int    `json:"limit"`
		Page      *int    `json:"page"`
	}

	GetRequest struct {
		ID string `json:"id"`
	}

	DeleteRequest struct {
		ID string `json:"id"`
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create:  MakeCreateEndpoint(s),
		Get:     MakeGetEndpoint(s),
		GetById: MakeGetByIdEndpoint(s),
		Update:  MakeUpdateEndpoint(s),
		Delete:  MakeDeleteEndpoint(s),
	}
}

func MakeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, req interface{}) (interface{}, error) {

		r := req.(CreateRequest)
		//validar campos podriamos
		course, err := s.Create(ctx, r.Name, r.StartDate, r.EndDate)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		return response.Created("Course Created successfully", course, nil), nil
	}
}

func MakeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(Filters)

		courses, meta, err := s.GetAll(ctx, &r, *r.Page, *r.Limit)

		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Ok("course fetched successfully", courses, meta), nil
	}
}

func MakeGetByIdEndpoint(s Service) Controller {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(GetRequest)
		course, err := s.GetById(ctx, r.ID)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		return response.Ok("course fetched successfully", course, nil), nil
	}
}

func MakeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(UpdateRequest)
		log.Println("REQ", *r.ID)
		course, err := s.Update(ctx, *r.ID, &r)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		return response.Ok("course updated successfully", course, nil), nil
	}
}
func MakeDeleteEndpoint(s Service) Controller {

	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(DeleteRequest)
		id, err := s.Delete(ctx, r.ID)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		return response.Ok("course deleted successfully", id, nil), nil
	}

}

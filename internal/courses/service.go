package courses

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/juanrojas09/gocourse_domain/domain"
	"github.com/juanrojas09/gocourse_meta/meta"
)

type (
	Service interface {
		Create(ctx context.Context, Name string, StartDate, EndDate string) (*domain.Course, error)
		GetAll(ctx context.Context, filters *Filters, offset int, limit int) ([]domain.Course, *meta.Metadata, error)
		GetById(ctx context.Context, id string) (*domain.Course, error)
		Update(ctx context.Context, id string, data *UpdateRequest) (*domain.Course, error)
		Delete(ctx context.Context, id string) (string, error)
		Count(filters Filters) (int, error)
	}

	service struct {
		repo   CourseRepository
		logger *log.Logger
	}
)

func NewService(repo CourseRepository, logger *log.Logger) Service {
	return &service{

		repo: repo, logger: logger,
	}
}

func (s *service) Create(ctx context.Context, Name string, StartDate, EndDate string) (*domain.Course, error) {

	startDateParsed, err := time.Parse("2006-02-02", StartDate)
	if err != nil {
		s.logger.Println(err)
		return nil, err
	}
	emdDateParsed, err := time.Parse("2006-02-02", EndDate)
	if err != nil {
		s.logger.Println(err)
		return nil, err
	}
	course := domain.Course{
		Name:      Name,
		StartDate: startDateParsed,
		EndDate:   emdDateParsed,
	}

	isValid, fields := course.ValidateFields()

	if !isValid {
		return nil, fmt.Errorf("there are fields that are null: %v", fields)
	}

	res, err := s.repo.Create(ctx, &course)

	if err != nil {
		return nil, err
	}
	return res, nil

}

func (s *service) GetAll(ctx context.Context, filters *Filters, perPage int, limit int) ([]domain.Course, *meta.Metadata, error) {

	var startDateParsed, endDateParsed time.Time
	var err error

	//convertir fecha en el formato a time.Time
	if filters.StartDate != nil {
		startDateParsed, err = time.Parse("2006-01-02", *filters.StartDate)
		if err != nil {
			s.logger.Println(err)
			return nil, nil, ErrParse
		}
		formatted := startDateParsed.Format("2006-02-02")
		filters.StartDate = &formatted
	}

	if filters.EndDate != nil {
		endDateParsed, err = time.Parse("2006-01-02", *filters.EndDate)
		if err != nil {
			s.logger.Println(err)
			return nil, nil, ErrParse
		}
		formatted := endDateParsed.Format("2006-02-02")
		filters.EndDate = &formatted
	}

	total, err := s.repo.Count(*filters)
	if err != nil {
		return nil, nil, err
	}
	defaultPage := os.Getenv("PAGINATION_PER_PAGE_DEFAULT")
	meta, err := meta.New(total, limit, perPage, defaultPage)
	if err != nil {
		return nil, nil, err
	}

	courses, err := s.repo.GetAll(ctx, filters, meta.Offset(), meta.Limit())
	if err != nil {
		return nil, nil, err
	}
	return courses, meta, err

}
func (s *service) GetById(ctx context.Context, id string) (*domain.Course, error) {

	course, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return course, nil

}
func (s *service) Update(ctx context.Context, id string, data *UpdateRequest) (*domain.Course, error) {

	course, err := s.repo.Update(ctx, id, data)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (s *service) Delete(ctx context.Context, id string) (string, error) {
	id, err := s.repo.Delete(ctx, id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *service) Count(filters Filters) (int, error) {
	return s.repo.Count(filters)
}

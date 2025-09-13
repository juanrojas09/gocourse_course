package courses

import (
	"context"
	"log"

	"github.com/juanrojas09/gocourse_domain/domain"
	"gorm.io/gorm"
)

type (
	CourseRepository interface {
		Create(ctx context.Context, course *domain.Course) (*domain.Course, error)
		GetAll(ctx context.Context, filters *Filters, offset int, limit int) ([]domain.Course, error)
		GetById(ctx context.Context, id string) (*domain.Course, error)
		Update(ctx context.Context, id string, data *UpdateRequest) (*domain.Course, error)
		Delete(ctx context.Context, id string) (string, error)
		Count(filters Filters) (int, error)
	}

	repository struct {
		db     *gorm.DB
		logger *log.Logger
	}
)

func NewRepository(db *gorm.DB, log *log.Logger) CourseRepository {
	return &repository{db: db, logger: log}
}

func (r *repository) Create(ctx context.Context, course *domain.Course) (*domain.Course, error) {

	tx := r.db.Model(&course).Create(course)
	if tx.Error != nil {
		r.logger.Panicln("There was an error creating course")
		return nil, tx.Error
	}
	r.logger.Println("Course created")
	return course, nil

}

func (r *repository) GetAll(ctx context.Context, filters *Filters, offset int, limit int) ([]domain.Course, error) {
	var c []domain.Course
	tx := r.db.Model(&c)
	filters.Limit = nil
	filters.Page = nil
	log.Println(filters)
	tx = domain.ApplyFilters(tx, filters)
	tx = tx.Offset(offset).Limit(limit).Order("created_at desc").Find(&c)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, &ErrCourseNotFound{*filters.Name}
		}
		return nil, tx.Error
	}

	return c, nil

}
func (r *repository) GetById(ctx context.Context, id string) (*domain.Course, error) {
	var c domain.Course
	tx := r.db.Model(&c).Where(&domain.Course{ID: id}).First(&c)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, &ErrCourseNotFound{id}
		}
		return nil, tx.Error
	}
	return &c, nil
}
func (r *repository) Update(ctx context.Context, id string, data *UpdateRequest) (*domain.Course, error) {

	data.ID = nil
	err := domain.ApplyChanges(r.db, data, id)

	if err != nil {
		return nil, err
	}

	course, err := r.GetById(ctx, id)
	if err != nil {
		log.Println("There was an error fetching course by id after update:", err)
	}
	return course, err

}
func (r *repository) Delete(ctx context.Context, id string) (string, error) {
	var course domain.Course
	tx := r.db.Model(course).Delete(&domain.Course{ID: id})
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return "", &ErrCourseNotFound{id}
		}
		return "", tx.Error
	}
	return id, nil
}

func (r *repository) Count(filters Filters) (int, error) {
	var users domain.Course
	var count int64
	tx := r.db.Model(&users)
	filters.Limit = nil
	filters.Page = nil
	tx = domain.ApplyFilters(tx, &filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

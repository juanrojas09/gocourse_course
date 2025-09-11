package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/juanrojas09/gocourse_course/internal/courses"
	"github.com/juanrojas09/gocourse_domain/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDb() (*gorm.DB, error) {
	godotenv.Load()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	//cargar db
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	//validar el auto migrate

	if err != nil {
		return nil, err
	}

	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}

	if os.Getenv("DATABASE_MIGRATE") == "true" {
		err := db.AutoMigrate(domain.Course{})
		if err != nil {
			return nil, err
		}

	}

	return db, nil

}

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

}

func InitCourses(db *gorm.DB, log *log.Logger) courses.Endpoints {
	courseRepo := courses.NewRepository(db, log)
	courseSvc := courses.NewService(courseRepo, log)
	coursesEndpoints := courses.MakeEndpoints(courseSvc)
	return coursesEndpoints
}

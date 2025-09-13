package courses

import (
	"errors"
	"fmt"
)

var ErrDecodingDto = errors.New("Error decoding request dto")
var ErrParse = errors.New("Error parsing value")
var ErrUpdatingRecord = errors.New("Error updating record")
var ErrInvalidStartDate = errors.New("invalid start_date value")
var ErrInvalidEndDate = errors.New("invalid end_date value")
var ErrEndLesserStart = errors.New("start date cannot be greather than end date")

type ErrCourseNotFound struct {
	CourseId string
}

func (e ErrCourseNotFound) Error() string {
	return fmt.Sprintf("course '%s' not found", e.CourseId)
}

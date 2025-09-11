package courses

import (
	"errors"
	"fmt"
)

var ErrDecodingDto = errors.New("Error decoding request dto")
var ErrParse = errors.New("Error parsing value")
var ErrUpdatingRecord = errors.New("Error updating record")

type ErrNotFound struct {
	CourseId string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("course '%s' not found", e.CourseId)
}

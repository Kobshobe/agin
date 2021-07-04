package agin

import "errors"

func GetOffsetAndLimit(page, size, maxSize int) (offset int, limit int, err error)  {
	if size > maxSize {
		err = errors.New("size to l")
		return
	}
	offset = (page - 1) * size
	limit = size
	if offset < 0 || limit <= 0 {
		err = errors.New("err lo")
		return
	}
	return
}
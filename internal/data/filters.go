package data

import "strings"

type Filters struct {
	Page     int
	PageSize int
	Sort     string
}

func (f Filters) sortColumn() string {
	switch f.Sort {
	case "id", "-id":
		return "id"
	case "title", "-title":
		return "title"
	case "year", "-year":
		return "year"
	case "runtime", "-runtime":
		return "runtime"
	case "created_at", "-created_at":
		return "created_at"
	default:
		return "id"
	}
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

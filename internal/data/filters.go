package data

import "strings"

type Filters struct {
	Sort         string
	SortSafeList []string
}

func (f Filters) sortColumn() string {
	for _, field := range f.SortSafeList {
		if f.Sort == field {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort) // STOP SQL INJECTION
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

package dto

type SortDirection uint

const (
	SORT_UNKNOWN SortDirection = iota
	SORT_ASC
	SORT_DESC
)

var dirToString = map[SortDirection]string{
	SORT_UNKNOWN: "UNKNOWN",
	SORT_ASC:     "ASC",
	SORT_DESC:    "DESC",
}

// GetSortDirectionFromString Will return the sort direction from the string
func GetSortDirectionFromString(value string) SortDirection {
	for k, v := range dirToString {
		if v == value {
			return k
		}
	}
	return SORT_UNKNOWN
}

// String Will return the string representation of the enum
func (s SortDirection) String() string {
	str, ok := dirToString[s]
	if !ok {
		return dirToString[SORT_UNKNOWN]
	}
	return str
}

// Index Will return enum numeric index
func (s SortDirection) Index() int {
	return int(s)
}

func (s SortDirection) IsValid() bool {
	return s == SORT_ASC || s == SORT_DESC
}

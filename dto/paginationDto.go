package dto

type SortDirection string

const (
	ASC  SortDirection = "ASC"
	DESC SortDirection = "DESC"
)

type PaginationParam struct {
	Offset    int           `json:"offset" default:"1"`
	Limit     int           `json:"limit" default:"10"`
	OrderBy   string        `json:"orderBy,omitempty"`
	Direction SortDirection `json:"direction,omitempty"`
	ShowAll   bool          `json:"showAll,omitempty"`
	Search    string        `json:"search,omitempty"`
}

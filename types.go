package gohetznerdns

import "fmt"

const zonesBasePath = "/zones"
const recordsBasePath = "/records"

type ZoneList struct {
	Zones []*Zone `json:"zones"`
	Meta  *Meta   `json:"meta"`
}

type Zone struct {
	Id              *string   `json:"id"`
	Name            *string   `json:"name"`
	TTL             *int      `json:"ttl"`
	NS              []*string `json:"ns"`
	Paused          *bool     `json:"paused"`
	Status          *string   `json:"status"`
	NumberOfRecords *int      `json:"records_count"`
}

type ZoneRequest struct {
	Name *string `json:"name"`
	TTL  *int    `json:"ttl"`
}

type ZoneResponse struct {
	Zone  *Zone  `json:"zone"`
	Error *Error `json:"error"`
}

type Meta struct {
	Pagination *Pagination `json:"pagination"`
}

type Pagination struct {
	Page         *int `json:"page"`
	PerPage      *int `json:"per_page"`
	PreviousPage *int `json:"previous_page"`
	NextPage     *int `json:"next_page"`
	LastPage     *int `json:"last_page"`
	TotalEntries *int `json:"total_entries"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"string"`
}

func (e *Error) Error() error {
	return fmt.Errorf("%d : %s", e.Code, e.Message)
}

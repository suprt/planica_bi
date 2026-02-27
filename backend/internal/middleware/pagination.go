package middleware

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// Pagination represents pagination parameters
type Pagination struct {
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Sort    string `json:"sort"`
	Order   string `json:"order"`
	Offset  int    `json:"-"`
}

// DefaultPagination returns default pagination parameters
func DefaultPagination() *Pagination {
	return &Pagination{
		Page:    1,
		PerPage: 20,
		Sort:    "id",
		Order:   "ASC",
		Offset:  0,
	}
}

// GetPagination extracts pagination parameters from request context
func GetPagination(c echo.Context) *Pagination {
	p := DefaultPagination()

	// Parse page
	if page := c.QueryParam("page"); page != "" {
		if pageNum, err := strconv.Atoi(page); err == nil && pageNum > 0 {
			p.Page = pageNum
		}
	}

	// Parse per_page
	if perPage := c.QueryParam("per_page"); perPage != "" {
		if pp, err := strconv.Atoi(perPage); err == nil && pp > 0 && pp <= 100 {
			p.PerPage = pp
		}
	}

	// Parse sort
	if sort := c.QueryParam("sort"); sort != "" {
		p.Sort = sort
	}

	// Parse order
	if order := c.QueryParam("order"); order != "" {
		if order == "ASC" || order == "DESC" {
			p.Order = order
		}
	}

	// Calculate offset
	p.Offset = (p.Page - 1) * p.PerPage

	return p
}

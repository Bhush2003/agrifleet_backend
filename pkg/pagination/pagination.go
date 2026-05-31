package pagination

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Params holds parsed pagination parameters.
type Params struct {
	Page     int
	PageSize int
	Offset   int
}

// FromContext parses page and page_size query params from the request.
func FromContext(c *fiber.Ctx) Params {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", strconv.Itoa(DefaultPageSize)))

	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}

	return Params{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}

// TotalPages calculates the number of pages given total records and page size.
func TotalPages(total int64, pageSize int) int {
	if pageSize == 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

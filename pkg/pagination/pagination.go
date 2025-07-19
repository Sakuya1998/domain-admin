package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Offset  int
	Limit   int
	OrderBy string
	Sort    string // ASC or DESC
}

type PageResult struct {
	Total int64       `json:"total"`
	Items interface{} `json:"items"`
}

func New(c *gin.Context) Pagination {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")
	orderBy := c.DefaultQuery("orderBy", "id")
	sort := c.DefaultQuery("sort", "asc")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	return Pagination{
		Offset:  offset,
		Limit:   limit,
		OrderBy: orderBy,
		Sort:    sort,
	}
}

func (p Pagination) GetOrderClause() string {
	return p.OrderBy + " " + p.Sort
}

func NewPageResult(total int64, items interface{}) PageResult {
	return PageResult{
		Total: total,
		Items: items,
	}
}

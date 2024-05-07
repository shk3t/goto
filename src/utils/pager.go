package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Pager struct {
	Start int
	Take  int
}

func NewPager(c *fiber.Ctx) *Pager {
	start, _ := strconv.Atoi(c.Query("start"))
	take, _ := strconv.Atoi(c.Query("take"))

	pager := &Pager{Start: start, Take: take}
	if take == 0 {
		pager.Take = 10
	}

	return pager
}

func (p *Pager) QuerySuffix() string {
	return " LIMIT " + strconv.Itoa(p.Take) + " OFFSET " + strconv.Itoa(p.Start)
}
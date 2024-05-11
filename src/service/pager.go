package service

import (
	sc "strconv"

	"github.com/gofiber/fiber/v2"
)

type Pager struct {
	Start int
	Take  int
}

func NewPager(fctx *fiber.Ctx) *Pager {
	start, _ := sc.Atoi(fctx.Query("start"))
	take, _ := sc.Atoi(fctx.Query("take"))

	pager := &Pager{Start: start, Take: take}
	if take == 0 {
		pager.Take = 10
	}

	return pager
}

func (p *Pager) QuerySuffix() string {
	return " LIMIT " + sc.Itoa(p.Take) + " OFFSET " + sc.Itoa(p.Start)
}
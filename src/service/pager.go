package service

import (
	sc "strconv"

	"github.com/gofiber/fiber/v2"
)

type Pager struct {
	start       int
	take        int
	QuerySuffix string
}

func NewPager(fctx *fiber.Ctx) *Pager {
	start, _ := sc.Atoi(fctx.Query("start"))
	take, _ := sc.Atoi(fctx.Query("take"))

	pager := &Pager{start: start, take: take}
	if take == 0 {
		pager.take = 10
	}

	pager.QuerySuffix = " LIMIT " + sc.Itoa(pager.take) + " OFFSET " + sc.Itoa(pager.start)

	return pager
}
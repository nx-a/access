package access

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"time"
)

type Access struct {
	storage *session.Store
}

func New() *Access {
	return &Access{
		storage: session.New(session.Config{
			Expiration: time.Hour,
		}),
	}
}

func (a *Access) Middleware(access, url string) func(ctx *fiber.Ctx) error {
	num := accessNumber(access)
	return func(ctx *fiber.Ctx) error {
		log.Debug(access, "\t", url)
		if num > 8 {
			return ctx.Next()
		}
		_sess, err := a.storage.Get(ctx)
		if err != nil {
			return fiber.ErrForbidden
		}
		_access := _sess.Get("access")
		if _access == nil {
			return fiber.ErrForbidden
		}
		_acc, ok := _access.(uint16)
		if !ok {
			return fiber.ErrForbidden
		}
		if num >= _acc {
			return ctx.Next()
		}
		return fiber.ErrForbidden
	}
}

func accessNumber(access string) uint16 {
	switch access {
	case "admin":
		return 1
	case "manager":
		return 2
	case "user":
		return 4
	case "tech":
		return 8
	}
	return 1024
}

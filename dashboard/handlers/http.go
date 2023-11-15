package handlers

import (
	"marketingBot/dashboard/adapters"

	"github.com/gofiber/fiber/v2"
)

type MetricsGetter interface {
	GetClickCount(plataform, questionKey, optionKey string) int
	GetRevelCount(plataform, code string) int
}

type dashHttpApp struct {
	metricsGetter MetricsGetter
}

func NewdashHttpApp() *dashHttpApp {
	return &dashHttpApp{
		metricsGetter: adapters.StatisticsRepoMemory,
	}
}
func (d *dashHttpApp) HandleClickCount(c *fiber.Ctx) error {
	plataform := c.Query("plataform")
	question := c.Query("question")
	option := c.Query("option")
	count := d.metricsGetter.GetClickCount(plataform, question, option)
	return c.JSON(fiber.Map{"clicks": count})
}

func (d *dashHttpApp) HandleCouponRevelCount(c *fiber.Ctx) error {
	plataform := c.Query("plataform")
	coupon_code := c.Query("code")
	count := d.metricsGetter.GetRevelCount(plataform, coupon_code)
	return c.JSON(fiber.Map{"revels": count})
}

package handlers

import (
	"marketingBot/dashboard/adapters"

	"github.com/gofiber/fiber/v2"
)

type MetricsGetter interface {
	GetClickCount(plataform, key string) int
}

type dashHttpApp struct {
	metricsGetter MetricsGetter
}

func NewdashHttpApp() *dashHttpApp {
	return &dashHttpApp{
		metricsGetter: adapters.ButtonStatisticsRepoMemory,
	}
}
func (d *dashHttpApp) HandleClickCount(c *fiber.Ctx) error {
	plataform := c.Query("plataform")
	key := c.Query("key")
	count := d.metricsGetter.GetClickCount(plataform, key)
	return c.JSON(fiber.Map{"clicks": count})
}

package handlers

import (
	"marketingBot/dashboard/adapters"
	"marketingBot/dashboard/flow"
	"net/http"

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

type flowRepo interface {
	Save(*flow.Flow) error
	Get(key string) *flow.Flow
}

type FlowHttpApp struct {
	repo flowRepo
}

func NewFlowApp() *FlowHttpApp {
	return &FlowHttpApp{repo: adapters.MemoryFlowRepo}
}

func (f *FlowHttpApp) HandleSaveFlow(c *fiber.Ctx) error {
	var fl flow.Flow
	err := c.BodyParser(&fl)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	err = f.repo.Save(&fl)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.JSON(fiber.Map{"sucess": true})
}

func (f *FlowHttpApp) HandleGetFLow(c *fiber.Ctx) error {
	flowKey := c.Params("key")
	fl := f.repo.Get(flowKey)
	if fl == nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "not found"})
	}

	return c.JSON(fiber.Map{"name": fl.Name, "key": fl.Key})
}

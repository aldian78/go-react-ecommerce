package handler

import (
	"log"
	"net/http"

	"github.com/aldian78/go-react-ecommerce/backend/internal/dto"
	"github.com/aldian78/go-react-ecommerce/backend/internal/service"
	"github.com/gofiber/fiber/v2"
)

type webhookHandler struct {
	webhookService service.IWebhookService
}

func (wh *webhookHandler) ReceiveInvoice(c *fiber.Ctx) error {
	var request dto.XenditInvoiceRequest
	err := c.BodyParser(&request)
	if err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusBadRequest)
	}

	err = wh.webhookService.ReceiveInvoice(c.UserContext(), &request)
	if err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}

func NewWebhookHandler(webhookService service.IWebhookService) *webhookHandler {
	return &webhookHandler{
		webhookService: webhookService,
	}
}

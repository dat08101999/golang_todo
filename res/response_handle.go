package res

import (
	"api_fiber/models"

	"github.com/gofiber/fiber/v2"
)

func Response(c *fiber.Ctx, data interface{}, err error, successMessage string) error {
	if err == nil {
		var success models.SuccessResponseModel = models.SuccessResponseModel{
			Code:    200,
			Message: successMessage,
			Status:  1,
			Data:    data,
		}
		return c.JSON(success)
	} else {
		var errorModel models.ErrorResponse = models.ErrorResponse{
			Code:    0,
			Message: err.Error(),
			Data: map[string]interface{}{
				"err": c.JSON(err),
			},
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorModel)
	}
}

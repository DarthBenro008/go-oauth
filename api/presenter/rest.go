package presenter

import "github.com/gofiber/fiber/v2"

func Success(data interface{}, description interface{}) *fiber.Map {
	if description == nil {
		description = "Successfully fetched data"
	}
	return &fiber.Map{
		"status":      true,
		"data":        data,
		"description": description,
	}
}

func Failure(description interface{}) *fiber.Map {
	return &fiber.Map{
		"status":      false,
		"data":        nil,
		"description": description,
	}
}

func TokenResponse(token string, data interface{}) *fiber.Map {
	return &fiber.Map{
		"status": true,
		"token":  token,
		"data": data,
	}
}

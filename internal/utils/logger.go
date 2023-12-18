package utils

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
	"time"
)

func LoggerConfig() logger.Config {
	// Create or open a log file
	filePath := "logs/" + time.Now().Format("2006-01-02") + ".log"
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %s", err)
	}

	return logger.Config{
		Format:     "${time} | ${ip} | ${method} | ${path} | ${status} | ${latency} |\nHeaders: ${reqHeaders}\nBody: ${body}\n\n",
		TimeFormat: "2006/01/02 15:04:05",
		TimeZone:   "Europe/Istanbul",
		Output:     file,
		CustomTags: map[string]logger.LogFunc{
			"body": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				var body interface{}
				if err := json.Unmarshal(c.Body(), &body); err != nil {
					return output.Write(c.Body())
				}

				if bodyMap, ok := body.(map[string]interface{}); ok {
					if _, ok := bodyMap["phone_number"]; ok {
						bodyMap["phone_number"] = "********"
					}
					if _, ok := bodyMap["password"]; ok {
						bodyMap["password"] = "********"
					}
					if _, ok := bodyMap["new_password"]; ok {
						bodyMap["new_password"] = "********"
					}
					if _, ok := bodyMap["new_password_confirm"]; ok {
						bodyMap["new_password_confirm"] = "********"
					}

					body, _ = json.Marshal(bodyMap)
					return output.Write(body.([]byte))
				}

				return output.Write(c.Body())
			},
		},
	}
}

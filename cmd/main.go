package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"log"
	"net/http"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/segment"
	"github.com/gofiber/fiber/v2"
	"github.com/otiai10/gosseract/v2"
	"github.com/tidwall/gjson"
)

func main() {
	app := fiber.New()
	app.Post("/ocr", func(c *fiber.Ctx) error {
		request := string(c.Body())
		if gjson.Valid(request) {
			// Get base64 from json request
			base64image := gjson.Get(request, "img").String()

			// Decode base64 to byte
			sDec, err := base64.StdEncoding.DecodeString(base64image)
			if err != nil {
				log.Fatal(err)
			}

			// Decode byte to image struct
			img, _, err := image.Decode(bytes.NewReader(sDec))
			if err != nil {
				log.Fatalln(err)
			}

			// Convert Image to grayscale
			grayscale := effect.Grayscale(img)

			// Convert Image to threshold segment
			threshold := segment.Threshold(grayscale, 128)

			// Convert Image to Bytes
			buf := new(bytes.Buffer)
			jpeg.Encode(buf, threshold, nil)

			// Initiation Gosseract new client
			client := gosseract.NewClient()

			// close client when the main function is finished running
			defer client.Close()

			// Read byte to image and set whitelist character
			client.SetImageFromBytes(buf.Bytes())
			client.SetWhitelist(" -:/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

			// Get text result from OCR
			text, _ := client.Text()

			// return the response
			return c.JSON(&fiber.Map{"status": "OK", "response": text})
		}
		return c.Status(http.StatusNotAcceptable).JSON(&fiber.Map{"status": "Request not JSON"})
	})
	log.Fatal(app.Listen(":8000"))
}

package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/segment"
	"github.com/otiai10/gosseract/v2"
)

func main() {
	var PANPath, AAdharCardPath string
	flag.StringVar(&PANPath, "p", "nil", "path to PAN Card image")
	flag.StringVar(&AAdharCardPath, "a", "nil", "path to AAdhar Card Image")
	flag.Parse()
	if PANPath != "nil" {
		fmt.Println("PAN Card image passed")
		PANCARDNumber := validate(PANPath, "PAN")
		fmt.Println("PAN Card Number: ", PANCARDNumber)
	} else {
		fmt.Println("PAN Card details not passed")
	}
	if AAdharCardPath != "nil" {
		fmt.Println("AAdhar Card image passed")
		AAdharCARDNumber := validate(AAdharCardPath, "AADHAR")
		fmt.Println("AAdhar Card Number: ", AAdharCARDNumber)
	} else {
		fmt.Println("AAdhar card details not passed")
	}

}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func validate(path, cardType string) string {
	// Get base64 from json request
	byt, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var pattern string
	var base64image string

	// Append the base64 encoded output
	base64image = toBase64(byt)
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
	client.SetLanguage("eng", "hin")
	// Get text result from OCR
	text, _ := client.Text()
	// return the response
	text = strings.Replace(text, "\n", " ", -1)
	if cardType == "AADHAR" {
		pattern = `\d{4}\s\d{4}\s\d{4}`
	}
	if cardType == "PAN" {
		pattern = "[A-Z]{5}[0-9]{4}[A-Z]{1}"
	}
	re := regexp.MustCompile(pattern)
	match := re.FindAllString(text, -1)
	// fmt.Println(text)
	//fmt.Println(match[0])
	return match[0]
}

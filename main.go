package main

import (
	"GreatAdventure/num2word"
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/fogleman/gg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	fmt.Println("Бот запущен")
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				captchaQuestion, answer := generateCaptcha()
				imageBytes, err := generateCaptchaImage(captchaQuestion)
				if err != nil {
					log.Println("Error generating image:", err)
					continue
				}

				photoFileBytes := tgbotapi.FileBytes{
					Name:  "captcha.png",
					Bytes: imageBytes,
				}
				msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, photoFileBytes)
				msg.Caption = "Пожалуйста, решите капчу, чтобы продолжить. Вводить ответ нужно цифрами"
				photoResponse, err := bot.Send(msg)
				if err != nil {
					log.Println("Error sending message:", err)
					continue
				}

				for update := range updates {
					if update.Message.Text == answer {
						thankYouMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Спасибо! Вы подтвердили, что вы не робот.")
						bot.Send(thankYouMsg)

						deleteConfig := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, photoResponse.MessageID)
						bot.Send(deleteConfig)

						break
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неправильный ответ. Попробуйте еще раз. Помните, что ответ нужно вводить цифрами")
						bot.Send(msg)
					}
				}

			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда. Попробуйте /start.")
				bot.Send(msg)
			}
		}
	}
}

func generateCaptcha() (string, string) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	num1 := rand.Intn(20)
	num2 := rand.Intn(20)
	firstNum, secondNum := num2word.Converter(num1, num2)
	question := fmt.Sprintf("%s + %s", firstNum, secondNum)
	answer := strconv.Itoa(num1 + num2)
	return question, answer
}

func generateCaptchaImage(text string) ([]byte, error) {
	const width = 500
	const height = 500

	dc := gg.NewContext(width, height)
	r := rand.Float64()
	g := rand.Float64()
	b := rand.Float64()
	dc.SetRGB(r, g, b)
	dc.Clear()

	for i := 0; i < 10000; i++ {
		x := rand.Float64() * width
		y := rand.Float64() * height
		dc.SetRGBA(0, 0, 0, 0.1)
		dc.DrawPoint(x, y, 1)
		dc.Fill()
	}

	dc.Push()
	angle := (rand.Float64() - 0.5) * (gg.Radians(30))
	dc.Rotate(angle)
	dc.SetColor(color.Black)
	fontPath := "resources/fonts/creepster.otf"
	if err := dc.LoadFontFace(fontPath, 48); err != nil {
		return nil, err
	}
	dc.DrawStringAnchored(text, width/2, height/2, 0.5, 0.5)
	dc.Pop()

	for i := 0; i < 5; i++ {
		x1 := rand.Float64() * width
		y1 := rand.Float64() * height
		x2 := rand.Float64() * width
		y2 := rand.Float64() * height
		dc.SetRGBA(0, 0, 0, 0.3)
		dc.SetLineWidth(1)
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	var imageBuffer bytes.Buffer
	if err := dc.EncodePNG(&imageBuffer); err != nil {
		return nil, err
	}
	return imageBuffer.Bytes(), nil
}

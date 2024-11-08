package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Типы для обработки данных, получаемых от API
type Chat struct {
	ID int64 `json:"id"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Response struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
}

// типы для создания клавиатуры
type KeyboardButton struct {
	Text string `json:"text"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
}

// Функция получения токена из переменной окружения
func getTOKEN() (string, error) {
	// загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		panic("Ошибка загрузки переменной окружения")
	}
	// загружаем токен бота из переменных окружения
	TOKEN, err := os.LookupEnv("TOKEN")
	if !err {
		panic("Ошибка загрузки переменной окружения 'TOKEN'")
	}

	return TOKEN, nil
}

// Функция для проверки токена аутентификации бота
func getMe(APIURL string) (bool, error) {
	resp, err := http.Get(APIURL + "getMe")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var user User

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return false, err
	}

	return true, nil
}

// Функция для получения обновлений
func getUpdates(offset int, APIURL string) ([]Update, error) {
	resp, err := http.Get(APIURL + fmt.Sprintf("getUpdates?offset=%d", offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response Response
	fmt.Println("Отправили get запрос")

	// Декодируем ответ
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if len(response.Result) != 0 {
		fmt.Println("Текст сообщения:", response.Result[0].Message.Text)
	}
	if !response.OK {
		return nil, fmt.Errorf("error: unable to get updates")
	}

	return response.Result, nil
}

// Функция для отправки сообщения
func sendMessage(chatID int64, text string, replyMarkup interface{}, APIURL string) error {
	apiURL := APIURL + "sendMessage"

	requestBody := map[string]interface{}{
		"chat_id":      chatID,
		"text":         text,
		"reply_markup": replyMarkup,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// Отправляем POST-запрос
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Основная функция
func main() {
	var offset int

	TOKEN, err := getTOKEN()
	if err != nil {
		panic(err)
	}
	APIURL := "https://api.telegram.org/bot" + TOKEN + "/"

	// Пингуем
	getMe(APIURL)

	for {
		updates, err := getUpdates(offset, APIURL)
		if err != nil {
			log.Println("Error fetching updates:", err)
			continue
		}

		for _, update := range updates {
			// Обновляем offset, чтобы не обрабатывать одно и то же сообщение дважды
			offset = update.UpdateID + 1

			chatID := update.Message.Chat.ID
			text := update.Message.Text

			switch text {
			case "/start":
				// Reply Keyboard
				replyKeyboard := ReplyKeyboardMarkup{
					Keyboard: [][]KeyboardButton{
						{
							{"Кнопка 1"},
							{"Кнопка 2"},
						},
						{
							{"Кнопка 3"},
						},
					},
					ResizeKeyboard:  true,
					OneTimeKeyboard: true,
				}

				if err := sendMessage(chatID, "Выберите опцию:", replyKeyboard, APIURL); err != nil {
					log.Println("Error sending message:", err)
				}
			default:
				reply := fmt.Sprintf("Вы сказали: %s", text)
				replyKeyboard := ReplyKeyboardMarkup{
					Keyboard: [][]KeyboardButton{
						{
							{"Кнопка 1"},
							{"Кнопка 2"},
						},
						{
							{"Кнопка 3"},
						},
					},
					ResizeKeyboard:  true,
					OneTimeKeyboard: true,
				}
				if err := sendMessage(chatID, reply, replyKeyboard, APIURL); err != nil {
					log.Println("Error sending message:", err)
				}
			}
		}
	}
}

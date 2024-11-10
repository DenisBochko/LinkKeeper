package TGinter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	db "LinkKeeper/database"

	"github.com/joho/godotenv"
)

type TGinter struct{
	OK bool
}

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

// Функция запуска телеграм-интерфейса
func (t TGinter) Start(ctx context.Context, saveChan, getChan, deleteChan chan<- db.Field, receiveChan <-chan []db.Field) {
	offset := 0
	timeout := 60

	TOKEN, err := t.getTOKEN()
	if err != nil {
		fmt.Println("telegramInterface: Ошибка получения токена")
	}
	APIURL := "https://api.telegram.org/bot" + TOKEN + "/"
loop:
	for {
		select {

		case <-ctx.Done():
			fmt.Println("telegramInterface: Отключаем работу TG бота")
			break loop

		default:
			// Пытаемся получить updates
			updates, err := t.getUpdates(offset, timeout, APIURL)
			if err != nil {
				log.Println("Error fetching updates:", err)
				continue
			}

			for _, update := range updates {
				// Обновляем offset, чтобы не обрабатывать одно и то же сообщение дважды
				offset = update.UpdateID + 1

				chatID := update.Message.Chat.ID
				text := update.Message.Text
				// Reply Keyboard
				replyKeyboard := ReplyKeyboardMarkup{
					Keyboard: [][]KeyboardButton{
						{
							{"/delete"},
							//{"Кнопка 2"},
						},
						{
							{"/get"},
						},
					},
					ResizeKeyboard:  true,
					OneTimeKeyboard: true,
				}
				// содержит ли text подстроку https://
				if strings.Contains(text, "https://") {
					saveChan <- db.Field{
						ID: 0,
						UserID: strconv.Itoa(int(chatID)),
						UserURL: text,
					}
					if err := t.sendMessage(chatID, "Вы успешно сохранили ссылку!", replyKeyboard, APIURL); err != nil {
						log.Println("Error sending message:", err)
					}
				} else {
					switch text {
					case "/delete":
						deleteChan <- db.Field{
							ID: 0,
							UserID: strconv.Itoa(int(chatID)),
							UserURL: text,
						}
						if err := t.sendMessage(chatID, "Вы успешно удалили ваши ссылки!", replyKeyboard, APIURL); err != nil {
							log.Println("Error sending message:", err)
						}
					case "/get":
						getChan <- db.Field{ID: 0, UserID: strconv.Itoa(int(chatID)), UserURL: text}
						select {
						case fields := <-receiveChan:
							t.sendMessage(chatID, fmt.Sprint("Вот ваши ссылки: ", fields), replyKeyboard, APIURL)
						case <-time.After(5 * time.Second):
							t.sendMessage(chatID, "Не удалось получить ссылки, попробуйте позже.", replyKeyboard, APIURL)
						}
					}
				}
			}
		}
	}
}

// Функция получения токена из переменной окружения
func (t TGinter) getTOKEN() (string, error) {
    if err := godotenv.Load(); err != nil {
        return "", fmt.Errorf("ошибка загрузки переменной окружения: %w", err)
    }
    TOKEN, exists := os.LookupEnv("TOKEN")
    if !exists {
        return "", fmt.Errorf("переменная окружения 'TOKEN' не найдена")
    }
    return TOKEN, nil
}

// Функция для проверки токена аутентификации бота
// func (t TGinter) getMe(APIURL string) (bool, error) {
// 	resp, err := http.Get(APIURL + "getMe")
// 	if err != nil {
// 		return false, err
// 	}
// 	defer resp.Body.Close()

// 	var user User

// 	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
// 		return false, err
// 	}

// 	return true, nil
// }

// Функция для получения обновлений
func (t TGinter) getUpdates(offset int, timeout int, APIURL string) ([]Update, error) {
	url := APIURL + fmt.Sprintf("getUpdates?offset=%d&timeout=%d", offset, timeout)

	// Устанавливаем контекст с таймаутом на 10 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Создаем запрос с контекстом
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Отправляем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Декодируем ответ
	var result struct {
		OK     bool     `json:"ok"`
		Result []Update `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

// Функция для отправки сообщения
func (t TGinter) sendMessage(chatID int64, text string, replyMarkup interface{}, APIURL string) error {
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

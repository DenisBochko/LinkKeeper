package analyzer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Field struct {
	CHATID       int64
	Urls         []string
	ResponseText string
}

type Analyzer struct {
	OK bool
}

func (a Analyzer) Start(ctx context.Context, inputChan <-chan Field, outputChan chan<- Field) error {
	defer close(outputChan)

	workerChan := make(chan Field, 100)
	defer close(workerChan)

	go a.worker(workerChan, outputChan, "http://localhost:1337/v1/chat/completions")
	go a.worker(workerChan, outputChan, "http://localhost:1338/v1/chat/completions")
	go a.worker(workerChan, outputChan, "http://localhost:1339/v1/chat/completions")

loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("analyzer: Завершаем работу", ctx.Err())
			break loop
		case body := <-inputChan:
			workerChan <- body
		}
	}
	return nil
}

// host: http://localhost:1337/v1/chat/completions
func (a Analyzer) worker(inputChan <-chan Field, outputChan chan<- Field, hostUrl string) error {
	defer close(outputChan)
	for task := range inputChan {
		fmt.Println("Таска захвачена воркером: ", hostUrl)
		select {
		case <-time.After(20 * time.Second):
			// хуета не работает это потому что я еблан
			fmt.Print("Превышено время ожидания ответа сервера! (20 секунд)")
			outputChan <- Field{
				CHATID:       task.CHATID,
				Urls:         task.Urls,
				ResponseText: "Превышено время ожидания ответа сервера! (20 секунд)",
			}
			continue
		default:
			response, err := a.request(task.Urls, hostUrl)
			if err != nil {
				fmt.Println(err)
				return err
			}
			outputChan <- Field{
				CHATID:       task.CHATID,
				Urls:         task.Urls,
				ResponseText: response,
			}
		}
	}
	return nil
}

// функция выполнения запроса
func (a Analyzer) request(links []string, url string) (string, error) {
	// Создаём контент запроса из ссылок
	requestString := "Проанализируй мои ссылки, опиши мои предпочтения и предложи мне новые, сновываясь на моих предпочтениях: "
	for _, link := range links {
		requestString += link + " "
	}

	// Структура данных, которую мы будем отправлять в JSON
	data := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": requestString,
			},
		},
		"model": "gpt-4o-mini",
	}

	// Сериализуем данные в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "analyzer: Ошибка при сериализации данных в JSON", err
	}

	// "http://localhost:1337/v1/chat/completions"
	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "analyzer: Ошибка создания запроса", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос с использованием HTTP-клиента
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "analyzer: Ошибка выполнения запроса", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "analyzer: Ошибка чтения ответа", err
	}

	// Создаем переменную для хранения распарсенных данных
	var result map[string]interface{}

	// Декодируем JSON в map
	if err := json.Unmarshal(body, &result); err != nil {
		return "analyzer: Ошибка декодирования", err
	}
	fmt.Println(result)
	// Извлекаем значение content
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choiceMap, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choiceMap["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}
	return "", nil
	//provider:RubiksAI ChatGptEs
}

func (a Analyzer) Printer(inputChan chan Field) {
	for v := range inputChan {
		fmt.Print(v.CHATID, "\n", v.ResponseText, "\n\n")
	}
}
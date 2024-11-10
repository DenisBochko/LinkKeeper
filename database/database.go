package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

type DataBase struct {
	ConnStr    string
	DriverName string
}

type Field struct {
	ID           int
	UserID       string
	UserURL      string
	DeleteNumber int
}

// connStr := "user=postgres dbname=LinkKeeper password=postgres host=localhost sslmode=disable"
// driver = "postgres"

func (d DataBase) Start(ctx context.Context, saveChan, getChan, deleteChan, deleteOfItemChan <-chan Field, receive chan<- []Field) {
	db, err := sql.Open(d.DriverName, d.ConnStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	defer close(receive)

loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("database: Закрываем соединение с БД")
			break loop
		case field := <-saveChan:
			_, err := d.Save(db, field.UserID, field.UserURL)
			if err != nil {
				fmt.Printf("database: Не удалось сохранить запись пользователя: %s \n", field.UserID)
			} else {
				fmt.Printf("database: Удалось сохранить запись пользователя: %s \n", field.UserID)
			}
		case field := <-getChan:
			fields, err := d.GetAllFieldsOfUserID(db, field.UserID)
			if err != nil {
				fmt.Printf("database: Не удалось получить запись пользователя: %s \n", field.UserID)
			} else {
				fmt.Printf("database: Удалось получить запись пользователя: %s \n", field.UserID)
			}
			receive <- fields
		case field := <-deleteChan:
			count, err := d.Delete(db, field.UserID)
			if err != nil {
				fmt.Printf("database: Не удалось удалить запись пользователя: %s \n", field.UserID)
			} else {
				fmt.Printf("database: Удалось удалить %d записей пользователя: %s \n", count, field.UserID)
			}
		case field := <-deleteOfItemChan:
			count, err := d.DeleteOneItem(db, field.UserID, field.DeleteNumber)
			if err != nil {
				fmt.Printf("database: Не удалось удалить запись пользователя: %s \n", field.UserID)
			} else {
				fmt.Printf("database: Удалось удалить %d записей пользователя: %s \n", count, field.UserID)
			}
		}
	}
}

// Функция вставки данных
func (d DataBase) Save(db *sql.DB, user_id string, user_url string) (int, error) {
	sqlInsert := `INSERT INTO sources (user_id, user_url) VALUES ($1, $2) RETURNING id`
	id := 0
	err := db.QueryRow(sqlInsert, user_id, user_url).Scan(&id)
	if err != nil {
		return 0, err
	}
	// fmt.Println("Новый ID:", id)
	return id, nil
}

// Функция запроса данных пользователя из таблицы sources
func (d DataBase) GetAllFieldsOfUserID(db *sql.DB, user_id string) ([]Field, error) {
	// SQL-запрос для выборки всех пользователей
	query := "SELECT id, user_id, user_url FROM sources WHERE user_id = $1"

	// Выполняем запрос
	rows, err := db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Создаем срез для хранения данных
	var fields []Field

	// Проходимся по строкам результата
	for rows.Next() {
		var field Field
		if err := rows.Scan(&field.ID, &field.UserID, &field.UserURL); err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}

	// Проверяем на ошибки после обработки всех строк
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return fields, nil
}

// Функция удаления данных
func (d DataBase) Delete(db *sql.DB, user_id string) (int64, error) {
	// SQL-запрос для удалиеня строки по user_id
	query := "DELETE FROM sources WHERE user_id = $1"

	// Выполняем запрос
	result, err := db.Exec(query, user_id)
	if err != nil {
		return 0, err
	}
	// Получаем количество удалённых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// Функция удаления определённого элемента пользователя
func (d DataBase) DeleteOneItem(db *sql.DB, user_id string, number int) (int64, error) {
	responseSlice, err := d.GetAllFieldsOfUserID(db, user_id)
	if err != nil {
		return 0, err
	}

	if number > len(responseSlice) || number < 1 {
		return 0, fmt.Errorf("Нет элемента с таким индексом")
	}
	deleteElement := responseSlice[number-1]

	// SQL-запрос для удалиеня строки по user_id
	query := "DELETE FROM sources WHERE user_id = $1 AND id = $2"

	// Выполняем запрос
	result, err := db.Exec(query, user_id, strconv.FormatInt(int64(deleteElement.ID), 10))
	if err != nil {
		return 0, err
	}
	// Получаем количество удалённых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// TEST Функция для еденичного запроса данных
func (d DataBase) GetOneField(db *sql.DB, user_id string) (string, string, error) {
	var userID string
	var userURL string

	query := "SELECT user_id, user_url FROM sources WHERE user_id = $1"
	rows, err := db.Query(query, user_id)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&userID, &userURL)
		if err != nil {
			return "", "", err
		}
		// fmt.Printf("userID %s, userURL: %s\n", userID, userURL)
	}
	return userID, userURL, nil
}

// TEST Функция запроса данных всех пользователей из таблицы sources
func (d DataBase) GetAllFields(db *sql.DB) ([]Field, error) {
	// SQL-запрос для выборки всех пользователей
	query := "SELECT id, user_id, user_url FROM sources"

	// Выполняем запрос
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Создаем срез для хранения данных
	var fields []Field

	// Проходимся по строкам результата
	for rows.Next() {
		var field Field
		if err := rows.Scan(&field.ID, &field.UserID, &field.UserURL); err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}

	// Проверяем на ошибки после обработки всех строк
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return fields, nil
}

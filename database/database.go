package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Field struct {
	ID      int
	UserID  string
	UserURL string
}

func main() {
	connStr := "user=postgres dbname=LinkKeeper password=postgres host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Вставка данных
	id, _ := Save(db, "88888", "https://axaxaxax.com")
	fmt.Println(id)

	// userID, userURL, err := Get(db, "1234")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(userID, userURL)

	fields2, err := GetAllFields(db)
	if err != nil {
		panic(err)
	}
	for _, fieild := range fields2 {
		fmt.Println(fieild.ID, fieild.UserID, fieild.UserURL)
	}


}

// Функция вставки данных
func Save(db *sql.DB, user_id string, user_url string) (int, error) {
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
func GetAllFieldsOfUser(db *sql.DB, user_id string) ([]Field, error) {
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
func Delete(db *sql.DB, user_id string) (int64, error) {
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

// TEST Функция для еденичного запроса данных
func Get(db *sql.DB, user_id string) (string, string, error) {
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
func GetAllFields(db *sql.DB) ([]Field, error) {
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

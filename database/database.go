package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=postgres dbname=LinkKeeper password=postgres host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Вставка данных
	sqlInsert := `INSERT INTO sources (user_name, user_url) VALUES ($1, $2) RETURNING id`
	id := 0
	err = db.QueryRow(sqlInsert, "Фёдор", "https://google.com").Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("Новый ID:", id)

	// Запрос данных
	// var name string
	// var age int
	// rows, err := db.Query("SELECT name, age FROM users WHERE age >= $1", 18)
	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	err := rows.Scan(&name, &age)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Printf("Имя: %s, Возраст: %d\n", name, age)
	// }
}

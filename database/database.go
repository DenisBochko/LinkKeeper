package database


import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func main() {
    connStr := "user=username dbname=mydb password=mypassword host=localhost sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Вставка данных
    sqlInsert := `INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id`
    id := 0
    err = db.QueryRow(sqlInsert, "Иван", 30).Scan(&id)
    if err != nil {
        panic(err)
    }
    fmt.Println("Новый ID:", id)

    // Запрос данных
    var name string
    var age int
    rows, err := db.Query("SELECT name, age FROM users WHERE age >= $1", 18)
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(&name, &age)
        if err != nil {
            panic(err)
        }
        fmt.Printf("Имя: %s, Возраст: %d\n", name, age)
    }
}
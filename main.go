package main

import (
    "fmt"
    "auth-service/internal/database"
)

func main() {
    dsn := "host=localhost port=5432 user=postgres password=secret dbname=auth sslmode=disable"
    db.InitDB(dsn)

    fmt.Println("Auth Service is running...")
}

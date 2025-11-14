package main

import (
    "fmt"
    "log"
    "planica_bi/database/config"
)

func main() {
    db, err := config.Connect("development")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    repo := config.NewRepository(db)
    version, err := repo.GetDatabaseVersion()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("✅ Подключение успешно! Версия MySQL: %s\n", version)
}


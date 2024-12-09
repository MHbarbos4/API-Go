package main

import (
    "database/sql"
    "log"
    _ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() {
    var err error
    db, err = sql.Open("sqlite", "./data.db")
    if err != nil {
        log.Fatal(err)
    }

    createTableQuery := `
    CREATE TABLE IF NOT EXISTS items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        price REAL NOT NULL
    );`
    if _, err := db.Exec(createTableQuery); err != nil {
        log.Fatal(err)
    }
}
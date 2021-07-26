package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var wg sync.WaitGroup
var db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/logs_database")

func main() {
	start := time.Now()
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(151)

	r, err := ioutil.ReadFile("data/log.txt")
	if err != nil {
		log.Fatal(err)
	}

	fullFile := string(r)
	fileLines := strings.Split(fullFile, "\n")
	for _, line := range fileLines {
		wg.Add(1)
		go CheckLine(line)
	}
	wg.Wait()
	fmt.Println("Finalizou", time.Since(start))
}

func CheckLine(line string) {
	defer wg.Done()

	if strings.Contains(line, "pam_unix(sshd:auth): check pass; user unknown") {
		insert, err := db.Query("INSERT INTO logs (log) VALUES ('" + line + "')")
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
	}
}

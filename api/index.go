package handler

import (
	sql "database/sql"
	fmt "fmt"
	log "log"
	http "net/http"
	url "net/url"
	os "os"
	strings "strings"

	pq "github.com/lib/pq"
)

func Handler(writer http.ResponseWriter, request *http.Request) {
	conn := os.Getenv("DATABASE_URL")
	db, _ := sql.Open("postgres", conn)
	searchPhrase, err := url.QueryUnescape(request.URL.Query().Get("character"))
	if err != nil {
		log.Fatal(err)
	}

	query := strings.Split(strings.ToLower(searchPhrase), " ")

	row := db.QueryRow("SELECT accounts.username, characters.name FROM characters JOIN accounts ON characters.account_id = accounts.id WHERE $1 <@ characters.search_phrases GROUP BY characters.name, accounts.username LIMIT 1", pq.Array(query))
	var username string
	var name string
	if err = row.Scan(&username, &name); err != nil {
		log.Fatal(err)
	}

	writer.Header().Set("Cache-Control", "s-maxage=1, stale-while-revalidate")
	if username == "" || name == "" {
		fmt.Fprint(writer, "character not found")
	} else {
		fmt.Fprintf(writer, "%s is played by twitch.tv/%s", name, username)
	}
}

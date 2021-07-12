package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/lib/pq"
)

type account struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type character struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Keywords []string `json:"keywords"`
	Player   account  `json:"player"`
}

func Handler(writer http.ResponseWriter, request *http.Request) {
	conn := os.Getenv("DATABASE_URL")
	db, _ := sql.Open("postgres", conn)
	searchPhrase := request.URL.Query().Get("id")
	query := strings.Split(strings.ToLower(searchPhrase), "_")

	row := db.QueryRow("SELECT characters.name, accounts.username, accounts.id, character.search_phrases, characters.id FROM characters JOIN accounts ON characters.account_id = accounts.id WHERE $1 <@ characters.search_phrases GROUP BY characters.name, accounts.username, accounts.id, character.search_phrases, characters.id LIMIT 1", pq.Array(query))

	var Id string
	var Name string
	var Keywords []string
	var pId string
	var username string

	row.Scan(&Name, &username, &pId, pq.Array(&Keywords), &Id)

	c := character{
		Id,
		Name,
		Keywords,
		account{
			pId,
			username,
		},
	}
	data, _ := json.Marshal(c)
	writer.Header().Set("Cache-Control", "s-maxage=1, stale-while-revalidate")
	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprint(writer, string(data))
}

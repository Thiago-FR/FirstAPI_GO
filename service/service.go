package service

import (
	"crud/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"name"`
	Email string `json:"email"`
}

// CreateUser Insere usuário no DB
func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição"))
		return
	}

	var user user
	if erro := json.Unmarshal(requestBody, &user); erro != nil {
		w.Write([]byte("Erro ao converter o usuário pra struct"))
		return
	}

	db, erro := db.Connect()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("INSERT INTO user (name, email) VALUES (?, ?)")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	insert, erro := statement.Exec(user.Nome, user.Email)
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro ao executar o statement!"))
		return
	}

	idInsert, erro := insert.LastInsertId()
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro ao obter o ID inserido!"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuário inserido com sucesso! Id: %d", idInsert)))
}

// FindUser Busca usuário no DB
func FindAllUser(w http.ResponseWriter, r *http.Request) {
	db, erro := db.Connect()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}
	defer db.Close()

	rows, erro := db.Query("SELECT * FROM user")
	if erro != nil {
		w.Write([]byte("Erro ao buscar usuários"))
		return
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var user user

		if erro := rows.Scan(&user.ID, &user.Nome, &user.Email); erro != nil {
			w.Write([]byte("Erro ao escanear o usuário"))
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(users); erro != nil {
		w.Write([]byte("Erro ao converter os usuários pra JSON"))
		return
	}
}

// FindUser Busca usuário no DB
func FindOneUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, erro := strconv.ParseUint(params["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter o parâmetro para inteiro"))
		return
	}

	db, erro := db.Connect()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}
	defer db.Close()

	row, erro := db.Query("SELECT * FROM user WHERE id = ?", ID)
	if erro != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Erro ao buscar o usuário!"))
		return
	}

	var user user
	if row.Next() {
		if erro := row.Scan(&user.ID, &user.Nome, &user.Email); erro != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Erro ao escanear o usuário"))
			return
		}
	}

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Usuário não encontrado!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(user); erro != nil {
		w.Write([]byte("Erro ao converter o usuário pra JSON"))
		return
	}
}

// UpdateUser Atualizar usuário
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, erro := strconv.ParseUint(params["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter o parâmetro para inteiro"))
		return
	}

	requestBody, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Erro ao ler o corpo da requisição!"))
	}

	var user user
	if erro := json.Unmarshal(requestBody, &user); erro != nil {
		w.Write([]byte("Erro ao converter o usuário pra struct"))
		return
	}

	db, erro := db.Connect()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("UPDATE user SET name = ?, email = ? WHERE id = ?")
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(user.Nome, user.Email, ID); erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro ao atualizar o usuário!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser Deleta usuário do DB
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, erro := strconv.ParseUint(params["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter o parâmetro para inteiro"))
		return
	}

	db, erro := db.Connect()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("DELETE FROM user WHERE id = ?")
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(ID); erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro ao atualizar o usuário!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

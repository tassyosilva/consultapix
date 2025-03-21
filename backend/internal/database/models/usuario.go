// internal/database/models/usuario.go
package models

import "time"

type Usuario struct {
	ID        int       `json:"id" db:"id"`
	Nome      string    `json:"nome" db:"nome"`
	CPF       string    `json:"cpf" db:"cpf"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // O campo senha nunca deve ser retornado no JSON
	Lotacao   string    `json:"lotacao" db:"lotacao"`
	Matricula string    `json:"matricula" db:"matricula"`
	Admin     bool      `json:"admin" db:"admin"`
}
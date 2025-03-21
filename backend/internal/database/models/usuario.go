package models

type Usuario struct {
	ID        int    `json:"id" db:"id"`
	Nome      string `json:"nome" db:"nome"`
	CPF       string `json:"cpf" db:"cpf"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"-" db:"password"`
	Lotacao   string `json:"lotacao" db:"lotacao"`
	Matricula string `json:"matricula" db:"matricula"`
	Admin     bool   `json:"admin" db:"admin"`
}
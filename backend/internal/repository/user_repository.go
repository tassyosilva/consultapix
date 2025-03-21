package repository

import (
	"database/sql"
	"errors"

	"github.com/tassyosilva/consultapix/internal/database"
	"github.com/tassyosilva/consultapix/internal/database/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		DB: database.GetDB(),
	}
}

func (r *UserRepository) FindByEmail(email string) (*models.Usuario, error) {
	var user models.Usuario
	query := `SELECT id, nome, cpf, email, password, lotacao, matricula, admin FROM usuario WHERE email = $1`
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Nome, &user.CPF, &user.Email, &user.Password, &user.Lotacao, &user.Matricula, &user.Admin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.Usuario) (int, error) {
	// Verificar se já existe usuário com o email, cpf ou matrícula
	var count int
	query := `SELECT COUNT(*) FROM usuario WHERE email = $1 OR cpf = $2 OR matricula = $3`
	err := r.DB.QueryRow(query, user.Email, user.CPF, user.Matricula).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("usuário já cadastrado")
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Password = string(hashedPassword)

	// Inserir usuário
	insertQuery := `
		INSERT INTO usuario (nome, cpf, email, password, lotacao, matricula, admin)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id int
	err = r.DB.QueryRow(
		insertQuery, 
		user.Nome, user.CPF, user.Email, user.Password, user.Lotacao, user.Matricula, user.Admin,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) Update(user *models.Usuario) error {
	// Hash da senha se fornecida
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
		
		updateQuery := `
			UPDATE usuario
			SET nome = $1, cpf = $2, email = $3, password = $4, lotacao = $5, matricula = $6, admin = $7
			WHERE id = $8
		`
		_, err = r.DB.Exec(
			updateQuery,
			user.Nome, user.CPF, user.Email, user.Password, user.Lotacao, user.Matricula, user.Admin, user.ID,
		)
		return err
	}
	
	// Atualização sem senha
	updateQuery := `
		UPDATE usuario
		SET nome = $1, cpf = $2, email = $3, lotacao = $4, matricula = $5, admin = $6
		WHERE id = $7
	`
	_, err := r.DB.Exec(
		updateQuery,
		user.Nome, user.CPF, user.Email, user.Lotacao, user.Matricula, user.Admin, user.ID,
	)
	return err
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM usuario WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *UserRepository) GetAll() ([]models.Usuario, error) {
	query := `SELECT id, nome, cpf, email, '', lotacao, matricula, admin FROM usuario ORDER BY nome ASC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.Usuario
	for rows.Next() {
		var user models.Usuario
		err := rows.Scan(
			&user.ID, &user.Nome, &user.CPF, &user.Email, &user.Password, &user.Lotacao, &user.Matricula, &user.Admin,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
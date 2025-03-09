package user

import (
	"database/sql"
	"fmt"

	"github.com/sikozonpc/ecom/types"
)

// Define a estrutura 'Store' que vai encapsular a conexão com o banco de dados.
type Store struct {
	db *sql.DB // Campo que armazena a conexão com o banco de dados SQL.
}

// Função para criar uma nova instância de 'Store' passando uma conexão de banco de dados.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db} // Retorna um ponteiro para uma nova instância de 'Store' com a conexão de banco de dados.
}

func (s *Store) CreateUser(user types.User) error {
	// Executa uma consulta SQL para inserir um novo usuário na tabela 'users'.
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password)

	// Se houver um erro na execução da consulta, retorna o erro.
	if err != nil {
		return err
	}

	// Caso contrário, retorna nil, indicando que o usuário foi criado com sucesso.
	return nil
}

// Função para buscar um usuário no banco de dados pelo seu e-mail.
func (s *Store) GetUserByEmail(email string) (*types.User, error) {

	// Executa uma consulta SQL para buscar um usuário com o e-mail fornecido.
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err // Se ocorrer um erro ao executar a consulta, retorna o erro.
	}

	// Cria uma nova instância de 'User' para armazenar os dados recuperados.
	u := new(types.User)

	// Itera sobre as linhas retornadas pela consulta.
	for rows.Next() {

		// Para cada linha, tenta mapear os dados para a estrutura 'User' usando a função scanRowsIntoUser.
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err // Se ocorrer um erro ao mapear os dados, retorna o erro.
		}
	}

	// Se o ID do usuário for 0, significa que o usuário não foi encontrado.
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found") // Retorna um erro informando que o usuário não foi encontrado.
	}

	// Retorna o usuário encontrado.
	return u, nil
}

// Função para buscar um usuário no banco de dados pelo seu ID.
func (s *Store) GetUserByID(id int) (*types.User, error) {

	// Executa uma consulta SQL para buscar um usuário com o ID fornecido.
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err // Se ocorrer um erro ao executar a consulta, retorna o erro.
	}

	// Cria uma nova instância de 'User' para armazenar os dados recuperados.
	u := new(types.User)

	// Itera sobre as linhas retornadas pela consulta.
	for rows.Next() {

		// Para cada linha, tenta mapear os dados para a estrutura 'User' usando a função scanRowsIntoUser.
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err // Se ocorrer um erro ao mapear os dados, retorna o erro.
		}
	}
	// Se o ID do usuário for 0, significa que o usuário não foi encontrado.
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	// Retorna o usuário encontrado.
	return u, nil
}

// Função auxiliar para mapear os dados de uma linha do banco de dados para uma estrutura 'User'.
func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {

	// Cria uma nova instância de 'User' para armazenar os dados mapeados.
	user := new(types.User)

	// Tenta mapear os dados da linha para os campos da estrutura 'User'.
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Retorna a instância de 'User' com os dados mapeados.
	return user, nil
}

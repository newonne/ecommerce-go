package auth

import (
	"golang.org/x/crypto/bcrypt" // Importa o pacote bcrypt, usado para hashing e verificação de senhas.
)

// HashPassword recebe uma senha em texto plano e retorna seu hash ou um erro.
func HashPassword(password string) (string, error) {
	// Gera o hash da senha fornecida usando o custo padrão do bcrypt.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// Retorna uma string vazia e o erro, caso ocorra alguma falha durante a geração do hash.
		return "", err
	}
	// Converte o hash (um slice de bytes) em uma string e o retorna junto com um valor nil indicando sucesso.
	return string(hash), nil
}

// ComparePasswords compara uma senha em texto plano com um hash armazenado e retorna um booleano indicando se coincidem.
func ComparePasswords(hashed string, plain []byte) bool {
	// Compara o hash fornecido com a senha em texto plano. Se forem equivalentes, retorna nil.
	err := bcrypt.CompareHashAndPassword([]byte(hashed), plain)
	// Retorna true se não houver erro (as senhas coincidem) ou false caso contrário.
	return err == nil
}

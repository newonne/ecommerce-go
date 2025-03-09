package auth

import (
	"testing"
)

// TestCreateJWT é uma função de teste para verificar a funcionalidade da criação de JWT.

func TestCreateJWT(t *testing.T) {
	// Define o segredo usado para assinar o JWT como uma sequência de bytes.
	secret := []byte("secret")

	// Chama a função CreateJWT (supostamente definida no mesmo pacote) com o segredo e um identificador de usuário.
	token, err := CreateJWT(secret, 1)
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	if token == "" {
		t.Error("expected token to be not empty")
	}
}

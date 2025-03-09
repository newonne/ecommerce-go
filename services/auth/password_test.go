package auth

import (
	"testing"
)

// Função de teste para verificar a funcionalidade de hashing de senhas.
func TestHashPassword(t *testing.T) {
	// Tenta gerar um hash para a senha "password". A função HashPassword retorna o hash e um possível erro.
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if hash == "" {
		t.Error("expected hash to be not empty")
	}

	if hash == "password" {
		t.Error("expected hash to be different from password")
	}
}

// Função de teste para verificar a comparação de senhas com hashes.
func TestComparePasswords(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if !ComparePasswords(hash, []byte("password")) {
		t.Errorf("expected password to match hash")
	}
	if ComparePasswords(hash, []byte("notpassword")) {
		t.Errorf("expected password to not match hash")
	}
}

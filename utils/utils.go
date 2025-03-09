package utils

import (
	"encoding/json" // Pacote para codificar e decodificar JSON
	"fmt"           // Pacote para formatação de strings e erros
	"net/http"      // Pacote para manipulação de requisições e respostas HTTP

	"github.com/go-playground/validator/v10" // Pacote para validação de dados (não utilizado diretamente neste código)
)

var Validate = validator.New()

// Função que escreve uma resposta HTTP em formato JSON
func WriteJSON(w http.ResponseWriter, status int, v any) error {

	// Adiciona o cabeçalho Content-Type como "application/json" para informar que a resposta será no formato JSON
	w.Header().Add("Content-Type", "application/json")

	// Define o código de status HTTP da resposta
	w.WriteHeader(status)

	// Codifica o valor 'v' como JSON e escreve na resposta HTTP
	return json.NewEncoder(w).Encode(v)
}

// Função que escreve uma resposta de erro em formato JSON
func WriteError(w http.ResponseWriter, status int, err error) {

	// Chama a função WriteJSON passando um mapa com a chave "error" e a mensagem de erro
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

// Função que faz o parse (leitura) do corpo da requisição HTTP e o decodifica como JSON
func ParseJSON(r *http.Request, v any) error {

	// Verifica se o corpo da requisição é nil, ou seja, não foi enviado nenhum corpo na requisição
	if r.Body == nil {

		// Retorna um erro caso o corpo esteja ausente
		return fmt.Errorf("missing request body")
	}

	// Decodifica o corpo da requisição JSON e preenche o valor de 'v' com os dados
	return json.NewDecoder(r.Body).Decode(v)
}

// Função que extrai o token de autenticação da requisição HTTP
func GetTokenFromRequest(r *http.Request) string {

	// Tenta obter o token do cabeçalho "Authorization"
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	// Se o token estiver presente no cabeçalho, retorna ele
	if tokenAuth != "" {
		return tokenAuth
	}

	// Se o token não estiver no cabeçalho, mas estiver na query string, retorna ele
	if tokenQuery != "" {
		return tokenQuery
	}

	// Caso nenhum dos dois tokens seja encontrado, retorna uma string vazia
	return ""
}

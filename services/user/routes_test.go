package user

import (
	"net/http"          // Pacote para manipulação de requisições e respostas HTTP.
	"net/http/httptest" // Pacote para criar testes de servidores HTTP.
	"testing"           // Pacote para escrever testes unitários.

	"github.com/gorilla/mux"          // Pacote de roteamento HTTP utilizado para manipulação de rotas.
	"github.com/sikozonpc/ecom/types" // Importa o pacote 'types' que define o tipo 'User'.
)

func TestUserServiceHandlers(t *testing.T) {
	// Cria um "mock" da camada de armazenamento de usuários (mockUserStore) que simula operações no banco de dados.
	userStore := &mockUserStore{}
	handler := NewHandler(userStore) // Cria um novo manipulador (handler) passando o "mock" como a camada de persistência.

	t.Run("should fail if the user ID is not a number", func(t *testing.T) {

		// Cria uma nova requisição GET para o endpoint "/user/abc" (ID inválido).
		req, err := http.NewRequest(http.MethodGet, "/user/abc", nil)
		if err != nil {
			t.Fatal(err) // Se ocorrer um erro ao criar a requisição, o teste falha.
		}
		// Cria um novo gravador de resposta (httptest.NewRecorder) para capturar a resposta da requisição.
		rr := httptest.NewRecorder()

		// Cria um novo roteador e registra o manipulador para o endpoint "/user/{userID}".
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}", handler.handleGetUser).Methods(http.MethodGet)

		// Executa a requisição HTTP usando o roteador e grava a resposta no gravador.
		router.ServeHTTP(rr, req)

		// Verifica se o código de status retornado é 400 (Bad Request), pois o ID não é válido.
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should handle get user by ID", func(t *testing.T) {

		// Cria uma nova requisição GET para o endpoint "/user/42" (ID válido).
		req, err := http.NewRequest(http.MethodGet, "/user/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}", handler.handleGetUser).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

type mockUserStore struct{}

func (m *mockUserStore) UpdateUser(u types.User) error {
	return nil
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return &types.User{}, nil
}

func (m *mockUserStore) CreateUser(u types.User) error {
	return nil
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return &types.User{}, nil
}

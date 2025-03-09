package user

import (
	"fmt"      // Pacote para formatação de strings e manipulação de erros.
	"net/http" // Pacote para manipulação de requisições e respostas HTTP.
	"strconv"  // Pacote para conversão de tipos, usado para converter strings em números.

	"github.com/go-playground/validator/v10"  // Pacote para validação de structs em Go.
	"github.com/gorilla/mux"                  // Pacote de roteamento HTTP, usado para definir rotas na aplicação.
	"github.com/sikozonpc/ecom/configs"       // Pacote de configuração, usado para acessar variáveis de ambiente e configurações da aplicação.
	"github.com/sikozonpc/ecom/services/auth" // Pacote de autenticação, contendo funções para criptografia de senhas e geração de JWTs.
	"github.com/sikozonpc/ecom/types"         // Pacote que contém os tipos usados no sistema, como o tipo User e os payloads de login e registro.
	"github.com/sikozonpc/ecom/utils"         // Pacote utilitário, que contém funções auxiliares para manipulação de JSON, validação, e manipulação de erros.
)

// Handler é a estrutura que irá conter os manipuladores de rotas relacionados a usuários.
type Handler struct {
	store types.UserStore // A estrutura Handler contém um campo 'store', que é uma interface para acessar dados de usuários (por exemplo, banco de dados).
}

// NewHandler cria uma nova instância de Handler, passando um objeto 'store' que implementa a interface 'UserStore'.
func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store} // Retorna um ponteiro para um novo Handler com a store fornecida.
}

// RegisterRoutes define as rotas que o servidor HTTP deve reconhecer e as associa aos respectivos manipuladores.
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Registra as rotas principais de login e registro de usuário, associando-as aos manipuladores correspondentes.
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")

	// Registra a rota de obtenção de informações de um usuário, exigindo autenticação JWT para o acesso.
	// A função 'auth.WithJWTAuth' é um middleware que valida o token JWT antes de chamar o manipulador real.
	router.HandleFunc("/users/{userID}", auth.WithJWTAuth(h.handleGetUser, h.store)).Methods(http.MethodGet)
}

// handleLogin é o manipulador que trata a requisição de login de um usuário.
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Cria uma variável para armazenar os dados do payload de login (e-mail e senha).
	var user types.LoginUserPayload

	// Faz o parsing do corpo da requisição para o tipo LoginUserPayload. Caso haja erro, responde com erro 400.
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Valida os dados usando o pacote validator. Caso haja erro, responde com erro 400.
	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)                                            // Converte o erro para um tipo de erro de validação.                                          // Converte o erro para um tipo de erro de validação.
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors)) // Responde com um erro de validação.
		return
	}
	// Tenta buscar o usuário no banco de dados pelo e-mail.
	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {

		// Se o usuário não for encontrado, responde com um erro de login.
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}
	// Compara a senha fornecida com a senha armazenada no banco de dados.
	if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		// Se a senha não for correta, responde com erro 400.
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}
	// Cria um token JWT para o usuário.
	secret := []byte(configs.Envs.JWTSecret)   // Obtém o segredo para o JWT da configuração.
	token, err := auth.CreateJWT(secret, u.ID) // Gera o token JWT com o ID do usuário.
	if err != nil {

		// Se houver erro ao criar o token, responde com erro 500.
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// Responde com o token JWT gerado.
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

// handleRegister é o manipulador que lida com o registro de um novo usuário.
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUserPayload // Variável para armazenar os dados do usuário a ser registrado.

	// Faz o parsing do corpo da requisição para o tipo RegisterUserPayload. Caso haja erro, responde com erro 400.
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Valida os dados usando o pacote validator. Caso haja erro, responde com erro 400.
	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// Verifica se já existe um usuário com o e-mail fornecido.
	_, err := h.store.GetUserByEmail(user.Email)
	if err == nil {

		// Se o usuário já existe, responde com erro 400.
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", user.Email))
		return
	}

	// Cria uma versão criptografada da senha.
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {

		// Se houver erro ao criptografar a senha, responde com erro 500.
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Cria o novo usuário no banco de dados com os dados fornecidos.
	err = h.store.CreateUser(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	})
	if err != nil {

		// Se houver erro ao criar o usuário, responde com erro 500.
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Responde com sucesso (código 201 Created) quando o usuário é registrado corretamente.
	utils.WriteJSON(w, http.StatusCreated, nil)
}

// handleGetUser é o manipulador que lida com a requisição para obter informações de um usuário específico.
func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {

	// Obtém o valor do parâmetro "userID" da URL.
	vars := mux.Vars(r)
	str, ok := vars["userID"]
	if !ok {
		// Se o parâmetro "userID" não for fornecido, responde com erro 400.
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	// Converte o ID do usuário de string para inteiro.
	userID, err := strconv.Atoi(str)
	if err != nil {
		// Se o ID não for um número válido, responde com erro 400.
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}
	// Tenta buscar o usuário no banco de dados pelo ID fornecido.
	user, err := h.store.GetUserByID(userID)
	if err != nil {
		// Se houver erro ao buscar o usuário, responde com erro 500.
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// Responde com os dados do usuário (status 200 OK).
	utils.WriteJSON(w, http.StatusOK, user)
}

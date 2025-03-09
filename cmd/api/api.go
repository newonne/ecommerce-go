//Esse arquivo representa o núcleo de uma API backend modular e escalável, adequada para uma aplicação de e-commerce. Ele é bem estruturado,
//com boas práticas de separação de responsabilidades, o que facilita a manutenção e a adição de novas funcionalidades.

// https://chatgpt.com/c/676c94f7-1ddc-800a-9e78-171cbe2d6e61
package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sikozonpc/ecom/services/cart"
	"github.com/sikozonpc/ecom/services/order"
	"github.com/sikozonpc/ecom/services/product"
	"github.com/sikozonpc/ecom/services/user"
)

// APIServer é a estrutura principal que representa o servidor da API.
// Ela contém o endereço do servidor (addr) e a conexão com o banco de dados (db).
type APIServer struct {
	addr string
	db   *sql.DB
}

// NewAPIServer cria e retorna uma nova instância de APIServer.
// Parametros:
// - addr: endereço no formato "host:porta" onde o servidor será iniciado.
// - db: conexão ao banco de dados.
// Retorna: ponteiro para uma instância de APIServer.
func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

// Run inicializa e executa o servidor da API.
// Configura as rotas, inicializa os handlers e inicia o servidor HTTP.
// Retorna um erro se o servidor falhar ao iniciar.
func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// Configuração do serviço de usuários.
	userStore := user.NewStore(s.db)          // Cria a camada de armazenamento para usuários.
	userHandler := user.NewHandler(userStore) // Cria o handler responsável por gerenciar rotas de usuários.
	userHandler.RegisterRoutes(subrouter)     // Registra as rotas relacionadas a usuários no subroteador.

	// Configuração do serviço de produtos.
	productStore := product.NewStore(s.db)                        // Cria a camada de armazenamento para produtos.
	productHandler := product.NewHandler(productStore, userStore) // Cria o handler para gerenciar produtos, integrando usuários.
	productHandler.RegisterRoutes(subrouter)                      // Registra as rotas de produtos no subroteador.

	// Configuração do serviço de pedidos.
	orderStore := order.NewStore(s.db) // Cria a camada de armazenamento para pedidos.

	// Configuração do serviço de carrinho de compras.
	cartHandler := cart.NewHandler(productStore, orderStore, userStore) // Cria o handler para carrinhos.
	cartHandler.RegisterRoutes(subrouter)                               // Registra as rotas de carrinhos no subroteador.

	// Serve static files
	// Qualquer rota que não coincida com as anteriores servirá arquivos da pasta "static".
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	log.Println("Listening on", s.addr)

	// Inicia o servidor HTTP e associa o roteador configurado.
	return http.ListenAndServe(s.addr, router)
}

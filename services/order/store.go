package order

import (
	"database/sql"

	"github.com/sikozonpc/ecom/types"
)

//****** Esse código, portanto, faz a inserção de pedidos e itens de pedidos em um banco de dados,
//retornando o ID do pedido e tratando erros quando necessário.**/

// Define a estrutura 'Store', que representa o armazenamento de dados (banco de dados) com um campo 'db' do tipo *sql.DB.
type Store struct {
	db *sql.DB // 'db' é uma referência para uma conexão com o banco de dados.
}

// Função que cria uma nova instância de 'Store' com uma conexão de banco de dados fornecida.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db} // Retorna um ponteiro para uma nova instância de 'Store' com o banco de dados associado.
}

// Método 'CreateOrder' da estrutura 'Store', que cria um novo pedido no banco de dados.
func (s *Store) CreateOrder(order types.Order) (int, error) {
	// Executa um comando SQL para inserir um novo pedido na tabela 'orders'.
	// Os valores são passados como parâmetros, substituindo os pontos de interrogação.
	res, err := s.db.Exec("INSERT INTO orders (userId, total, status, address) VALUES (?, ?, ?, ?)", order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		return 0, err
	}

	// Recupera o ID do último registro inserido (o ID do pedido recém-criado).
	id, err := res.LastInsertId()
	if err != nil {
		// Se houver um erro ao obter o ID, retorna 0 e o erro.
		return 0, err
	}

	// Retorna o ID do pedido como um inteiro, junto com um valor de erro nulo (sem erro).
	return int(id), nil
}

// Método 'CreateOrderItem' da estrutura 'Store', que cria um item de pedido no banco de dados.
func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO order_items (orderId, productId, quantity, price) VALUES (?, ?, ?, ?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	return err
}

package auth

import (
	"context"  // Importa o pacote 'context', utilizado para propagação de informações entre as requisições.
	"fmt"      // Importa o pacote 'fmt', usado para formatar mensagens e erros.
	"log"      // Importa o pacote 'log', usado para logar mensagens de erro e outras informações.
	"net/http" // Importa o pacote 'http', que oferece funcionalidades para manipulação de requisições HTTP.
	"strconv"  // Importa o pacote 'strconv', usado para converter valores entre tipos de dados, como string para int.
	"time"     // Importa o pacote 'time', utilizado para manipulação de datas e horários.

	"github.com/golang-jwt/jwt/v5"      // Importa a biblioteca para trabalhar com JSON Web Tokens (JWTs).
	"github.com/sikozonpc/ecom/configs" // Importa o pacote 'configs', onde estão configuradas variáveis de ambiente.
	"github.com/sikozonpc/ecom/types"   // Importa o pacote 'types', onde estão definidos tipos como o modelo 'UserStore'.
	"github.com/sikozonpc/ecom/utils"   // Importa o pacote 'utils', que contém funções auxiliares como a manipulação de erros e a extração de tokens.
)

// Define um tipo 'contextKey' como um alias de 'string', utilizado para evitar conflitos com chaves de contexto.
type contextKey string

// Declara uma constante 'UserKey', que será usada como chave para armazenar o 'userID' no contexto.
const UserKey contextKey = "userID"

// Função 'WithJWTAuth' que adiciona autenticação JWT à rota.
// Ela recebe uma função de manipulação de requisição (handlerFunc)
// e um repositório de usuários (store).
func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Extrai o token JWT da requisição. A função 'utils.GetTokenFromRequest' é responsável por verificar
		// o cabeçalho ou a query da requisição.
		tokenString := utils.GetTokenFromRequest(r)

		// Valida o token JWT utilizando a função 'validateJWT'. A função retorna o token validado ou um erro.
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		// Extrai as "claims" do token (informações armazenadas dentro do JWT). 'jwt.MapClaims
		//' é uma estrutura que facilita o acesso a essas informações.
		claims := token.Claims.(jwt.MapClaims)

		// Extrai o 'userID' das claims, que foi armazenado como uma string no JWT.
		str := claims["userID"].(string)

		// Converte o 'userID' de string para int. O 'userID' foi codificado como string no JWT, mas é tratado como inteiro aqui.
		userID, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("failed to convert userID to int: %v", err)
			permissionDenied(w)
			return
		}

		// Busca o usuário no banco de dados usando o 'userID'. A função 'store.GetUserByID' é chamada para isso.
		u, err := store.GetUserByID(userID)
		if err != nil { // Se não encontrar o usuário, loga o erro e retorna permissão negada.
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

		// Cria um novo contexto, armazenando o 'userID' no contexto da requisição. O contexto será propagado para as próximas etapas.
		ctx := r.Context()
		// Usa a chave 'UserKey' para associar o 'userID' ao contexto. Esse valor estará disponível em qualquer parte do código onde o contexto for acessado.
		ctx = context.WithValue(ctx, UserKey, u.ID)
		// Atualiza a requisição (r) com o novo contexto que contém o 'userID'.
		r = r.WithContext(ctx)

		// Chama a função original de manipulação de requisição (handlerFunc), que agora pode acessar o 'userID' do contexto.
		handlerFunc(w, r)
	}
}

// Função para criar um token JWT para um usuário com base no 'userID'. Recebe o 'secret' para assinar o token e o 'userID' do usuário.
func CreateJWT(secret []byte, userID int) (string, error) {
	// Define a expiração do token com base no valor configurado (em segundos) no arquivo de configurações.
	expiration := time.Second * time.Duration(configs.Envs.JWTExpirationInSeconds)

	// Cria um novo JWT com o método de assinatura HS256 e as claims do token, incluindo o 'userID' e o tempo de expiração.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),         // Converte o 'userID' de int para string para armazenar no JWT.
		"expiresAt": time.Now().Add(expiration).Unix(), // Calcula o tempo de expiração do token (em segundos desde a época Unix).
	})

	// Gera o token assinado com o 'secret'. A função 'SignedString' assina o token usando a chave secreta e retorna o token como string.
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err // Retorna um erro se a assinatura falhar.
	}

	return tokenString, err // Retorna o token assinado e o erro (se houver).
}

// Função para validar um token JWT. Recebe o token como string e tenta analisá-lo e verificar sua autenticidade.
func validateJWT(tokenString string) (*jwt.Token, error) {
	// Tenta analisar o token e fornecer a chave de assinatura necessária para validar o token.
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifica se o método de assinatura do token é HMAC, que é o esperado.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// Retorna erro se o método de assinatura for inesperado.
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Retorna a chave secreta para validar a assinatura do token.
		return []byte(configs.Envs.JWTSecret), nil
	})
}

// Função para retornar uma resposta de "permissão negada" (HTTP 403) quando a autenticação falha.
func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

// Função para obter o 'userID' armazenado no contexto da requisição. Retorna o 'userID' ou -1 caso não seja encontrado.
func GetUserIDFromContext(ctx context.Context) int {

	// Tenta acessar o valor 'userID' no contexto. Se não for encontrado, retorna -1.
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1 // Retorna -1 se o 'userID' não for encontrado ou se houver um erro de conversão de tipo.
	}

	return userID
}

package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain/repository"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/logger"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rest_err"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	JWT_SECRET_KEY = "JWT_SECRET_KEY"
)

func NewAuthService(
	userRepository repository.UserRepository,
) AuthService {
	return &authService{userRepository}
}

type authService struct {
	userRepository repository.UserRepository
}

type AuthService interface {
	SignUp(domain.UserDomainInterface) (
		domain.UserDomainInterface,
		*rest_err.RestErr)

	SignIn(
		email string,
		password string,
	) (string, *rest_err.RestErr)

	SignOut() *rest_err.RestErr

	VerifyEmail(
		email, token string,
	) *rest_err.RestErr
}

func (as *authService) SignUp(userDomain domain.UserDomainInterface) (domain.UserDomainInterface, *rest_err.RestErr) {
	user, _ := as.findUserByEmail(userDomain.GetEmail())
	if user != nil {
		return nil, rest_err.NewBadRequestError("Email already registered in another account")
	}

	userDomain.EncryptPassword()
	userDomain.SetCreatedAt(time.Now())
	userDomain.SetUpdatedAt(time.Now())
	user, err := as.userRepository.CreateUser(userDomain)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (as *authService) SignIn(email, password string) (token string, err *rest_err.RestErr) {
	return "", nil
}

func (as *authService) SignOut() *rest_err.RestErr {
	return nil
}

func (as *authService) VerifyEmail(email, token string) *rest_err.RestErr {
	return nil
}

func (as *authService) findUserByEmail(email string) (domain.UserDomainInterface, *rest_err.RestErr) {
	return as.userRepository.FindUserByEmail(email)
}

func (as *authService) GenerateToken(ud domain.UserDomainInterface) (string, *rest_err.RestErr) {
	secret := os.Getenv(JWT_SECRET_KEY)

	claims := jwt.MapClaims{
		"id":         ud.GetID(),
		"email":      ud.GetEmail(),
		"first_name": ud.GetFirstName(),
		"last_name":  ud.GetLastName(),
		"exp":        time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", rest_err.NewInternalServerError(
			fmt.Sprintf("error trying to generate jwt token, err=%s", err.Error()),
		)
	}
	return tokenString, nil
}

func VerifyToken(tokenValue string) (domain.UserDomainInterface, *rest_err.RestErr) {
	secret := os.Getenv(JWT_SECRET_KEY)

	token, err := jwt.Parse(RemoveBearerPrefix(tokenValue), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(secret), nil
		}

		return nil, rest_err.NewBadRequestError("invalid token")
	})
	if err != nil {
		return nil, rest_err.NewUnauthorizedRequestError("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, rest_err.NewUnauthorizedRequestError("invalid token")
	}

	userDomain := domain.NewUserTokenDomain(claims["id"].(string), claims["email"].(string), claims["first_name"].(string), claims["last_name"].(string))

	return userDomain, nil
}

func VerifyTokenMiddleware(c *gin.Context) {
	secret := os.Getenv(JWT_SECRET_KEY)
	tokenValue := RemoveBearerPrefix(c.Request.Header.Get("Authorization"))

	token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(secret), nil
		}

		return nil, rest_err.NewBadRequestError("invalid token")
	})
	if err != nil {
		errRest := rest_err.NewUnauthorizedRequestError("invalid token")
		c.JSON(errRest.Code, errRest)
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		errRest := rest_err.NewUnauthorizedRequestError("invalid token")
		c.JSON(errRest.Code, errRest)
		c.Abort()
		return
	}

	userDomain := domain.NewUserTokenDomain(claims["id"].(string), claims["email"].(string), claims["first_name"].(string), claims["last_name"].(string))

	logger.Info(fmt.Sprintf("User authenticated: %#v", userDomain))
}

func RemoveBearerPrefix(token string) string {
	return strings.TrimPrefix(token, "Bearer ")
}

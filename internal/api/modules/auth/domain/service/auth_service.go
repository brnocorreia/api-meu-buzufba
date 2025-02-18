package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain/repository"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/logger"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/mail"
	"github.com/brnocorreia/api-meu-buzufba/internal/api/shared/rest_err"
	"go.uber.org/zap"

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

	SignOut(c *gin.Context) *rest_err.RestErr

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

	logger.Info("user created in database", zap.String("user_id", user.GetID()))

	html, htmlErr := mail.ParseWelcomeTemplate(mail.WelcomeEmailData{
		Name:         user.GetFirstName(),
		DashboardURL: "https://buzufba.condosnap.com.br",
	})
	if htmlErr != nil {
		logger.Error("error parsing welcome template", htmlErr, zap.String("user_id", user.GetID()))
	}

	mailId, mailErr := mail.Send(mail.EmailParams{
		To:      user.GetEmail(),
		Subject: "Bem-vindo ao Meu Buzufba",
		Html:    html,
	})
	if mailErr != nil {
		logger.Error("error sending welcome email", mailErr, zap.String("user_id", user.GetID()))
		return user, nil
	}

	logger.Info("welcome email sent", zap.String("mail_id", mailId), zap.String("user_id", user.GetID()))

	return user, nil
}

func (as *authService) SignIn(email, password string) (token string, err *rest_err.RestErr) {
	user, err := as.findUserByEmail(email)
	if err != nil {
		return "", err
	}

	if !user.ComparePassword(password) {
		return "", rest_err.NewUnauthorizedRequestError("invalid credentials")
	}

	token, err = as.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Bearer %s", token), nil
}

func (as *authService) SignOut(c *gin.Context) *rest_err.RestErr {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
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
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", rest_err.NewInternalServerError("error trying to generate jwt token")
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
	c.Set("user_domain", userDomain)

	logger.Info(fmt.Sprintf("User authenticated: %#v", userDomain))
}

func RemoveBearerPrefix(token string) string {
	return strings.TrimPrefix(token, "Bearer ")
}

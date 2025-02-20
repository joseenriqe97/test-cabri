package application

import (
	"net/http"
	"time"

	"github.com/ansel1/merry"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joseenriqe97/test-cabri/config"
	"github.com/joseenriqe97/test-cabri/pkg/infrastructure/database"
	"github.com/joseenriqe97/test-cabri/pkg/infrastructure/domain"
	"github.com/joseenriqe97/test-cabri/pkg/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"` // Opcional: mensaje adicional
	Data    interface{} `json:"data,omitempty"`    // Opcional: datos de la respuesta
}

func GenerateResponse(c echo.Context, status int, data interface{}, message string) error {
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	return c.JSON(status, response)
}

type UserRequest struct {
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SignatureOfAuth struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}
type AuthReponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserInterface interface {
	Created(c echo.Context) error
	Authenticate(c echo.Context) error
	GetByID(c echo.Context) error
}
type userApplication struct {
	userDatabase database.UserRepositoryInteface
	servicePLD   service.PldServiceInterface
}

// Validate validates the user request
func (request *UserRequest) Validate() error {
	return validation.ValidateStruct(request,
		validation.Field(&request.Name, validation.Required),
		validation.Field(&request.LastName, validation.Required),
		validation.Field(&request.Email, is.Email),
		validation.Field(&request.Password, validation.Required),
	)
}

// Validate validates the auth request
func (request *AuthRequest) Validate() error {
	return validation.ValidateStruct(request,
		validation.Field(&request.Email, validation.Required),
		validation.Field(&request.Password, validation.Required),
	)
}

func NewUserApplication() UserInterface {
	return &userApplication{
		userDatabase: database.UserDatabase,
		servicePLD:   service.PldService,
	}
}

func (app *userApplication) Created(ctx echo.Context) error {
	var req UserRequest

	if err := ctx.Bind(&req); err != nil {
		return merry.Wrap(err).WithHTTPCode(http.StatusBadRequest)
	}

	if err := req.Validate(); err != nil {
		return merry.Wrap(err).WithHTTPCode(http.StatusBadRequest)
	}

	existUser, err := app.userDatabase.GetByEmail(req.Email)
	if err != nil {
		return merry.Wrap(err)
	}

	if existUser.ID != "" {
		return GenerateResponse(ctx, http.StatusConflict, nil, "Email already exists")
	}

	isBlackList, err := app.servicePLD.CheckPld(service.RequestPld{
		FirstName: req.Name,
		LastName:  req.LastName,
		Email:     req.Email,
	})
	if err != nil {
		return merry.Wrap(err)
	}
	if isBlackList.IsInBlacklist {
		return GenerateResponse(ctx, http.StatusForbidden, nil, "User is blacklisted")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return merry.Wrap(err)
	}

	user := &domain.User{
		Name:     req.Name,
		LastName: req.LastName,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	newUser, err := app.userDatabase.Create(user)
	if err != nil {
		log.WithError(err).Errorf("Error creating user: %s", user.Email)
		return GenerateResponse(ctx, http.StatusInternalServerError, nil, "Error creating user")
	}

	return GenerateResponse(ctx, http.StatusCreated, map[string]interface{}{
		"id":       newUser.ID,
		"name":     newUser.Name,
		"lastName": newUser.LastName,
		"email":    newUser.Email,
	}, "User created successfully")
}

func (app *userApplication) Authenticate(ctx echo.Context) error {
	var request AuthRequest
	if err := ctx.Bind(&request); err != nil {
		return merry.Wrap(err).WithHTTPCode(http.StatusBadRequest)
	}

	if err := request.Validate(); err != nil {
		return merry.Wrap(err).WithHTTPCode(http.StatusBadRequest)
	}

	user, err := app.userDatabase.GetByEmail(request.Email)
	if err != nil {
		return merry.Wrap(err).WithHTTPCode(http.StatusNotFound)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(request.Password),
	); err != nil {
		return GenerateResponse(ctx, http.StatusNotFound, nil, "Incorrect email or password")
	}

	if user.ID == "" {
		return GenerateResponse(ctx, http.StatusNotFound, nil, "User not found")
	}

	token, err := generateToken(user.ID)
	if err != nil {
		return merry.Wrap(err)
	}

	responseToken := map[string]interface{}{
		"token": token,
		"user": User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}
	return GenerateResponse(ctx, http.StatusOK, responseToken, "User authenticated successfully")

}

func (app *userApplication) GetByID(ctx echo.Context) error {
	userId := ctx.Get("id").(string)
	user, err := app.userDatabase.GetByID(userId)

	if err != nil {
		return merry.Wrap(err)
	}

	userResponse := User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	return GenerateResponse(ctx, http.StatusOK, userResponse, "User found successfully")

}

func generateToken(userId string) (string, error) {
	expirationTime := jwt.NewNumericDate(time.Now().Add(72 * time.Hour))

	claims := &SignatureOfAuth{
		ID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetSecretJWT()))
}

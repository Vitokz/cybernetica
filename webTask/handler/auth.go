package handler

import (
	"crypto/sha256"
	"fmt"
	"main/model"
	"main/proto"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h *Handler) Auth(c echo.Context) error {
	ctx := c.Request().Context()            //Создаю context для ограничения времени поиска в бд
	log := logrus.WithFields(logrus.Fields{ //Создаю ивент для обозначения текущих действий
		"event": "registration new user",
	})

	user := new(model.User) //Регистрационная модель юзера

	err := c.Bind(&user) //Присваиваю полученный Json моей модели юзера
	if err != nil {
		log.Printf("Failed processing add user request: %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	user.Password = hash(user.Password) //хэширую пароль

	log.WithFields(logrus.Fields{ //Вывожу в лог инфу о юзере
		"login":    user.Login,
		"name":     user.Name,
		"id":       user.Id,
		"password": user.Password,
	}).Println()

	user, err = h.Db.Auth(ctx, user) //Отправляю запрос в бд для добавления нового юзера
	if err != nil {
		log.WithFields(logrus.Fields{
			"login": user.Login,
			"name":  user.Name,
		}).WithError(err).Error("failed to insert new user")
	}

	log.Info("new user created")
	return c.JSON(http.StatusCreated, user) //Вывожу ответ в виде инфы о юзере без пароля
}

func (h *Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()            //Создаю context для ограничения по времени запроса в бд
	log := logrus.WithFields(logrus.Fields{ //Отправляю в лог инфу о новом ивенте
		"event": "login user",
	})

	login := new(model.Login) //Создаю модель авторизации юзера
	err := c.Bind(&login)     //Присваиваю json своей модели
	if err != nil {
		log.Printf("Failed processing add login request: %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	login.Password = hash(login.Password) //Хэширую пароль

	login, err = h.Db.Login(ctx, login) //Создаю запрос в бд для поиска юзера
	if err != nil {
		log.WithError(err).Error("failed to selecting user")
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	token, err := createJwtToken(*login) //Создаю JWT токен
	if err != nil {
		log.WithError(err).Error("failed to create jwt token")
		return c.String(http.StatusInternalServerError, "err in create a jwt token")
	}

	c.SetCookie(createJwtCookie(token.Token)) //Создаю Куку с этим токеном в значении

	log.Info("user login is competed")
	return c.JSON(http.StatusOK, token) //Возвращаю токен юзеру
}

func hash(value string) string { //Функция хэширования
	hash := sha256.New()
	hash.Write([]byte(value))
	hash.Write([]byte(proto.SALT))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func createJwtToken(user model.Login) (proto.LoginResponse, error) { //Функция создания токена
	claims := model.JwtClaims{ //Поле для значений у токена
		user.Id,
		user.Login,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //Тут идет создание header и payload у токена
	token, err := rawToken.SignedString([]byte(proto.JWTKEY))     //Тут я присваиваю  секретное значение к токену
	if err != nil {
		return proto.LoginResponse{}, err
	}
	return proto.LoginResponse{
		Token: token,
	}, nil //Возвращаю токен
}

func createJwtCookie(token string) *http.Cookie { //Функция создания куки
	cookie := &http.Cookie{}

	cookie.Name = "JWTCookie"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	return cookie
}

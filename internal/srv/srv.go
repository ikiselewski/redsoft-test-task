package srv

import (
	"fmt"
	"math/rand"
	"net/http"
	"redsoft-test-task/internal/database"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/masonkmeyer/agify"
	"github.com/masonkmeyer/genderize"
	"github.com/masonkmeyer/nationalize"
)

type SrvConfig struct {
}

type Dependencies struct {
	Database    database.ReadModel
	Nationalize *nationalize.Client
	Genderize   *genderize.Client
	Agify       *agify.Client
}

type server struct {
	db   database.ReadModel
	r    *gin.Engine
	cfg  SrvConfig
	deps Dependencies
}

func New(cfg *SrvConfig, deps *Dependencies) (ServerInterface, error) {
	r := gin.Default()
	return &server{db: deps.Database, r: r}, nil
}

func (srv *server) ListUsers(c *gin.Context) {
	fmt.Printf("got request %s %s %s", c.Request.URL, c.Request.Method, c.HandlerName())
}

func (srv *server) CreateUser(c *gin.Context) {

	fmt.Printf("got request %s %s %s", c.Request.URL, c.Request.Method, c.HandlerName())

	var body CreateUserJSONRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("can't parse request body due to %s", err.Error()))
		return
	}

	id := uuid.New()

	emails := convertOAPIEmailsToReadmodel(body.Emails, id)

	var wg sync.WaitGroup
	wg.Add(3)

	var (
		age         int
		nationality string
		gender      string
	)

	go func() {
		prediction, _, err := srv.deps.Nationalize.Predict(fmt.Sprintf("%s %s", body.FirstName, body.Surname))
		if err != nil || len(prediction.Country) == 0 {
			nationality = "unknown"
			return
		}
		nationality = prediction.Country[0]

		wg.Done()
	}()

	go func() {
		prediction, _, err := srv.deps.Agify.Predict(fmt.Sprintf("%s %s", body.FirstName, body.Surname))
		if err != nil {
			age = rand.Intn(100)
			return
		}
		age = prediction.Age
		wg.Done()
	}()

	go func() {
		prediction, _, err := srv.deps.Genderize.Predict(fmt.Sprintf("%s %s", body.FirstName, body.Surname))
		if err != nil {
			gender = "unknown"
			return
		}
		gender = prediction.Gender
		wg.Done()
	}()

	wg.Wait()

	srv.db.CreateUser(c.Request.Context(), &database.User{
		ID:         id,
		FirstName:  body.FirstName,
		Surname:    body.Surname,
		Patronymic: body.Patronymic,
		Age:        age,
		Gender:     gender,
		Email:      emails,
	})

}

func (srv *server) UpdatePerson(c *gin.Context, id int64) {

	fmt.Printf("got request %s %s %s", c.Request.URL, c.Request.Method, c.HandlerName())
}

func (srv *server) ListFriends(c *gin.Context, id int64) {

	fmt.Printf("got request %s %s %s", c.Request.URL, c.Request.Method, c.HandlerName())
}

func (srv *server) AddFriend(c *gin.Context, id int64, params AddFriendParams) {

	fmt.Printf("got request %s %s %s", c.Request.URL, c.Request.Method, c.HandlerName())
}

func (srv *server) SearchPersonsBySurname(c *gin.Context, surname string) {

	fmt.Printf("got request %s %s %s", c.Request.URL, c.Request.Method, c.HandlerName())
}

func convertOAPIEmailsToReadmodel(emails *[]string, userID uuid.UUID) []*database.Email {
	if emails == nil {
		return nil
	}
	res := make([]*database.Email, len(*emails))
	for i, email := range *emails {
		res[i] = &database.Email{
			UserID:  userID,
			Address: email,
		}
	}
	return res
}

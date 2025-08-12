package srv

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"redsoft-test-task/internal/database"
	"redsoft-test-task/internal/misc"
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
	r    *gin.Engine
	cfg  SrvConfig
	deps Dependencies
}

func New(cfg *SrvConfig, deps *Dependencies) (ServerInterface, error) {
	r := gin.Default()
	return &server{r: r}, nil
}

func (srv *server) ListUsers(c *gin.Context, params ListUsersParams) {
	if params.Limit < 0 || params.Offset < 0 {
		c.JSON(http.StatusBadRequest, "invalid pagination")
		return
	}

	users, total, err := srv.deps.Database.GetAllUsers(c.Request.Context(), params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("can't parse request body due to %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, ListUsersResponse{Users: convertDBUsersToOAPI(users), Total: total})

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
		nationality = prediction.Country[0].CountryId

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

	srv.deps.Database.CreateUser(c.Request.Context(), &database.User{
		ID:          id,
		FirstName:   body.FirstName,
		Surname:     body.Surname,
		Patronymic:  body.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
		Emails:      misc.StrSlicePtrToStrSlice(body.Emails),
	})

}

func (srv *server) UpdateUser(c *gin.Context, id int64) {
	c.String(http.StatusNotImplemented, "not yet")
}

func (srv *server) ListFriends(c *gin.Context, id int64) {
	c.String(http.StatusNotImplemented, "not yet")
}

func (srv *server) AddFriend(c *gin.Context, id int64, params AddFriendParams) {
	c.String(http.StatusNotImplemented, "not yet")
}

func (srv *server) SearchUsersBySurname(c *gin.Context, surname string) {
	usr, err := srv.deps.Database.SearchUsersBySurname(c.Request.Context(), surname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, "Not found")
			return
		}
		c.JSON(http.StatusInternalServerError, "Something went wrong")
		return
	}

	c.JSON(http.StatusOK, usr)
}

func convertDBUsersToOAPI(users []*database.User) []User {
	res := make([]User, len(users))
	for i, user := range users {

		res[i] = User{
			Age:         user.Age,
			Emails:      &[]string{},
			Gender:      user.Gender,
			Id:          user.ID,
			Name:        user.FirstName,
			Nationality: user.Nationality,
			Patronymic:  misc.StrPtrToStr(user.Patronymic),
			Surname:     user.Surname,
		}
	}
	return res
}

package controller

import (
	"authorization/model"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"authorization/common"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type UserCtrl struct {
	Router        *mux.Router
	DB            *gorm.DB
	JwtMiddleware *jwtmiddleware.JWTMiddleware
}

func (ctrl UserCtrl) InitializeRoutes() {
	ctrl.Router.HandleFunc("/api/login", ctrl.userLogin).Methods("POST")
	ctrl.Router.HandleFunc("/api/signup", ctrl.createUser).Methods("POST")
	ctrl.Router.HandleFunc("/api/profileImg", ctrl.getProfileImg).Methods("GET")
	ctrl.Router.HandleFunc("/api/profileImg", ctrl.uploadProfileImg).Methods("POST")
}

func (ctrl UserCtrl) uploadProfileImg(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("profileImg")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()
	f, err := os.OpenFile("images/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}

func (ctrl UserCtrl) getProfileImg(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("images/download.jpeg")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Conent-type", "image/jpeg")
	w.Write(file)
}

func (ctrl UserCtrl) getUser(user *model.User) (err error) {
	if err := ctrl.DB.Where("email = ?", user.Email).First(user).Error; err != nil {
		return err
	}

	return nil
}

func (ctrl UserCtrl) userLogin(w http.ResponseWriter, r *http.Request) {
	var user model.User
	var userVM = make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := ctrl.getUser(&user); err == nil {
		userVM["email"] = user.Email
		userVM["token"] = getToken()
		common.RespondWithJSON(w, http.StatusOK, userVM)
	} else {
		common.RespondWithError(w, http.StatusInternalServerError, "Can't find user, Please sign up")
	}
}

func getToken() interface{} {
	var rawResponse map[string]interface{}

	url := "https://brian-880120.auth0.com/oauth/token"
	payload := strings.NewReader("{\"client_id\":\"Bxfw9p2eZx1QV9kFDFeze3xqXRyQ1EqV\",\"client_secret\":\"b8nGJhKYrIeB9B3H6IMTDSfTmw_JqbtNd-_uauvm1zU30mOTK16M8a3Vb554Q3nq\",\"audience\":\"academy-auth\",\"grant_type\":\"client_credentials\"}")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	json.Unmarshal(body, &rawResponse)
	return rawResponse
}

func (ctrl UserCtrl) createUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := ctrl.DB.Create(&user).Error; err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	common.RespondWithJSON(w, http.StatusOK, user)
}

package controller

import (
	"authorization/model"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
	ctrl.Router.HandleFunc("/api/profileImg/{id:[0-9]+}", ctrl.getProfileImg).Methods("GET")
	ctrl.Router.HandleFunc("/api/profileImg/{id:[0-9]+}", ctrl.uploadProfileImg).Methods("POST")
	ctrl.Router.HandleFunc("/api/users/{id:[0-9]+}", ctrl.updateUser).Methods("PUT")
	ctrl.Router.HandleFunc("/api/users/{id:[0-9]+}", ctrl.getUserById).Methods("GET")
}

func (ctrl UserCtrl) getUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 8)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	user := model.User{}
	user.ID = uint(id)
	if err := ctrl.DB.First(&user).Error; err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, "User doesn't exist")
	} else {
		common.RespondWithJSON(w, http.StatusOK, user)
	}
}

func (ctrl UserCtrl) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 8)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user := model.User{}
	user.ID = uint(id)

	var updateUser model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateUser); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	ctrl.DB.Model(&user).Updates(updateUser)
	common.RespondWithJSON(w, http.StatusOK, user)
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

	// save image path
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 8)
	user := model.User{}
	user.ID = uint(id)
	ctrl.DB.First(&user)
	user.PhotoPath = "images/" + handler.Filename
	ctrl.DB.Save(&user)
}

func (ctrl UserCtrl) getProfileImg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 8)
	user := model.User{}
	user.ID = uint(id)
	ctrl.DB.First(&user)

	file, err := os.Open(user.PhotoPath)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	defer file.Close()

	fInfo, _ := file.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	fReader := bufio.NewReader(file)
	fReader.Read(buf)

	imgBase64Str := base64.StdEncoding.EncodeToString(buf)

	w.Write([]byte(imgBase64Str))
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
		userVM["user"] = user
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

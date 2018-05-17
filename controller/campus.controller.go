package controller

import (
	"authorization/common"
	"authorization/model"
	"encoding/json"
	"net/http"

	"regexp"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type CampCtrl struct {
	Router        *mux.Router
	DB            *gorm.DB
	JwtMiddleware *jwtmiddleware.JWTMiddleware
}

func (ctrl CampCtrl) InitializeRoutes() {
	ctrl.Router.HandleFunc("/api/camp/registerStudent", ctrl.registerStudent).Methods("POST")
}

func (ctrl CampCtrl) registerStudent(w http.ResponseWriter, r *http.Request) {
	var student model.Student
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&student); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if student.Name == "" || student.Email == "" || student.Phone == "" {
		common.RespondWithError(w, http.StatusInternalServerError, "Required fields cannot be empty")
		return
	}

	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !re.MatchString(student.Email) {
		common.RespondWithError(w, http.StatusInternalServerError, "invalid email format")
		return
	}

	defer r.Body.Close()

	if err := ctrl.DB.Create(&student).Error; err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	common.RespondWithJSON(w, http.StatusOK, student)
}

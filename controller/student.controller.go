package controller

import (
	"authorization/model"
	"net/http"

	"authorization/common"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type StudentCtrl struct {
	Router        *mux.Router
	DB            *gorm.DB
	JwtMiddleware *jwtmiddleware.JWTMiddleware
}

func (ctrl StudentCtrl) InitializeRoutes() {
	ctrl.Router.HandleFunc("/api/students", ctrl.getStuentsByCamp).Queries("campId", "{camp}").Methods("GET")
	ctrl.Router.HandleFunc("/api/testLocale", ctrl.testLocale).Methods("GET")
}

func (ctrl StudentCtrl) testLocale(w http.ResponseWriter, r *http.Request) {
	locationString := r.Header.Get("Accept-Language")
	common.RespondWithJSON(w, http.StatusOK, locationString)
}

func (ctrl StudentCtrl) getStuentsByCamp(w http.ResponseWriter, r *http.Request) {
	campId := r.URL.Query().Get("campId")
	userVMs := []model.StudentVM{}
	err := ctrl.DB.Model(&model.Student{}).Joins("join campus on students.campus_id = campus.id").
		Where("students.campus_id = ?", campId).
		Select("students.name as name, students.email as email, students.phone as phone, students.company as company, campus.name as campusname").
		Scan(&userVMs).Error

	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	common.RespondWithJSON(w, http.StatusOK, userVMs)
}

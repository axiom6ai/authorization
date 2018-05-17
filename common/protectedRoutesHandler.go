package common

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
)

type routesHandler func(w http.ResponseWriter, r *http.Request)

func HandlePrivateRoutes(jwtMiddleware *jwtmiddleware.JWTMiddleware, fn routesHandler) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(fn)))
}

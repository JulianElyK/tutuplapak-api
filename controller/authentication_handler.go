package controller

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("jsrs6d5rRpHsDmIp")
var tokenName = "token"

type Claims struct {
	ID       string `json:id`
	Name     string `json:name`
	UserType string `json:user_type`
	jwt.StandardClaims
}

func generateToken(w http.ResponseWriter, email string, name string, userType string) {
	expiry := 24 * time.Hour
	if userType == "A" {
		expiry = 30 * time.Minute
	}
	tokenExpiryTime := time.Now().Add(expiry)

	claims := &Claims{
		ID:       email,
		Name:     name,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    signedToken,
		Expires:  tokenExpiryTime,
		Secure:   false,
		HttpOnly: true,
	})
}

func resetUsersToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    "",
		Expires:  time.Now(),
		Secure:   false,
		HttpOnly: true,
	})
}

func Authenticate(next http.HandlerFunc, accessType string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isValidToken := validateUserToken(w, r, accessType)
		if !isValidToken {
			sendUnAuthorizedResponse(w)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func validateUserToken(w http.ResponseWriter, r *http.Request, accessType string) bool {
	isAccessTokenValid, userType := validateTokenFromCookies(r)
	if isAccessTokenValid {
		return accessType == userType
	}
	return false
}

func validateTokenFromCookies(r *http.Request) (bool, string) {
	if coockie, err := r.Cookie(tokenName); err == nil {
		accessToken := coockie.Value
		accessClaims := &Claims{}
		parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			return true, accessClaims.UserType
		}
	}
	return false, ""
}

func getEmailType(r *http.Request) (string, string) {
	if coockie, err := r.Cookie(tokenName); err == nil {
		accessToken := coockie.Value
		accessClaims := &Claims{}
		parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			return accessClaims.ID, accessClaims.UserType
		}
	}
	return "", ""
}

package jwtmidware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func Login(w http.ResponseWriter, req *http.Request) {
	log.Println("Login")

	//Creating Access Token
	//os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = "rich"
	atClaims["exp"] = time.Now().Add(time.Minute * 2).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, _ := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	mapD := map[string]string{"apple": token}
	jsonresp, _ := json.Marshal(mapD)
	w.Write([]byte(jsonresp))
}

func ServeProtected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the Authheader
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 { //Authheader is bad
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		}

		//check the token
		jwtToken := authHeader[1]
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			value := os.Getenv("ACCESS_SECRET")
			return []byte(value), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "props", claims)
			// Access context values in handlers like this
			// props, _ := r.Context().Value("props").(jwt.MapClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		}
	})
}

func Logout(w http.ResponseWriter, req *http.Request) {
	log.Println("Logout")
	io.WriteString(w, "Hello , just logged out\n")
}

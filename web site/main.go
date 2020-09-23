package main

import (
	"./database"
	"context"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	//db []*database.User
	tokenstring = []byte("secret")
	db *sql.DB
)

func indexHandler(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil{
		fmt.Println("template pars error: ", err)
	}
	tmpl.Execute(w, nil)
}

func saveloginHandler(w http.ResponseWriter, r *http.Request){
	db := database.GetDB()
	defer db.Close()
	login := r.FormValue("login")
	pass := r.FormValue("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil{
		fmt.Println("hash error: ", err)
	}
	result, err := db.Exec("INSERT INTO users(login, password) VALUES($1,$2) ", r.FormValue("login"), string(hashedPassword))
	if err != nil{
		fmt.Println("insert db error: ", err)
	}
	fmt.Println(result.RowsAffected())
	token, _ := GenerateJWT(login)
	cookie := http.Cookie{
		Name: "login",
		Value: token,
	}
	http.SetCookie(w, &cookie)
	ctx := context.WithValue(context.Background(), "Name", login)
	//r.Header.Set("Location", "http://localhost:8000/profile")
	http.Redirect(w,r.WithContext(ctx),"/profile", 302)
}

func profileHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println(strings.TrimLeft(r.Header["Cookie"][0], "login="))
	ctx := r.Context()
	login := ctx.Value("Name")
	fmt.Fprintf(w, "%s", login)
}

func GenerateJWT(login string)(string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"auth": true,
		"user": login,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	})

	tokenString, err := token.SignedString(tokenstring)
	if err != nil {
		fmt.Errorf("Something went wrong :%s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func isAuthorized(endpoint func (w http.ResponseWriter, r *http.Request))http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		if r.Header["Cookie"] != nil{
			token, err := jwt.Parse(strings.TrimLeft(r.Header["Cookie"][0], "login="), func(token *jwt.Token) (interface{},error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("there was an error")
				}
				return tokenstring, nil
			})
			if err != nil{
				fmt.Fprintf(w, err.Error())
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
				ctx := context.WithValue(context.Background(), "Name", claims["user"])
				endpoint(w, r.WithContext(ctx))
			}
		}else{
			//fmt.Fprintf(w, "Not authorized")
			http.Redirect(w, r, "/login", 302)
		}
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request){
	if r.Header["Cookie"] != nil{
		http.SetCookie(w, &http.Cookie{
			Name: "login",
			Value: "",
			MaxAge: -1,
			Expires: time.Now().Add(-100 * time.Hour),// set expires for older versions
		})
		http.Redirect(w,r,"/login", 302)
	}else{
		fmt.Fprintf(w, "You are not registred")
	}
}


func main(){
	router := chi.NewRouter()
	router.Get("/login", indexHandler)
	router.Post("/Save", saveloginHandler)
	router.Handle("/profile", isAuthorized(profileHandler))
	router.Get("/logout", logoutHandler)

	log.Fatal(http.ListenAndServe(":8000", router))
}
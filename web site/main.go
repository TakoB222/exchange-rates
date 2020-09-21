package main

import (
	"./database"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	db []*database.User
	tokenstring = []byte("secret")
)

func indexHandler(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil{
		fmt.Println("template pars error: ", err)
	}
	tmpl.Execute(w, nil)
}

func saveloginHandler(w http.ResponseWriter, r *http.Request){
	/*db := database.GetDB()
	defer db.Close()
	err := db.QueryRow("INSERT INTO users(login, pass) VALUES($1,$2) ", r.FormValue("login"), r.FormValue("password"))
	if err != nil{
		fmt.Println("insert db error: ", err)
	}*/
	login := r.FormValue("login")
	pass := r.FormValue("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil{
		fmt.Println("hash error: ", err)
	}
	token, _ := GenerateJWT(login)
	db = append(db, &database.User{Login: login, HashPassword: hashedPassword, JWT: token})
	fmt.Printf("%+v", db)
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
	ctx := r.Context()
	fmt.Fprintf(w, ctx.Value("Name").(string))
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
			if cookie, err := r.Cookie("login"); err != nil{
				fmt.Println("cookie error: ", err)
			}
			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{},error) {
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
			fmt.Fprintf(w, "Not authorized")
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
	//router.Get("/profile", profileHandler)
	router.Handle("/profile", isAuthorized(profileHandler))

	log.Fatal(http.ListenAndServe(":8000", router))
}
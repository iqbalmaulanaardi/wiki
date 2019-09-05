package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
	"wiki/driver"
	"wiki/helper"
	"wiki/models"
	"wiki/repository"
	"wiki/repository/user"
	jwt "github.com/dgrijalva/jwt-go"
)

func NewUserHandler(db *driver.DB) *User{
	return &User{
		repo: user.NewSQLUserRepo(db.SQL),
	}
}
type User struct {
	repo repsitory.UserRepo
}
type JwtObject struct {
	ID int64 `json:"id"`
	Email string `json:"email"`
	ExpiredDate int64 `json:"expiredDate"`
}
func (u *User) Login(w http.ResponseWriter, r *http.Request){
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user) //return struct not byte
	resUser,errGetByEmail := u.repo.GetByEmail(r.Context(), user.Email)// return struct
	if errGetByEmail != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Account not found")
	}
	resUserPwdByte := []byte(resUser.Password)
	userPwdByte := []byte(user.Password)
	errCompare := bcrypt.CompareHashAndPassword(resUserPwdByte,userPwdByte)
	if errCompare != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Incorrect Password")
	}else{
		//generate token
		timein := time.Now().Local().Add(time.Minute * time.Duration(5)).Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"ID": resUser.ID,
			"Email": resUser.Email,
			"ExpiredDate": timein,
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, errSign := token.SignedString([]byte("trygloenv"))
		if errSign != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to sign jwt")
		}
		//try parsing
		parsedToken, _ := jwt.Parse(tokenString, func(tok *jwt.Token)(interface{}, error) {
			if jwt.GetSigningMethod("HS256") != tok.Method {
				return nil, models.ErrNotFound
			} //return jwt methods
			return []byte("trygloenv"), nil
		})
		claims, ok := token.Claims.(jwt.MapClaims); //return map
		if ok && parsedToken.Valid {
			fmt.Println(claims["Email"])
		} else {
			fmt.Println("ok")
		}
		//end try parsing
		helper.RespondwithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
	}
}
func (u *User) Fetch(w http.ResponseWriter, r *http.Request){
	payload, _ := u.repo.Fetch(r.Context(),5)
	helper.RespondwithJSON(w, http.StatusOK, payload)
}
func(u *User) GetById(w http.ResponseWriter, r *http.Request){
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	payload, err := u.repo.GetByID(r.Context(), int64(id))
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Content Not Found")
	}
	helper.RespondwithJSON(w, http.StatusOK, payload)

}
func(u *User) Signup(w http.ResponseWriter, r *http.Request){
	now := time.Now().Unix()
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)
	user.CreatedAt = now
	hashedPassword, errHashedPassword := bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.MinCost)
	if errHashedPassword != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "error generate from password")
	}
	user.Password = string(hashedPassword)
	_, errGetByEmail := u.repo.GetByEmail(r.Context(), user.Email)
	if errGetByEmail != nil {
		_, err := u.repo.Create(r.Context(), &user)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Server Error")
		}
		helper.RespondwithJSON(w, http.StatusCreated, user)
	}else{
		helper.RespondWithError(w, http.StatusInternalServerError, "Email Already Taken")
	}

}
func(u *User) Update(w http.ResponseWriter, r *http.Request){
	fmt.Println("try update user")
}
func(u *User) Delete(w http.ResponseWriter, r *http.Request){
	fmt.Println("try delete user")
}
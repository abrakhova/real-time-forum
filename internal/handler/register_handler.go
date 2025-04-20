package handler

import (
	"forum/internal/model"
	"forum/internal/user"
	"forum/internal/util"
	"net/http"
	"strconv"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		util.ExecuteJSON(w, model.MsgData{"Invalid request method"}, http.StatusMethodNotAllowed)
		return
	}

	// Collect form data
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	ageStr := r.FormValue("age")
	gender := r.FormValue("gender")

	// Validate required fields
	if username == "" || email == "" || password == "" || 
	   firstName == "" || lastName == "" || ageStr == "" || gender == "" {
		util.ExecuteJSON(w, model.MsgData{"All fields are required"}, http.StatusBadRequest)
		return
	}

	// Parse and validate age
	age, err := strconv.Atoi(ageStr)
	if err != nil || age <= 0 {
		util.ExecuteJSON(w, model.MsgData{"Invalid age"}, http.StatusBadRequest)
		return
	}

	// Validate email format
	if !util.IsValidEmail(email) {
		util.ExecuteJSON(w, model.MsgData{"Invalid email format"}, http.StatusBadRequest)
		return
	}

	// Check if username exists
	exists, err := user.UsernameExists(username)
	if err != nil || exists {
		util.ExecuteJSON(w, model.MsgData{"Username already taken"}, http.StatusConflict)
		return
	}

	// Check if email exists
	exists, err = user.EmailExists(email)
	if err != nil || exists {
		util.ExecuteJSON(w, model.MsgData{"Email already taken"}, http.StatusConflict)
		return
	}

	// Hash the password
	hashedPassword, err := user.HashPassword(password)
	if err != nil {
		util.ExecuteJSON(w, model.MsgData{"Password hashing failed"}, http.StatusInternalServerError)
		return
	}

	// Save user to the database
	if err := user.SaveUser(username, email, hashedPassword, firstName, lastName, gender, age); err != nil {
		util.ExecuteJSON(w, model.MsgData{"User registration failed"}, http.StatusInternalServerError)
		return
	}

	// Return successful registration response
	util.ExecuteJSON(w, model.MsgData{"Registration successful!"}, http.StatusOK)
}
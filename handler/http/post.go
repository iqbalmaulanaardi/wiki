package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"wiki/driver"
	"wiki/helper"
	"wiki/models"
	repository "wiki/repository"
	"wiki/repository/post"
)

// NewPostHandler ...
func NewPostHandler(db *driver.DB) *Post {
	return &Post{
		repo: post.NewSQLPostRepo(db.SQL),
	}
}
// Post ...
type Post struct {
	repo repository.PostRepo
}

// Fetch all post data
func (p *Post) Fetch(w http.ResponseWriter, r *http.Request) {
	payload, _ := p.repo.Fetch(r.Context(), 5)

	helper.RespondwithJSON(w, http.StatusOK, payload)
}

// Create a new post
func (p *Post) Create(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}
	json.NewDecoder(r.Body).Decode(&post)

	newID, err := p.repo.Create(r.Context(), &post)
	fmt.Println(newID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Server Error")
	}

	helper.RespondwithJSON(w, http.StatusCreated, post)
}

// Update a post by id
func (p *Post) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	data := models.Post{ID: int64(id)}
	json.NewDecoder(r.Body).Decode(&data)
	payload, err := p.repo.Update(r.Context(), &data)

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Server Error")
	}

	helper.RespondwithJSON(w, http.StatusOK, payload)
}

// GetByID returns a post details
func (p *Post) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	payload, err := p.repo.GetByID(r.Context(), int64(id))
	if err != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Content not found")
	}

	helper.RespondwithJSON(w, http.StatusOK, payload)
}

// Delete a post
func (p *Post) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	//find first
	_, errGet := p.repo.GetByID(r.Context(), int64(id))
	if errGet != nil {
		helper.RespondWithError(w, http.StatusNotFound, "Content not found")
	}else{
		_, err := p.repo.Delete(r.Context(), int64(id))
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Server Error")
		}

		helper.RespondwithJSON(w, http.StatusMovedPermanently, map[string]string{"message": "Delete Successfully"})
	}
	//find end
}



package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/caturandi-labs/go-social/internal/store"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required"`
	Content string   `json:"content" validate:"required"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var post CreatePostPayload
	if err := readJSON(w, r, &post); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(post); err != nil {
		validationErr := formatValidationErrors(err)
		app.unprocessableEntityResponse(w, r, validationErr)
		return
	}

	newPost := &store.Post{
		Title:   post.Title,
		Content: post.Content,
		UserID:  1,
		Tags:    post.Tags,
	}
	ctx := r.Context()
	if err := app.store.Posts.Create(ctx, newPost); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, newPost); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	paramId := chi.URLParam(r, "id")
	fmt.Println(paramId)
	id, err := strconv.ParseInt(paramId, 10, 64)
	ctx := r.Context()
	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	comments, err := app.store.Comments.GetByPostID(ctx, id)
	if err != nil {
		app.internalServerError(w, r, err)
	}
	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"required,min=3"`
	Content *string `json:"content" validate:"required,min=3"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		validationErr := formatValidationErrors(err)
		app.unprocessableEntityResponse(w, r, validationErr)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
	}
	
	ctx := r.Context()

	if err := app.store.Posts.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	_ = app.jsonResponse(w, http.StatusNoContent, nil)
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paramId := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(paramId, 10, 64)
		ctx := r.Context()
		post, err := app.store.Posts.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
		valContext := context.WithValue(ctx, "post", post)
		next.ServeHTTP(w, r.WithContext(valContext))
	})
}

func getPostFromContext(r *http.Request) *store.Post {
	post, _ := r.Context().Value("post").(*store.Post)
	return post
}

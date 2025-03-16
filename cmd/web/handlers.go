package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"toramanomer/snippetbox/internal/models"
	"unicode/utf8"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData()
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData()
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Form = snippetCreateForm{Expires: 365}
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var (
		title        = strings.TrimSpace(r.Form.Get("title"))
		content      = strings.TrimSpace(r.Form.Get("content"))
		expires, err = strconv.Atoi(r.Form.Get("expires"))
	)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:       title,
		Content:     content,
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	if title == "" {
		form.FieldErrors["title"] = "Title cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		form.FieldErrors["title"] = "Title cannot be more than 100 characters long"
	}

	if content == "" {
		form.FieldErrors["content"] = "Content cannot be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		form.FieldErrors["expires"] = "Expiration must be 1, 7 or 365"
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

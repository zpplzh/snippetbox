package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zappel/snippetbox/pkg/forms"
	"github.com/zappel/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, req *http.Request) {

	/*if req.URL.Path != "/" {
		app.notFound(w)
		return
	}*/ //dihapus karena pat library jadi ga usah manual check

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, req, "home.page.tmpl", &templateData{
		Snippets: s,
	})

	/*
		data := &templateData{Snippets: s}

		files := []string{
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/footer.partial.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {*/
	/*app.errorLog.Println(err.Error())
	http.Error(w, "Internal Server Error", 500)*/
	/*	app.serverError(w, err)
			return
		}

		err = ts.Execute(w, data)
		if err != nil {*/
	/*app.errorLog.Println(err.Error())
	http.Error(w, "Internal Server Error", 500)*/
	/*	app.serverError(w, err)
		}*/

}

func (app *application) showSnippet(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		//http.NotFound(w, req)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, req, "show.page.tmpl", &templateData{
		Snippet: s,
	})

	/*data := &templateData{Snippet: s}

	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}*/

	// Pass in the templateData struct when executing the template.
	/*err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}*/
}

func (app *application) createSnippetForm(w http.ResponseWriter, req *http.Request) {
	app.render(w, req, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})

}

func (app *application) createSnippet(w http.ResponseWriter, req *http.Request) {

	/*if req.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) // ini buat masukin ke header kalo yang allow itu POST
		// .Set untuk nambahin tapi kalau ada yang sama .Set akan replace yang sudah ada
		// .Add untuk nambahin jadi bisa banyak ga akan ngereplace yang sudah ada
		//w.WriteHeader(405) //ganti pake http.Error
		//w.Write([]byte("Method Not Allowed")) //ganti pake http.Error lebih singkat
		app.clientError(w, http.StatusMethodNotAllowed)
		//http.Error(w, "{}", 405)
		return
	}*/

	err := req.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	/*title := req.PostForm.Get("title")
	content := req.PostForm.Get("content")
	expires := req.PostForm.Get("expires")*/

	form := forms.New(req.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, req, "create.page.tmpl", &templateData{Form: form})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// bikin penampungan error validasi
	/*errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "this field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "this field is too long"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "this field cannot be blank"
	}

	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "this field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "this field not valid"
	}

	if len(errors) > 0 {
		app.render(w, req, "create.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   req.PostForm,
		})

		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}*/

	app.session.Put(req, "flash", "Snippet successfully added")

	http.Redirect(w, req, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Display the user signup form...")
}

func (app *application) signupUser(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}

func (app *application) loginUserForm(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Display the user login form...")
}

func (app *application) loginUser(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) logoutUser(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}

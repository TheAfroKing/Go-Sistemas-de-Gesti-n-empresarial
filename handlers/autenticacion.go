package handlers

import (
	"Go-Sistemas-de-Gestion-empresarial/models"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		cliente, err := models.Login(email, password)
		if err != nil {
			log.Println("Error de login:", err)
			http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
			return
		}

		expiration := time.Now().Add(24 * time.Hour)
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    fmt.Sprintf("%d", cliente.ID),
			Expires:  expiration,
			Path:     "/",
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "user_perfil",
			Value:    cliente.Perfil,
			Expires:  expiration,
			Path:     "/",
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "true",
			Expires:  expiration,
			Path:     "/",
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/login.html")
	if err != nil {
		log.Println("Error al cargar el template de login", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	data := struct {
		Error      bool
		Registered bool
		LoginToken bool
		Perfil     string
	}{
		Error:      r.URL.Query().Get("error") != "",
		Registered: r.URL.Query().Get("registered") == "true",
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println("Error al ejecutar el template", err)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := models.CreateCliente(r.FormValue("nombre"), r.FormValue("email"), r.FormValue("password"), r.FormValue("direccion"), r.FormValue("telefono"))
		if err != nil {
			log.Println("Error al registrar:", err)
			http.Redirect(w, r, "/register?error=register_failed", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/login?registered=true", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/register.html")
	if err != nil {
		log.Println("Error al cargar el template de registro", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	data := struct {
		Error      bool
		LoginToken bool
		Perfil     string
	}{
		Error: r.URL.Query().Get("error") != "",
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println("Error al ejecutar el template", err)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_perfil",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetSessionData(r *http.Request) (bool, string, string) {
	tokenCookie, err := r.Cookie("token")
	if err != nil || tokenCookie.Value != "true" {
		return false, "", ""
	}
	perfilCookie, err := r.Cookie("user_perfil")
	perfil := ""
	if err == nil {
		perfil = perfilCookie.Value
	}
	idCookie, err := r.Cookie("user_id")
	id := ""
	if err == nil {
		id = idCookie.Value
	}
	return true, perfil, id
}

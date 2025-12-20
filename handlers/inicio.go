package handlers

import (
	"Go-Sistemas-de-Gestion-empresarial/models"
	"html/template"
	"log"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// HomeHandler muestra la página principal con los productos activos disponibles.
	// Filtra productos por `Activo` y `Stock > 0`, carga templates y renderiza la vista.
	productos, err := models.GetAllProductos()
	if err != nil {
		log.Println("Error al obtener los productos", err)
		return
	}

	loggedIn, perfil, _ := GetSessionData(r)

	var activeProductos []models.Producto
	for _, p := range productos {
		if p.Activo && p.Stock > 0 {
			activeProductos = append(activeProductos, p)
		}
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/cliente/productos.html")
	if err != nil {
		log.Println("Error al cargar el template de home", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	data := struct {
		Productos  []models.Producto
		LoginToken bool
		Perfil     string
		ItemAdded  bool
	}{
		Productos:  activeProductos,
		LoginToken: loggedIn,
		Perfil:     perfil,
		ItemAdded:  r.URL.Query().Get("added") == "true",
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println("Error al ejecutar el template", err)
	}
}

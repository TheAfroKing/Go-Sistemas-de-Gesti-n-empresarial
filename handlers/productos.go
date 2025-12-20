package handlers

import (
	"Go-Sistemas-de-Gestion-empresarial/models"
	"log"
	"net/http"
	"strconv"
)

func AgregarItemCarrito(w http.ResponseWriter, r *http.Request) {
	// AgregarItemCarrito procesa la solicitud POST para añadir un producto
	// al carrito del cliente autenticado. Crea un carrito si no existe.

	if r.Method == "POST" {
		loggedIn, _, userIDStr := GetSessionData(r)
		if !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		id, _ := strconv.Atoi(userIDStr)
		carrito, _ := models.GetCarritoByClienteID(id)
		log.Println("Carrito ID encontrado:", carrito.ID)

		if carrito.ID == 0 {
			log.Println("Carrito no encontrado, creando nuevo carrito para cliente:", id)
			errCarrito := models.CreateCarrito(id)
			if errCarrito != nil {
				log.Println("Error al crear el carrito:", errCarrito)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			carrito, _ = models.GetCarritoByClienteID(id)
		}

		rawID := r.FormValue("id_producto")
		idProducto, _ := strconv.Atoi(rawID)
		cantidad, _ := strconv.Atoi(r.FormValue("cantidad"))

		log.Println("DEBUG: Raw id_producto form value:", rawID)
		log.Println("Agregando producto:", idProducto, "Cantidad:", cantidad, "a Carrito:", carrito.ID)

		err := models.AgregarItemCarrito(carrito.ID, idProducto, cantidad)
		if err != nil {
			log.Println("Error al registrar item en carrito:", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/?added=true", http.StatusSeeOther)
		return
	}

}

package handlers

import (
	"Go-Sistemas-de-Gestion-empresarial/models"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ClientProductDetail(w http.ResponseWriter, r *http.Request) {
	loggedIn, perfil, _ := GetSessionData(r)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	producto, err := models.GetProductoByID(id)
	if err != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/client/product_detail.html")
	if err != nil {
		log.Println("Error cargando template client product detail:", err)
		http.Error(w, "Error cargando templates", http.StatusInternalServerError)
		return
	}

	data := struct {
		Producto   models.Producto
		LoginToken bool
		Perfil     string
	}{
		Producto:   producto,
		LoginToken: loggedIn,
		Perfil:     perfil,
	}

	tmpl.ExecuteTemplate(w, "base", data)
}

func ClientCart(w http.ResponseWriter, r *http.Request) {
	loggedIn, perfil, userIDStr := GetSessionData(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(userIDStr)

	models.CreateCarrito(userID)

	carrito, err := models.GetCarritoByClienteID(userID)
	if err != nil {
		log.Println("Error obteniendo carrito:", err)
	}

	items, err := models.GetItemsByCarritoID(carrito.ID)
	if err != nil {
		log.Println("Error obteniendo items del carrito:", err)
	}

	type CartItemDetail struct {
		models.ItemCarrito
		Producto models.Producto
		Subtotal float64
	}

	var cartDetails []CartItemDetail
	var totalCart float64

	for _, item := range items {
		prod, _ := models.GetProductoByID(item.IDProducto)
		subtotal := float64(item.Cantidad) * prod.Precio
		cartDetails = append(cartDetails, CartItemDetail{
			ItemCarrito: item,
			Producto:    prod,
			Subtotal:    subtotal,
		})
		totalCart += subtotal
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/client/cart.html")
	if err != nil {
		log.Println("Error cargando template client cart:", err)
		http.Error(w, "Error cargando templates", http.StatusInternalServerError)
		return
	}

	data := struct {
		CartItems  []CartItemDetail
		Total      float64
		LoginToken bool
		Perfil     string
	}{
		CartItems:  cartDetails,
		Total:      totalCart,
		LoginToken: loggedIn,
		Perfil:     perfil,
	}

	tmpl.ExecuteTemplate(w, "base", data)
}

func RemoveItemFromCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, _ = strconv.Atoi(vars["id"])

	http.Redirect(w, r, "/carrito", http.StatusSeeOther)
}

func ClientCheckout(w http.ResponseWriter, r *http.Request) {
	loggedIn, perfil, userIDStr := GetSessionData(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(userIDStr)
	carrito, _ := models.GetCarritoByClienteID(userID)
	items, _ := models.GetItemsByCarritoID(carrito.ID)

	var totalCart float64
	for _, item := range items {
		prod, _ := models.GetProductoByID(item.IDProducto)
		totalCart += float64(item.Cantidad) * prod.Precio
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/client/checkout.html")
	if err != nil {
		log.Println("Error cargando template client checkout:", err)
		http.Error(w, "Error cargando templates", http.StatusInternalServerError)
		return
	}

	data := struct {
		Total      float64
		LoginToken bool
		Perfil     string
	}{
		Total:      totalCart,
		LoginToken: loggedIn,
		Perfil:     perfil,
	}

	tmpl.ExecuteTemplate(w, "base", data)
}

func ProcessCheckout(w http.ResponseWriter, r *http.Request) {
	loggedIn, _, userIDStr := GetSessionData(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		userID, _ := strconv.Atoi(userIDStr)
		metodoPago := r.FormValue("metodo_pago") // tarjeta, transferencia, etc

		carrito, _ := models.GetCarritoByClienteID(userID)
		items, _ := models.GetItemsByCarritoID(carrito.ID)

		var totalCart float64
		for _, item := range items {
			prod, _ := models.GetProductoByID(item.IDProducto)
			totalCart += float64(item.Cantidad) * prod.Precio
		}

		transaccionID := "imulado_123" // Simulado
		pedidoID, err := models.CreatePedido(userID, totalCart, metodoPago, transaccionID)
		if err != nil {
			log.Println("Error creando pedido:", err)
			http.Error(w, "Error procesando pedido", http.StatusInternalServerError)
			return
		}

		for _, item := range items {
			prod, _ := models.GetProductoByID(item.IDProducto)
			models.CreateDetallePedido(pedidoID, item.IDProducto, item.Cantidad, prod.Precio)

			newStock := prod.Stock - item.Cantidad
			models.UpdateProducto(prod.ID, prod.Nombre, prod.Descripcion, prod.Precio, newStock, prod.SKU, prod.Activo)
		}

		err = models.EmptyCarrito(carrito.ID)
		if err != nil {
			log.Println("Error vaciando carrito:", err)
		}

		http.Redirect(w, r, "/perfil?order_success=true", http.StatusSeeOther)
	}
}

func ClientProfile(w http.ResponseWriter, r *http.Request) {
	loggedIn, perfil, userIDStr := GetSessionData(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(userIDStr)
	cliente, _ := models.GetClienteByID(userID)

	myPedidos, err := models.GetPedidosByClienteID(userID)
	if err != nil {
		log.Println("Error obteniendo pedidos del cliente:", err)
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/client/profile.html")
	if err != nil {
		log.Println("Error cargando template client profile:", err)
		http.Error(w, "Error cargando templates", http.StatusInternalServerError)
		return
	}

	data := struct {
		Cliente    models.Cliente
		Pedidos    []models.Pedido
		LoginToken bool
		Perfil     string
		Success    bool
	}{
		Cliente:    cliente,
		Pedidos:    myPedidos,
		LoginToken: loggedIn,
		Perfil:     perfil,
		Success:    r.URL.Query().Get("order_success") == "true",
	}

	tmpl.ExecuteTemplate(w, "base", data)
}

func ClientProfileEdit(w http.ResponseWriter, r *http.Request) {
	loggedIn, perfil, userIDStr := GetSessionData(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userID, _ := strconv.Atoi(userIDStr)

	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		email := r.FormValue("email")
		telefono := r.FormValue("telefono")
		direccion := r.FormValue("direccion")

		err := models.UpdateCliente(userID, nombre, email, direccion, telefono)
		if err != nil {
			log.Println("Error actualizando perfil:", err)
			http.Error(w, "Error actualizando perfil", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/perfil", http.StatusSeeOther)
		return
	}

	cliente, _ := models.GetClienteByID(userID)

	tmpl, err := template.ParseFiles("templates/base.html", "templates/client/profile_edit.html")
	if err != nil {
		log.Println("Error cargando template client profile edit:", err)
		http.Error(w, "Error cargando templates", http.StatusInternalServerError)
		return
	}

	data := struct {
		Cliente    models.Cliente
		LoginToken bool
		Perfil     string
	}{
		Cliente:    cliente,
		LoginToken: loggedIn,
		Perfil:     perfil,
	}

	tmpl.ExecuteTemplate(w, "base", data)
}

func ClientOrderDetail(w http.ResponseWriter, r *http.Request) {
	loggedIn, perfil, userIDStr := GetSessionData(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userID, _ := strconv.Atoi(userIDStr)

	vars := mux.Vars(r)
	orderID, _ := strconv.Atoi(vars["id"])

	pedido, err := models.GetPedidoByID(orderID)
	if err != nil {
		http.Error(w, "Pedido no encontrado", http.StatusNotFound)
		return
	}

	if pedido.IDCliente != userID {
		http.Error(w, "No autorizado", http.StatusForbidden)
		return
	}

	detalles, _ := models.GetDetallesByPedidoID(orderID)

	tmpl, err := template.ParseFiles("templates/base.html", "templates/client/order_detail.html")
	if err != nil {
		log.Println("Error cargando template client order detail:", err)
		http.Error(w, "Error cargando templates", http.StatusInternalServerError)
		return
	}

	data := struct {
		Pedido     models.Pedido
		Detalles   []models.DetallePedido
		LoginToken bool
		Perfil     string
	}{
		Pedido:     pedido,
		Detalles:   detalles,
		LoginToken: loggedIn,
		Perfil:     perfil,
	}

	tmpl.ExecuteTemplate(w, "base", data)
}

package handlers

import (
	"Go-Sistemas-de-Gestion-empresarial/db"
	"Go-Sistemas-de-Gestion-empresarial/models"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	_, perfil, _ := GetSessionData(r)

	stats, err := getAdminStats()
	if err != nil {
		log.Println("Error obteniendo estadísticas:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/admin/layout.html", "templates/admin/dashboard.html")
	if err != nil {
		log.Println("Error cargando templates admin:", err)
		http.Error(w, "Error cargando templates", http.StatusInternalServerError)
		return
	}

	data := struct {
		Perfil          string
		Stats           AdminStats
		DashboardActive bool
		ProductosActive bool
		PedidosActive   bool
		ClientesActive  bool
	}{
		Perfil:          perfil,
		Stats:           stats,
		DashboardActive: true,
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}

type AdminStats struct {
	TotalClientes  int
	TotalProductos int
	TotalPedidos   int
	TotalVentas    float64
}

type sqlNullFloat64 struct {
	sql.NullFloat64
}

func getAdminStats() (AdminStats, error) {
	var stats AdminStats
	database, err := db.Connect()
	if err != nil {
		return stats, err
	}
	defer database.Close()

	database.QueryRow("SELECT COUNT(*) FROM clientes").Scan(&stats.TotalClientes)

	database.QueryRow("SELECT COUNT(*) FROM productos").Scan(&stats.TotalProductos)

	database.QueryRow("SELECT COUNT(*) FROM pedidos").Scan(&stats.TotalPedidos)

	var total sql.NullFloat64
	database.QueryRow("SELECT SUM(total) FROM pedidos").Scan(&total)
	stats.TotalVentas = total.Float64

	return stats, nil
}

func AdminProducts(w http.ResponseWriter, r *http.Request) {
	_, perfil, _ := GetSessionData(r)

	productos, err := models.GetAllProductos()
	if err != nil {
		log.Println("Error obteniendo productos:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/admin/layout.html", "templates/admin/products_list.html")
	if err != nil {
		log.Println("Error cargando templates admin products:", err)
		http.Error(w, "Error cargando templates", http.StatusInternalServerError)
		return
	}

	data := struct {
		Perfil          string
		Productos       []models.Producto
		DashboardActive bool
		ProductosActive bool
		PedidosActive   bool
		ClientesActive  bool
	}{
		Perfil:          perfil,
		Productos:       productos,
		ProductosActive: true,
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}

func AdminProductCreate(w http.ResponseWriter, r *http.Request) {
	_, perfil, _ := GetSessionData(r)

	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		descripcion := r.FormValue("descripcion")
		precio, _ := strconv.ParseFloat(r.FormValue("precio"), 64)
		stock, _ := strconv.Atoi(r.FormValue("stock"))
		sku := r.FormValue("sku")
		activo := r.FormValue("activo") == "on"

		err := models.CreateProducto(nombre, descripcion, precio, stock, sku, activo)
		if err != nil {
			log.Println("Error creando producto:", err)
			http.Error(w, "Error creando producto", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/productos", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/admin/layout.html", "templates/admin/product_form.html")
	if err != nil {
		log.Println("Error cargando template admin product form:", err)
		http.Error(w, "Error cargando templates", http.StatusInternalServerError)
		return
	}

	data := struct {
		Perfil          string
		IsEdit          bool
		ProductosActive bool
		DashboardActive bool
		PedidosActive   bool
		ClientesActive  bool
	}{
		Perfil:          perfil,
		IsEdit:          false,
		ProductosActive: true,
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}

func AdminProductEdit(w http.ResponseWriter, r *http.Request) {
	_, perfil, _ := GetSessionData(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "POST":
		nombre := r.FormValue("nombre")
		descripcion := r.FormValue("descripcion")
		precio, _ := strconv.ParseFloat(r.FormValue("precio"), 64)
		stock, _ := strconv.Atoi(r.FormValue("stock"))
		sku := r.FormValue("sku")
		activo := r.FormValue("activo") == "on"

		err := models.UpdateProducto(id, nombre, descripcion, precio, stock, sku, activo)
		if err != nil {
			log.Println("Error actualizando producto:", err)
			http.Error(w, "Error actualizando producto", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/productos", http.StatusSeeOther)

	case "GET":
		producto, err := models.GetProductoByID(id)
		if err != nil {
			http.Error(w, "Producto no encontrado", http.StatusNotFound)
			return
		}

		tmpl, err := template.ParseFiles("templates/admin/layout.html", "templates/admin/product_form.html")
		if err != nil {
			log.Println("Error cargando template admin product form:", err)
			http.Error(w, "Error cargando templates", http.StatusInternalServerError)
			return
		}

		data := struct {
			Perfil          string
			IsEdit          bool
			Producto        models.Producto
			DashboardActive bool
			ProductosActive bool
			PedidosActive   bool
			ClientesActive  bool
		}{
			Perfil:          perfil,
			IsEdit:          true,
			Producto:        producto,
			ProductosActive: true,
		}

		err = tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			log.Println("Error ejecutando template admin product form:", err)
		}

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func AdminProductDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := models.DeleteProducto(id)
	if err != nil {
		log.Println("Error eliminando producto:", err)
	}
	http.Redirect(w, r, "/admin/productos", http.StatusSeeOther)
}

func AdminOrders(w http.ResponseWriter, r *http.Request) {
	_, perfil, _ := GetSessionData(r)

	pedidos, err := models.GetAllPedidos()
	if err != nil {
		log.Println("Error obteniendo pedidos:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/admin/layout.html", "templates/admin/orders_list.html")
	if err != nil {
		log.Println("Error cargando templates admin orders:", err)
		http.Error(w, "Error cargando templates", http.StatusInternalServerError)
		return
	}

	data := struct {
		Perfil          string
		Pedidos         []models.Pedido
		DashboardActive bool
		ProductosActive bool
		PedidosActive   bool
		ClientesActive  bool
	}{
		Perfil:        perfil,
		Pedidos:       pedidos,
		PedidosActive: true,
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Println("Error ejecutando template admin orders:", err)
	}
}

func AdminOrderDetail(w http.ResponseWriter, r *http.Request) {
	_, perfil, _ := GetSessionData(r)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	pedido, err := models.GetPedidoByID(id)
	if err != nil {
		http.Error(w, "Pedido no encontrado", http.StatusNotFound)
		return
	}

	detalles, err := models.GetDetallesByPedidoID(id)
	if err != nil {
		log.Println("Error obteniendo detalles del pedido:", err)
	}

	cliente, err := models.GetClienteByID(pedido.IDCliente)
	if err != nil {
		log.Println("Error obteniendo cliente:", err)
	}

	tmpl, err := template.ParseFiles("templates/admin/layout.html", "templates/admin/order_detail.html")
	if err != nil {
		log.Println("Error cargando template admin order detail:", err)
		http.Error(w, "Error cargando templates", http.StatusInternalServerError)
		return
	}

	data := struct {
		Perfil          string
		Pedido          models.Pedido
		Detalles        []models.DetallePedido
		Cliente         models.Cliente
		DashboardActive bool
		ProductosActive bool
		PedidosActive   bool
		ClientesActive  bool
	}{
		Perfil:        perfil,
		Pedido:        pedido,
		Detalles:      detalles,
		Cliente:       cliente,
		PedidosActive: true,
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Println("Error ejecutando template admin order detail:", err)
	}
}

func AdminClients(w http.ResponseWriter, r *http.Request) {
	_, perfil, _ := GetSessionData(r)

	clientes, err := models.GetAllClientes()
	if err != nil {
		log.Println("Error obteniendo clientes:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/admin/layout.html", "templates/admin/clients_list.html")
	if err != nil {
		log.Println("Error cargando templates admin clients:", err)
		http.Error(w, "Error cargando templates", http.StatusInternalServerError)
		return
	}

	data := struct {
		Perfil          string
		Clientes        []models.Cliente
		DashboardActive bool
		ProductosActive bool
		PedidosActive   bool
		ClientesActive  bool
	}{
		Perfil:         perfil,
		Clientes:       clientes,
		ClientesActive: true,
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Println("Error ejecutando template admin clients:", err)
	}
}

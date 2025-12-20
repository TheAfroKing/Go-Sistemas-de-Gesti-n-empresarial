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
	// AdminDashboard muestra el dashboard de administración con estadísticas generales.
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

// sqlNullFloat64 es un wrapper para sql.NullFloat64 usado internamente.

type sqlNullFloat64 struct {
	sql.NullFloat64
}

func getAdminStats() (AdminStats, error) {
	// getAdminStats obtiene estadísticas agregadas de la base de datos para el admin.
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
	// AdminProducts lista todos los productos en la vista de administración.
	_, perfil, _ := GetSessionData(r)

	productos, err := models.GetAllProductos()
	if err != nil {
		log.Println("Error obteniendo productos:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/admin/layout.html", "templates/admin/productos.html")
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
	// AdminProductCreate maneja la creación de productos desde el panel admin.
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

	tmpl, err := template.ParseFiles("templates/admin/layout.html", "templates/admin/formulario_producto.html")
	if err != nil {
		log.Println("Error cargando template admin product form:", err)
		http.Error(w, "Error cargando templates", http.StatusInternalServerError)
		return
	}

	data := struct {
		Perfil          string
		IsEdit          bool
		Producto        models.Producto
		ProductosActive bool
		DashboardActive bool
		PedidosActive   bool
		ClientesActive  bool
	}{
		Perfil:          perfil,
		IsEdit:          false,
		Producto:        models.Producto{},
		ProductosActive: true,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println("Error ejecutando template admin product form:", err)
		http.Error(w, "Error ejecutando template", http.StatusInternalServerError)
		return
	}
}

func AdminProductEdit(w http.ResponseWriter, r *http.Request) {
	// AdminProductEdit permite editar un producto existente o mostrar el formulario.
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

		tmpl, err := template.ParseFiles("templates/admin/layout.html", "templates/admin/formulario_producto.html")
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
	// AdminProductDelete elimina un producto por su ID y redirige a la lista.
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := models.DeleteProducto(id)
	if err != nil {
		log.Println("Error eliminando producto:", err)
	}
	http.Redirect(w, r, "/admin/productos", http.StatusSeeOther)
}

func AdminOrders(w http.ResponseWriter, r *http.Request) {
	// AdminOrders lista todos los pedidos en la vista de administración.
	_, perfil, _ := GetSessionData(r)

	pedidos, err := models.GetAllPedidos()
	if err != nil {
		log.Println("Error obteniendo pedidos:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/admin/layout.html", "templates/admin/ordenes.html")
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
	// AdminOrderDetail muestra los detalles de un pedido específico en admin.
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

	tmpl, err := template.ParseFiles("templates/admin/layout.html", "templates/admin/detalle_orden.html")
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
	// AdminClients lista todos los clientes en el panel de administración.
	_, perfil, _ := GetSessionData(r)

	clientes, err := models.GetAllClientes()
	if err != nil {
		log.Println("Error obteniendo clientes:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/admin/layout.html", "templates/admin/clientes.html")
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

func AdminOrderStatus(w http.ResponseWriter, r *http.Request) {
	// AdminOrderStatus actualiza el estado de un pedido (p.ej. PAGADO, ENTREGADO).
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if r.Method == "POST" {
		nuevoEstado := r.FormValue("estado") // PAGADO, ENTREGADO
		err := models.UpdatePedidoStatus(id, nuevoEstado)
		if err != nil {
			log.Println("Error actualizando estado del pedido:", err)
			http.Error(w, "Error actualizando estado", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/admin/pedidos", http.StatusSeeOther)
}

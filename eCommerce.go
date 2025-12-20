package main

import (
	"Go-Sistemas-de-Gestion-empresarial/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// main inicia el servidor web y registra las rutas principales del eCommerce.
// Usa Gorilla Mux para el enrutamiento y carga variables de entorno con godotenv.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Nota: No se pudo cargar el archivo .env, usando variables de entorno del sistema")
	}

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET")
	r.HandleFunc("/producto/{id:[0-9]+}", handlers.ClientProductDetail).Methods("GET")

	r.HandleFunc("/carrito", handlers.ClientCart).Methods("GET")
	r.HandleFunc("/producto/agregar-carrito", handlers.AgregarItemCarrito).Methods("POST")
	r.HandleFunc("/carrito/eliminar/{id:[0-9]+}", handlers.RemoveItemFromCart).Methods("GET", "POST") // Allow POST for better practice, but kept GET to match link
	r.HandleFunc("/checkout", handlers.ClientCheckout).Methods("GET")
	r.HandleFunc("/checkout", handlers.ProcessCheckout).Methods("POST")

	r.HandleFunc("/perfil", handlers.ClientProfile).Methods("GET")
	r.HandleFunc("/perfil/editar", handlers.ClientProfileEdit).Methods("GET", "POST")
	r.HandleFunc("/pedidos/{id:[0-9]+}", handlers.ClientOrderDetail).Methods("GET")

	r.HandleFunc("/admin/dashboard", handlers.AdminDashboard).Methods("GET")
	r.HandleFunc("/admin/productos", handlers.AdminProducts).Methods("GET")
	r.HandleFunc("/admin/productos/nuevo", handlers.AdminProductCreate).Methods("GET", "POST")
	r.HandleFunc("/admin/productos/editar/{id}", handlers.AdminProductEdit).Methods("GET", "POST")
	r.HandleFunc("/admin/productos/eliminar/{id}", handlers.AdminProductDelete)
	r.HandleFunc("/admin/pedidos", handlers.AdminOrders).Methods("GET")
	r.HandleFunc("/admin/pedidos/{id}", handlers.AdminOrderDetail).Methods("GET")
	r.HandleFunc("/admin/pedidos/{id}/status", handlers.AdminOrderStatus).Methods("POST")
	r.HandleFunc("/admin/clientes", handlers.AdminClients).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Servidor iniciado en puerto :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

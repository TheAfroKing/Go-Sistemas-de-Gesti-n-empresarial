/*
@Autor: Fabian Paredes
@Fecha: 07/12/2025
@Descripcion: Código para gestion de eCommerce
*/

package main

import (
	"errors"
	"fmt"
	"time"
	//"encoding/csv"
	//"os"
)

/*
		==========================================
		Clase cliente con sus propiedades
	   	==========================================
*/
type Cliente struct {
	idCliente int
	nombre    string
	email     string
	direccion string
}

// Getters y Setters de Cliente
func (c *Cliente) GetIdCliente() int    { return c.idCliente }
func (c *Cliente) GetNombre() string    { return c.nombre }
func (c *Cliente) GetEmail() string     { return c.email }
func (c *Cliente) GetDireccion() string { return c.direccion }

func (c *Cliente) SetIdCliente(id int)           { c.idCliente = id }
func (c *Cliente) SetNombre(nombre string)       { c.nombre = nombre }
func (c *Cliente) SetEmail(email string)         { c.email = email }
func (c *Cliente) SetDireccion(direccion string) { c.direccion = direccion }

// Métodos específicos del Diagrama
func (c *Cliente) ActualizarDireccion(nuevaDireccion string) {
	c.direccion = nuevaDireccion
	fmt.Printf("Dirección actualizada para %s.\n", c.nombre)
}

/*
		==========================================
		Clase producto con sus propiedades
	   	==========================================
*/
type Producto struct {
	codigo     int
	nombre     string
	precio     float64
	cantidad   float64
	categorias []string
}

// Getters y Setters
func (p *Producto) GetCodigo() int          { return p.codigo }
func (p *Producto) GetNombre() string       { return p.nombre }
func (p *Producto) GetPrecio() float64      { return p.precio }
func (p *Producto) GetCantidad() float64    { return p.cantidad }
func (p *Producto) GetCategorias() []string { return p.categorias }

func (p *Producto) SetCodigo(codigo int)         { p.codigo = codigo }
func (p *Producto) SetNombre(nombre string)      { p.nombre = nombre }
func (p *Producto) SetPrecio(precio float64)     { p.precio = precio }
func (p *Producto) SetCantidad(cantidad float64) { p.cantidad = cantidad }
func (p *Producto) SetCategorias(cats []string)  { p.categorias = cats }

func (p *Producto) RevisarStock() bool {
	return p.cantidad > 0
}

func (p *Producto) QuitarStock(cantidad float64) error {
	if p.cantidad >= cantidad {
		p.cantidad -= cantidad
		return nil
	}
	return errors.New("Stock insuficiente")
}

/*
		==========================================
		Clase carrito con sus propiedades
	   	==========================================
*/
type Carrito struct {
	idCarrito int
	idCliente int
	fecha     string
	productos []*Producto
}

// Getters y Setters
func (c *Carrito) GetIdCarrito() int         { return c.idCarrito }
func (c *Carrito) GetIdCliente() int         { return c.idCliente }
func (c *Carrito) GetFecha() string          { return c.fecha }
func (c *Carrito) GetProductos() []*Producto { return c.productos }

func (c *Carrito) SetIdCarrito(id int)        { c.idCarrito = id }
func (c *Carrito) SetIdCliente(idCliente int) { c.idCliente = idCliente }
func (c *Carrito) SetFecha(fecha string)      { c.fecha = fecha }

func (c *Carrito) AgregarProductoCarrito(prod *Producto) {
	c.productos = append(c.productos, prod)
	fmt.Printf("Producto %s agregado al carrito.\n", prod.nombre)
}

func (c *Carrito) QuitarProductoCarrito(codigoProducto int) {
	for i, p := range c.productos {
		if p.GetCodigo() == codigoProducto {
			c.productos = append(c.productos[:i], c.productos[i+1:]...)
			fmt.Println("Producto eliminado del carrito.")
			return
		}
	}
	fmt.Println("Producto no encontrado en el carrito.")
}

func (c *Carrito) CalcularTotal() float64 {
	var total float64 = 0
	for _, p := range c.productos {
		total += p.GetPrecio()
	}
	return total
}

func (c *Carrito) CancelarCarrito() {
	c.productos = []*Producto{} // Vaciar slice
	fmt.Println("Carrito vaciado.")
}

/*
		==========================================
		Clase pedido con sus propiedades
	   	==========================================
*/
type Pedido struct {
	idPedido  int
	idCliente int
	productos []*Producto
	fecha     string
	estado    string
	precio    float64
}

// Getters y Setters
func (p *Pedido) GetIdPedido() int          { return p.idPedido }
func (p *Pedido) GetIdCliente() int         { return p.idCliente }
func (p *Pedido) GetProductos() []*Producto { return p.productos }
func (p *Pedido) GetFecha() string          { return p.fecha }
func (p *Pedido) GetEstado() string         { return p.estado }
func (p *Pedido) GetPrecio() float64        { return p.precio }

func (p *Pedido) SetIdPedido(id int)                 { p.idPedido = id }
func (p *Pedido) SetIdCliente(idCliente int)         { p.idCliente = idCliente }
func (p *Pedido) SetProductos(productos []*Producto) { p.productos = productos }
func (p *Pedido) SetFecha(fecha string)              { p.fecha = fecha }
func (p *Pedido) SetEstado(estado string)            { p.estado = estado }
func (p *Pedido) SetPrecio(precio float64)           { p.precio = precio }

func (p *Pedido) PagarPedido() {
	p.estado = "PAGADO"
	fmt.Printf("El pedido %d ha sido pagado. Total: %.2f\n", p.idPedido, p.precio)
}

func (p *Pedido) EntregarPedido() {
	if p.estado == "PAGADO" {
		p.estado = "ENTREGADO"
		fmt.Printf("Pedido %d entregado con éxito.\n", p.idPedido)
	} else {
		fmt.Println("No se puede entregar un pedido no pagado.")
	}
}

func (p *Pedido) CancelarPedido() {
	p.estado = "CANCELADO"
	fmt.Printf("Pedido %d cancelado.\n", p.idPedido)
}

/*
		==========================================
		Clase tienda con sus propiedades
	   	==========================================
*/
type Tienda struct {
	clientes  []*Cliente
	productos []*Producto
	carritos  []*Carrito
	pedidos   []*Pedido
	lastId    int
}

// Métodos de control (simulando la base de datos de la tienda)

func (t *Tienda) AgregarCliente(nombre, email, direccion string) error {
	if nombre == "" || email == "" {
		return errors.New("Datos del cliente incompletos")
	}
	t.lastId++
	nuevoCliente := &Cliente{}
	nuevoCliente.SetIdCliente(t.lastId)
	nuevoCliente.SetNombre(nombre)
	nuevoCliente.SetEmail(email)
	nuevoCliente.SetDireccion(direccion)
	t.clientes = append(t.clientes, nuevoCliente)
	return nil
}

// Equivalente a Productos.agregar() del diagrama
func (t *Tienda) AgregarProducto(nombre string, precio, cantidad float64, categorias ...string) error {
	if nombre == "" || precio <= 0 {
		return errors.New("Datos del producto inválidos")
	}
	t.lastId++
	nuevoProd := &Producto{}
	nuevoProd.SetCodigo(t.lastId)
	nuevoProd.SetNombre(nombre)
	nuevoProd.SetPrecio(precio)
	nuevoProd.SetCantidad(cantidad)
	nuevoProd.SetCategorias(categorias)
	t.productos = append(t.productos, nuevoProd)
	return nil
}

// Gestionar lógica de crear carrito o recuperar existente
func (t *Tienda) ObtenerCarrito(idCliente int) *Carrito {
	for _, c := range t.carritos {
		if c.GetIdCliente() == idCliente {
			return c
		}
	}
	// Si no existe, crea uno
	t.lastId++
	nuevoCarrito := &Carrito{}
	nuevoCarrito.SetIdCarrito(t.lastId)
	nuevoCarrito.SetIdCliente(idCliente)
	nuevoCarrito.SetFecha(time.Now().Format("2006-01-02"))
	t.carritos = append(t.carritos, nuevoCarrito)
	return nuevoCarrito
}

// Convertir Carrito a Pedido
func (t *Tienda) CrearPedido(idCliente int) error {
	carrito := t.ObtenerCarrito(idCliente)
	if len(carrito.GetProductos()) == 0 {
		return errors.New("El carrito está vacío")
	}

	t.lastId++
	nuevoPedido := &Pedido{
		idPedido:  t.lastId,
		idCliente: idCliente,
		productos: carrito.GetProductos(), // Copia los productos
		fecha:     time.Now().Format("2006-01-02"),
		estado:    "PENDIENTE",
		precio:    carrito.CalcularTotal(),
	}

	// Verificar Stock antes de confirmar
	for _, prod := range nuevoPedido.productos {
		if err := prod.QuitarStock(1); err != nil {
			return fmt.Errorf("Error de stock en producto %s", prod.nombre)
		}
	}

	t.pedidos = append(t.pedidos, nuevoPedido)
	carrito.CancelarCarrito()
	fmt.Println("Pedido creado exitosamente.")
	return nil
}

// Listados
func (t *Tienda) ListarProductos() {
	fmt.Println("--- CATÁLOGO DE PRODUCTOS ---")
	for _, p := range t.productos {
		fmt.Printf("[ID: %d] %s - $%.2f (Stock: %.0f)\n", p.GetCodigo(), p.GetNombre(), p.GetPrecio(), p.GetCantidad())
	}
}

func (t *Tienda) ListarPedidos() {
	fmt.Println("--- LISTA DE PEDIDOS ---")
	for _, p := range t.pedidos {
		// Buscar nombre cliente
		nombreCliente := "Desconocido"
		for _, c := range t.clientes {
			if c.GetIdCliente() == p.idCliente {
				nombreCliente = c.GetNombre()
			}
		}
		fmt.Printf("Pedido #%d | Cliente: %s | Estado: %s | Total: $%.2f\n", p.idPedido, nombreCliente, p.estado, p.precio)
	}
}

func (t *Tienda) ListarClientes() {
	fmt.Println("--- LISTA DE CLIENTES ---")
	if len(t.clientes) == 0 {
		fmt.Println("No hay clientes registrados.")
		return
	}
	for _, c := range t.clientes {
		fmt.Printf("[ID: %d] %s | Email: %s | Dirección: %s\n",
			c.GetIdCliente(), c.GetNombre(), c.GetEmail(), c.GetDireccion())
	}
}

/*
		==========================================
		MAIN MENU
	   	==========================================
*/
func main() {
	tienda := &Tienda{}

	// Datos precargados para pruebas rápidas
	tienda.AgregarProducto("Laptop Gamer", 1500.00, 10, "Tecnología")
	tienda.AgregarProducto("Mouse Inalámbrico", 25.50, 50, "Accesorios")
	tienda.AgregarCliente("Juan Perez", "juan@mail.com", "Av. Central 123")

	for {
		fmt.Println("\n=== SISTEMA E-COMMERCE ===")
		fmt.Println("--- GESTIÓN ---")
		fmt.Println("1. Registrar Cliente")
		fmt.Println("2. Agregar Nuevo Producto")
		fmt.Println("--- COMPRAS ---")
		fmt.Println("3. Comprar (Agregar al Carrito)")
		fmt.Println("4. Ver Mi Carrito y Total")
		fmt.Println("5. Confirmar Pedido (Checkout)")
		fmt.Println("6. Pagar/Entregar Pedidos (Admin)")
		fmt.Println("--- REPORTES Y LISTADOS ---")
		fmt.Println("7. Listar Clientes")
		fmt.Println("8. Listar Productos")
		fmt.Println("9. Listar Pedidos")
		fmt.Println("10. Salir")
		fmt.Print("Ingrese una opción: ")

		var opcion int
		fmt.Scan(&opcion)

		switch opcion {
		case 1:
			var nombre, email, direccion string
			fmt.Println("--- REGISTRO CLIENTE ---")
			fmt.Print("Nombre: ")
			fmt.Scan(&nombre)
			fmt.Print("Email: ")
			fmt.Scan(&email)
			fmt.Print("Dirección: ")
			fmt.Scan(&direccion) // Nota: Scan corta en espacios, usar bufio para frases completas es mejor, pero mantenemos Scan por simplicidad aquí
			tienda.AgregarCliente(nombre, email, direccion)
			fmt.Println("Cliente registrado.")

		case 2:
			var nombre string
			var precio, cantidad float64
			fmt.Println("--- NUEVO PRODUCTO ---")
			fmt.Print("Nombre: ")
			fmt.Scan(&nombre)
			fmt.Print("Precio: ")
			fmt.Scan(&precio)
			fmt.Print("Cantidad inicial: ")
			fmt.Scan(&cantidad)
			tienda.AgregarProducto(nombre, precio, cantidad, "General")
			fmt.Println("Producto agregado.")

		case 3:
			var idCliente, idProd int
			fmt.Println("--- COMPRAR ---")
			fmt.Print("ID Cliente: ")
			fmt.Scan(&idCliente)
			tienda.ListarProductos()
			fmt.Print("ID Producto a comprar: ")
			fmt.Scan(&idProd)

			// Buscar producto
			var prodEncontrado *Producto
			for _, p := range tienda.productos {
				if p.GetCodigo() == idProd {
					prodEncontrado = p
					break
				}
			}

			if prodEncontrado != nil && prodEncontrado.RevisarStock() {
				carrito := tienda.ObtenerCarrito(idCliente)
				carrito.AgregarProductoCarrito(prodEncontrado)
			} else {
				fmt.Println("Producto no encontrado o sin stock.")
			}

		case 4:
			var idCliente int
			fmt.Print("ID Cliente para ver carrito: ")
			fmt.Scan(&idCliente)
			carrito := tienda.ObtenerCarrito(idCliente)
			fmt.Printf("--- CARRITO DE CLIENTE %d ---\n", idCliente)
			if len(carrito.GetProductos()) == 0 {
				fmt.Println("El carrito está vacío.")
			} else {
				for _, p := range carrito.GetProductos() {
					fmt.Printf("- %s ($%.2f)\n", p.GetNombre(), p.GetPrecio())
				}
				fmt.Printf("TOTAL A PAGAR: $%.2f\n", carrito.CalcularTotal())
			}

		case 5:
			var idCliente int
			fmt.Print("ID Cliente para Checkout: ")
			fmt.Scan(&idCliente)
			if err := tienda.CrearPedido(idCliente); err != nil {
				fmt.Println("Error al crear pedido:", err)
			}

		case 6:
			tienda.ListarPedidos()
			var idPedido int
			var accion int
			fmt.Print("Ingrese ID Pedido a gestionar: ")
			fmt.Scan(&idPedido)
			fmt.Print("Acción (1. Pagar | 2. Entregar): ")
			fmt.Scan(&accion)

			found := false
			for _, p := range tienda.pedidos {
				if p.GetIdPedido() == idPedido {
					found = true
					if accion == 1 {
						p.PagarPedido()
					} else if accion == 2 {
						p.EntregarPedido()
					} else {
						fmt.Println("Acción no válida.")
					}
				}
			}
			if !found {
				fmt.Println("Pedido no encontrado.")
			}

		case 7:
			tienda.ListarClientes()

		case 8:
			tienda.ListarProductos()

		case 9:
			tienda.ListarPedidos()

		case 10:
			fmt.Println("Saliendo del sistema...")
			return
		default:
			fmt.Println("Opción no válida, intente de nuevo.")
		}
	}
}

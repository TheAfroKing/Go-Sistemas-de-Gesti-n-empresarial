/*
@Autor: Fabian Paredes
@Fecha: 07/12/2025
@Descripcion: Código para gestion de eCommerce
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
==========================================
DEFINICIÓN DE ERRORES PERSONALIZADOS
==========================================
*/
var (
	ErrStockInsuficiente    = errors.New("error: stock insuficiente para el producto solicitado")
	ErrCarritoVacio         = errors.New("error: no se puede procesar un pedido con el carrito vacío")
	ErrProductoNoEncontrado = errors.New("error: producto no encontrado")
	ErrDatoInvalido         = errors.New("error: datos de entrada inválidos")
)

/*
==========================================
INTERFACES
==========================================
Interfaz MetodoPago:
Define el comportamiento que debe tener cualquier forma de pago.
Esto permite que el sistema sea escalable (podemos agregar PayPal, Crypto, etc. sin romper el código).
*/
type MetodoPago interface {
	ProcesarPago(monto float64) error
	NombreMetodo() string
}

// Implementación 1: Pago en Efectivo
type PagoEfectivo struct{}

func (p PagoEfectivo) ProcesarPago(monto float64) error {
	fmt.Printf("Recibiendo %.2f en efectivo...\n", monto)
	return nil
}
func (p PagoEfectivo) NombreMetodo() string { return "Efectivo" }

// Implementación 2: Pago con Tarjeta
type PagoTarjeta struct {
	Numero string
}

func (p PagoTarjeta) ProcesarPago(monto float64) error {
	if len(p.Numero) < 4 {
		return errors.New("número de tarjeta inválido")
	}
	fmt.Printf("Procesando cobro de %.2f a la tarjeta terminada en *%s...\n", monto, p.Numero[len(p.Numero)-4:])
	return nil
}

func (p PagoTarjeta) NombreMetodo() string { return "Tarjeta de Crédito" }

/*
==========================================
Clase cliente
==========================================
*/
type Cliente struct {
	idCliente int
	nombre    string
	email     string
	direccion string
}

// Getters y Setters
func (c *Cliente) GetIdCliente() int    { return c.idCliente }
func (c *Cliente) GetNombre() string    { return c.nombre }
func (c *Cliente) GetEmail() string     { return c.email }
func (c *Cliente) GetDireccion() string { return c.direccion }

func (c *Cliente) SetIdCliente(id int)           { c.idCliente = id }
func (c *Cliente) SetNombre(nombre string)       { c.nombre = nombre }
func (c *Cliente) SetEmail(email string)         { c.email = email }
func (c *Cliente) SetDireccion(direccion string) { c.direccion = direccion }

/*
==========================================
Clase producto
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
func (p *Producto) GetCodigo() int       { return p.codigo }
func (p *Producto) GetNombre() string    { return p.nombre }
func (p *Producto) GetPrecio() float64   { return p.precio }
func (p *Producto) GetCantidad() float64 { return p.cantidad }

func (p *Producto) SetCodigo(codigo int)         { p.codigo = codigo }
func (p *Producto) SetNombre(nombre string)      { p.nombre = nombre }
func (p *Producto) SetPrecio(precio float64)     { p.precio = precio }
func (p *Producto) SetCantidad(cantidad float64) { p.cantidad = cantidad }
func (p *Producto) SetCategorias(cats []string)  { p.categorias = cats }

func (p *Producto) RevisarStock() bool {
	return p.cantidad > 0
}

// QuitarStock ahora devuelve nuestro error personalizado
func (p *Producto) QuitarStock(cantidad float64) error {
	if p.cantidad >= cantidad {
		p.cantidad -= cantidad
		return nil
	}
	return ErrStockInsuficiente
}

/*
==========================================
Clase carrito
==========================================
*/
type Carrito struct {
	idCarrito int
	idCliente int
	fecha     string
	productos []*Producto
}

func (c *Carrito) GetIdCliente() int         { return c.idCliente }
func (c *Carrito) GetProductos() []*Producto { return c.productos }
func (c *Carrito) SetIdCarrito(id int)       { c.idCarrito = id }
func (c *Carrito) SetIdCliente(id int)       { c.idCliente = id }
func (c *Carrito) SetFecha(f string)         { c.fecha = f }

func (c *Carrito) AgregarProductoCarrito(prod *Producto) {
	c.productos = append(c.productos, prod)
	fmt.Printf(">> Producto '%s' agregado al carrito.\n", prod.nombre)
}

func (c *Carrito) CalcularTotal() float64 {
	var total float64 = 0
	for _, p := range c.productos {
		total += p.precio
	}
	return total
}

func (c *Carrito) CancelarCarrito() {
	c.productos = []*Producto{}
}

/*
==========================================
Clase pedido
==========================================
*/
type Pedido struct {
	idPedido    int
	idCliente   int
	productos   []*Producto
	fecha       string
	estado      string
	precio      float64
	detallePago string
}

// Getters basicos
func (p *Pedido) GetIdPedido() int   { return p.idPedido }
func (p *Pedido) GetEstado() string  { return p.estado }
func (p *Pedido) GetPrecio() float64 { return p.precio }

// PagarPedido recibe una INTERFAZ MetodoPago.
// Esto permite pagar con cualquier struct que implemente "ProcesarPago".
func (p *Pedido) PagarPedido(metodo MetodoPago) error {
	if p.estado == "PAGADO" {
		return errors.New("el pedido ya está pagado")
	}

	err := metodo.ProcesarPago(p.precio)
	if err != nil {
		return fmt.Errorf("fallo al procesar el pago: %w", err)
	}
	p.estado = "PAGADO"
	p.detallePago = metodo.NombreMetodo()

	fmt.Printf("El pedido %d ha sido pagado usando %s.\n", p.idPedido, metodo.NombreMetodo())
	return nil
}

func (p *Pedido) EntregarPedido() {
	if p.estado == "PAGADO" {
		p.estado = "ENTREGADO"
		fmt.Printf("Pedido %d entregado con éxito.\n", p.idPedido)
	} else {
		fmt.Println("ERROR: No se puede entregar un pedido no pagado.")
	}
}

/*
==========================================
Clase tienda (Controlador Principal)
==========================================
*/
type Tienda struct {
	clientes  []*Cliente
	productos []*Producto
	carritos  []*Carrito
	pedidos   []*Pedido
	lastId    int
}

func (t *Tienda) AgregarCliente(nombre, email, direccion string) error {
	if nombre == "" || email == "" {
		return ErrDatoInvalido
	}
	t.lastId++
	nuevoCliente := &Cliente{
		idCliente: t.lastId,
		nombre:    nombre,
		email:     email,
		direccion: direccion,
	}
	t.clientes = append(t.clientes, nuevoCliente)
	return nil
}

func (t *Tienda) AgregarProducto(nombre string, precio, cantidad float64, categorias ...string) error {
	if nombre == "" || precio <= 0 {
		return ErrDatoInvalido
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

func (t *Tienda) ObtenerCarrito(idCliente int) *Carrito {
	for _, c := range t.carritos {
		if c.GetIdCliente() == idCliente {
			return c
		}
	}
	t.lastId++
	nuevoCarrito := &Carrito{}
	nuevoCarrito.SetIdCarrito(t.lastId)
	nuevoCarrito.SetIdCliente(idCliente)
	nuevoCarrito.SetFecha(time.Now().Format("2006-01-02"))
	t.carritos = append(t.carritos, nuevoCarrito)
	return nuevoCarrito
}

/*
CrearPedido:
Esta es una función crítica. Maneja la lógica de "Checkout".
1. Obtiene el carrito.
2. Verifica que no esté vacío.
3. Intenta reservar el stock de CADA producto.
  - Si falla el stock de un producto, debemos decidir qué hacer (aquí retornamos error).

4. Si todo sale bien, crea el objeto Pedido y limpia el carrito.
*/
func (t *Tienda) CrearPedido(idCliente int) error {
	carrito := t.ObtenerCarrito(idCliente)
	if len(carrito.GetProductos()) == 0 {
		return ErrCarritoVacio
	}

	// Verificar Stock antes de confirmar
	// Recorremos los productos del carrito e intentamos descontar stock
	for _, prod := range carrito.GetProductos() {
		// Intentamos quitar 1 unidad
		if err := prod.QuitarStock(1); err != nil {
			// Si falla, retornamos un error envuelto con contexto
			return fmt.Errorf("no se pudo procesar el producto '%s': %w", prod.nombre, err)
		}
	}

	t.lastId++
	nuevoPedido := &Pedido{
		idPedido:  t.lastId,
		idCliente: idCliente,
		productos: carrito.GetProductos(),
		fecha:     time.Now().Format("2006-01-02"),
		estado:    "PENDIENTE",
		precio:    carrito.CalcularTotal(),
	}

	t.pedidos = append(t.pedidos, nuevoPedido)
	carrito.CancelarCarrito()
	fmt.Println("¡Pedido creado exitosamente!")
	return nil
}

// Funciones de listado (Simplificadas para usar en menú)
func (t *Tienda) ListarProductos() {
	fmt.Println("\n--- CATÁLOGO ---")
	for _, p := range t.productos {
		fmt.Printf("[ID: %d] %-20s | $%.2f | Stock: %.0f\n", p.codigo, p.nombre, p.precio, p.cantidad)
	}
}

func (t *Tienda) ListarPedidos() {
	fmt.Println("\n--- PEDIDOS ---")
	for _, p := range t.pedidos {
		if p.estado == "PENDIENTE" {
			fmt.Printf("Pedido #%d | Cliente ID: %d | Estado: %s | Total: $%.2f\n", p.idPedido, p.idCliente, p.estado, p.precio)
		} else {
			fmt.Printf("Pedido #%d | Cliente ID: %d | Estado: %s| Metodo de pago:  %s | Total: $%.2f\n", p.idPedido, p.idCliente, p.estado, p.detallePago, p.precio)
		}
	}
}

func (t *Tienda) ListarClientes() {
	fmt.Println("\n--- CLIENTES ---")
	for _, c := range t.clientes {
		fmt.Printf("[ID: %d] %s (Email: %s)\n", c.idCliente, c.nombre, c.email)
	}
}

/*
==========================================
HELPERS DE ENTRADA (Manejo de Errores)
==========================================
*/

// leerTexto usa bufio para leer líneas completas (incluyendo espacios)
// Esto soluciona el problema de fmt.Scan cortando nombres compuestos.
func leerTexto(mensaje string, reader *bufio.Reader) string {
	fmt.Print(mensaje)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// leerEntero lee un número y maneja el error si el usuario escribe letras
func leerEntero(mensaje string, reader *bufio.Reader) int {
	for {
		fmt.Print(mensaje)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		numero, err := strconv.Atoi(text)
		if err == nil {
			return numero
		}
		fmt.Println(">> Error: Por favor ingrese un número válido.")
	}
}

func leerFloat(mensaje string, reader *bufio.Reader) float64 {
	for {
		fmt.Print(mensaje)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		numero, err := strconv.ParseFloat(text, 64)
		if err == nil {
			return numero
		}
		fmt.Println(">> Error: Por favor ingrese un valor decimal válido (ej. 10.50).")
	}
}

/*
==========================================
MAIN MENU
==========================================
*/
func main() {
	tienda := &Tienda{}
	reader := bufio.NewReader(os.Stdin)

	tienda.AgregarProducto("Laptop Gamer", 1500.00, 10, "Tecnología")
	tienda.AgregarProducto("Mouse Inalámbrico", 25.50, 50, "Accesorios")
	tienda.AgregarCliente("Juan Perez", "juan@mail.com", "Av. Central 123")

	for {
		fmt.Println("\n=== SISTEMA E-COMMERCE (MEJORADO) ===")
		fmt.Println("1. Registrar Cliente")
		fmt.Println("2. Agregar Nuevo Producto")
		fmt.Println("3. Comprar (Agregar al Carrito)")
		fmt.Println("4. Ver Mi Carrito")
		fmt.Println("5. Confirmar Pedido (Checkout)")
		fmt.Println("6. Pagar Pedido (Usando Interfaces)")
		fmt.Println("7. Entregar Pedido (Admin)")
		fmt.Println("8. Reportes (Listar Todo)")
		fmt.Println("9. Salir")

		opcion := leerEntero("Ingrese una opción: ", reader)

		switch opcion {
		case 1:
			n := leerTexto("Nombre: ", reader)
			e := leerTexto("Email: ", reader)
			d := leerTexto("Dirección: ", reader)
			if err := tienda.AgregarCliente(n, e, d); err != nil {
				fmt.Println(">> Error al registrar:", err)
			} else {
				fmt.Println(">> Cliente registrado.")
			}

		case 2:
			n := leerTexto("Nombre Producto: ", reader)
			p := leerFloat("Precio: ", reader)
			c := leerFloat("Cantidad inicial: ", reader)
			tienda.AgregarProducto(n, p, c, "General")
			fmt.Println(">> Producto agregado.")

		case 3:
			tienda.ListarClientes()
			idCliente := leerEntero("ID Cliente: ", reader)
			tienda.ListarProductos()
			idProd := leerEntero("ID Producto a comprar: ", reader)

			// Busqueda manual básica
			var prodEncontrado *Producto
			for _, p := range tienda.productos {
				if p.GetCodigo() == idProd {
					prodEncontrado = p
					break
				}
			}

			if prodEncontrado != nil {
				if prodEncontrado.RevisarStock() {
					carrito := tienda.ObtenerCarrito(idCliente)
					carrito.AgregarProductoCarrito(prodEncontrado)
				} else {
					fmt.Println(">> Error: Producto sin stock disponible.")
				}
			} else {
				fmt.Println(">> Error: Producto no encontrado.")
			}

		case 4:
			tienda.ListarClientes()
			idCliente := leerEntero("ID Cliente: ", reader)
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
			tienda.ListarClientes()
			idCliente := leerEntero("ID Cliente para Checkout: ", reader)
			// Manejo de errores explícito al crear pedido
			if err := tienda.CrearPedido(idCliente); err != nil {
				fmt.Printf(">> NO SE PUDO CREAR EL PEDIDO: %s\n", err.Error())
			}

		case 6: // Pagar usando Interfaces
			tienda.ListarPedidos()
			idPedido := leerEntero("Ingrese ID Pedido a pagar: ", reader)

			// Buscar pedido
			var pedido *Pedido
			for _, p := range tienda.pedidos {
				if p.idPedido == idPedido {
					pedido = p
					break
				}
			}

			if pedido == nil {
				fmt.Println(">> Pedido no encontrado.")
				continue
			}

			fmt.Println("Seleccione método de pago:")
			fmt.Println("1. Efectivo")
			fmt.Println("2. Tarjeta de Crédito")
			tipoPago := leerEntero("Opción: ", reader)

			var metodo MetodoPago // Declaración de la Interfaz

			if tipoPago == 1 {
				metodo = PagoEfectivo{}
			} else if tipoPago == 2 {
				num := leerTexto("Ingrese num tarjeta (4 dígitos min): ", reader)
				metodo = PagoTarjeta{Numero: num}
			} else {
				fmt.Println("Método no válido")
				continue
			}

			// Polimorfismo: PagarPedido no sabe si es tarjeta o efectivo, solo llama a ProcesarPago
			if err := pedido.PagarPedido(metodo); err != nil {
				fmt.Println(">> Error en el pago:", err)
			}

		case 7:
			tienda.ListarPedidos()
			id := leerEntero("ID Pedido a entregar: ", reader)
			found := false
			for _, p := range tienda.pedidos {
				if p.idPedido == id {
					p.EntregarPedido()
					found = true
				}
			}
			if !found {
				fmt.Println("Pedido no encontrado.")
			}

		case 8:
			tienda.ListarClientes()
			tienda.ListarProductos()
			tienda.ListarPedidos()

		case 9:
			fmt.Println("Saliendo...")
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}

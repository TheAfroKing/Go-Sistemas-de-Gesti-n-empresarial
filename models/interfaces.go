package models

// ClienteRepository define la interfaz para el manejo de datos de clientes.
// Esto permite desacoplar la lógica de negocio de la implementación de base de datos.
type ClienteRepository interface {
	GetByID(id int) (Cliente, error)
	GetByEmail(email string) (Cliente, error)
	GetAll() ([]Cliente, error)
	Create(cliente Cliente) error
	Update(cliente Cliente) error
	Delete(id int) error
}

// ProductoRepository define la interfaz para el manejo de datos de productos.
type ProductoRepository interface {
	GetByID(id int) (Producto, error)
	GetAll() ([]Producto, error)
	Update(producto Producto) error
}

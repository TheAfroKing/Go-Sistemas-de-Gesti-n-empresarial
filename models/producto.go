package models

import (
	"Go-Sistemas-de-Gestion-empresarial/db"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Producto representa un artículo disponible en la tienda.
type Producto struct {
	ID            int
	Nombre        string
	Descripcion   string
	Precio        float64
	Stock         int
	SKU           string
	Activo        bool
	FechaCreacion time.Time
}

// ProductoCategoria representa la relación entre productos y categorías.
type ProductoCategoria struct {
	IDProducto  int
	IDCategoria int
}

// GetProductoByID devuelve un producto por su identificador o un error si no existe.
func GetProductoByID(id int) (Producto, error) {
	var producto Producto
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return producto, err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("SELECT id_producto, nombre, descripcion, precio, stock, sku, activo, fecha_creacion FROM productos WHERE id_producto = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return producto, err
	}
	defer stmt.Close()

	var descripcion, sku sql.NullString
	row := stmt.QueryRow(id)
	err = row.Scan(&producto.ID, &producto.Nombre, &descripcion, &producto.Precio, &producto.Stock, &sku, &producto.Activo, &producto.FechaCreacion)
	if err != nil {
		if err == sql.ErrNoRows {
			return producto, fmt.Errorf("producto no encontrado con ID: %d", id)
		}
		log.Println("Error al escanear la consulta sql", err)
		return producto, err
	}
	producto.Descripcion = descripcion.String
	producto.SKU = sku.String

	log.Println("Producto obtenido", producto)
	return producto, nil
}

// GetAllProductos devuelve la lista completa de productos en la base de datos.
func GetAllProductos() ([]Producto, error) {
	var productos []Producto
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return productos, err
	}
	defer DB.Close()

	rows, err := DB.Query("SELECT id_producto, nombre, descripcion, precio, stock, sku, activo, fecha_creacion FROM productos")
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return productos, err
	}
	defer rows.Close()

	for rows.Next() {
		var producto Producto
		var descripcion, sku sql.NullString
		err = rows.Scan(&producto.ID, &producto.Nombre, &descripcion, &producto.Precio, &producto.Stock, &sku, &producto.Activo, &producto.FechaCreacion)
		if err != nil {
			log.Println("Error al escanear la consulta sql", err)
			return productos, err
		}
		producto.Descripcion = descripcion.String
		producto.SKU = sku.String
		productos = append(productos, producto)
	}
	if err = rows.Err(); err != nil {
		log.Println("Error al obtener los productos", err)
		return productos, err
	}
	log.Println("Productos obtenidos", productos)
	return productos, nil
}

// CreateProducto inserta un nuevo producto en la base de datos.
func CreateProducto(nombre, descripcion string, precio float64, stock int, sku string, activo bool) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("INSERT INTO productos (nombre, descripcion, precio, stock, sku, activo) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(nombre, descripcion, precio, stock, sku, activo)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return err
	}
	log.Println("Producto creado exitosamente")
	return nil
}

// UpdateProducto actualiza la información de un producto existente.
func UpdateProducto(id int, nombre, descripcion string, precio float64, stock int, sku string, activo bool) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("UPDATE productos SET nombre = ?, descripcion = ?, precio = ?, stock = ?, sku = ?, activo = ? WHERE id_producto = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(nombre, descripcion, precio, stock, sku, activo, id)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return err
	}
	log.Println("Producto actualizado exitosamente")
	return nil
}

// DeleteProducto elimina un producto por su ID.
func DeleteProducto(id int) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("DELETE FROM productos WHERE id_producto = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return err
	}
	log.Println("Producto eliminado exitosamente")
	return nil
}

// AsignarCategoria asigna una categoría a un producto en la tabla intermedia.
func AsignarCategoria(idProducto, idCategoria int) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("INSERT INTO producto_categorias (id_producto, id_categoria) VALUES (?, ?)")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(idProducto, idCategoria)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return err
	}
	log.Println("Categoría asignada al producto exitosamente")
	return nil
}

package models

import (
	"Go-Sistemas-de-Gestion-empresarial/db"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Carrito struct {
	ID            int
	IDCliente     int
	FechaCreacion time.Time
}

type ItemCarrito struct {
	ID         int
	IDCarrito  int
	IDProducto int
	Cantidad   int
}

// GetCarritoByID obtiene un carrito por su ID
func GetCarritoByID(id int) (Carrito, error) {
	var carrito Carrito
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return carrito, err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("SELECT id_carrito, id_cliente, fecha_creacion FROM carritos WHERE id_carrito = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return carrito, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	err = row.Scan(&carrito.ID, &carrito.IDCliente, &carrito.FechaCreacion)
	if err != nil {
		if err == sql.ErrNoRows {
			return carrito, fmt.Errorf("carrito no encontrado con ID: %d", id)
		}
		log.Println("Error al escanear la consulta sql", err)
		return carrito, err
	}
	log.Println("Carrito obtenido", carrito)
	return carrito, nil
}

func GetCarritoByClienteID(id int) (Carrito, error) {
	var carrito Carrito
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return carrito, err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("SELECT id_carrito, id_cliente, fecha_creacion FROM carritos WHERE id_cliente = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return carrito, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	err = row.Scan(&carrito.ID, &carrito.IDCliente, &carrito.FechaCreacion)
	if err != nil {
		if err == sql.ErrNoRows {
			return carrito, fmt.Errorf("carrito no encontrado con ID: %d", id)
		}
		log.Println("Error al escanear la consulta sql", err)
		return carrito, err
	}
	log.Println("Carrito obtenido", carrito)
	return carrito, nil
}

func CreateCarrito(idCliente int) error {
	log.Println("Cliente ID", idCliente)

	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return err
	}
	defer DB.Close()
	carrito, err := GetCarritoByClienteID(idCliente)
	if err != nil {
		log.Println("Error al obtener el carrito", err)
	}
	log.Println("Carrito verificado para cliente", idCliente, "Encontrado ID:", carrito.ID)
	if carrito.ID != 0 {
		log.Println("Carrito ya existe")
		return nil
	}
	stmt, err := DB.Prepare("INSERT INTO carritos (id_cliente) VALUES (?)")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(idCliente)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return err
	}
	log.Println("Carrito creado exitosamente")
	return nil
}

func DeleteCarrito(id int) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("DELETE FROM carritos WHERE id_carrito = ?")
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
	log.Println("Carrito eliminado exitosamente")
	return nil
}

func AgregarItemCarrito(idCarrito, idProducto, cantidad int) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return err
	}
	defer DB.Close()
	log.Println("Intentando agregar item: CarritoID:", idCarrito, "ProductoID:", idProducto, "Cantidad:", cantidad)
	stmt, err := DB.Prepare("INSERT INTO items_carrito (id_carrito, id_producto, cantidad) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(idCarrito, idProducto, cantidad)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return err
	}
	log.Println("Item agregado al carrito exitosamente")
	return nil
}

func GetItemsByCarritoID(idCarrito int) ([]ItemCarrito, error) {
	var items []ItemCarrito
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return items, err
	}
	defer DB.Close()

	rows, err := DB.Query("SELECT id_item, id_carrito, id_producto, cantidad FROM items_carrito WHERE id_carrito = ?", idCarrito)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item ItemCarrito
		err = rows.Scan(&item.ID, &item.IDCarrito, &item.IDProducto, &item.Cantidad)
		if err != nil {
			log.Println("Error al escanear la consulta sql", err)
			return items, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		log.Println("Error al obtener los items del carrito", err)
		return items, err
	}
	log.Println("Items del carrito obtenidos", items)
	return items, nil
}

func UpdateItemCarrito(idItem, cantidad int) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("UPDATE items_carrito SET cantidad = ? WHERE id_item = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cantidad, idItem)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return err
	}
	log.Println("Item actualizado exitosamente")
	return nil
}

func EmptyCarrito(idCarrito int) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("DELETE FROM items_carrito WHERE id_carrito = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(idCarrito)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return err
	}
	log.Println("Carrito vaciado exitosamente")
	return nil
}

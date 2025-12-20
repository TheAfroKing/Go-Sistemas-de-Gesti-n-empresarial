package models

import (
	"Go-Sistemas-de-Gestion-empresarial/db"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Pedido struct {
	ID            int
	IDCliente     int
	Fecha         time.Time
	Estado        string
	Total         float64
	MetodoPago    string
	TransaccionID string
}

type DetallePedido struct {
	ID             int
	IDPedido       int
	IDProducto     int
	Cantidad       int
	PrecioUnitario float64
	// DetallePedido representa una línea de un pedido con cantidad y precio unitario.
	Subtotal float64
}

func GetPedidoByID(id int) (Pedido, error) {
	var pedido Pedido
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return pedido, err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("SELECT id_pedido, id_cliente, fecha, estado, total, metodo_pago, transaccion_id FROM pedidos WHERE id_pedido = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return pedido, err
	}
	defer stmt.Close()

	var metodoPago, transaccionID sql.NullString
	row := stmt.QueryRow(id)
	err = row.Scan(&pedido.ID, &pedido.IDCliente, &pedido.Fecha, &pedido.Estado, &pedido.Total, &metodoPago, &transaccionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return pedido, fmt.Errorf("pedido no encontrado con ID: %d", id)
		}
		log.Println("Error al escanear la consulta sql", err)
		return pedido, err
	}
	pedido.MetodoPago = metodoPago.String
	pedido.TransaccionID = transaccionID.String

	log.Println("Pedido obtenido", pedido)
	return pedido, nil
}

func GetAllPedidos() ([]Pedido, error) {
	var pedidos []Pedido
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return pedidos, err
	}
	defer DB.Close()

	rows, err := DB.Query("SELECT id_pedido, id_cliente, fecha, estado, total, metodo_pago, transaccion_id FROM pedidos")
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return pedidos, err
	}
	defer rows.Close()

	for rows.Next() {
		var pedido Pedido
		var metodoPago, transaccionID sql.NullString
		err = rows.Scan(&pedido.ID, &pedido.IDCliente, &pedido.Fecha, &pedido.Estado, &pedido.Total, &metodoPago, &transaccionID)
		if err != nil {
			log.Println("Error al escanear la consulta sql", err)
			return pedidos, err
		}
		pedido.MetodoPago = metodoPago.String
		pedido.TransaccionID = transaccionID.String
		pedidos = append(pedidos, pedido)
	}
	if err = rows.Err(); err != nil {
		log.Println("Error al obtener los pedidos", err)
		return pedidos, err
	}
	log.Println("Pedidos obtenidos", pedidos)
	return pedidos, nil
}

func CreatePedido(idCliente int, total float64, metodoPago, transaccionID string) (int, error) {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return 0, err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("INSERT INTO pedidos (id_cliente, total, metodo_pago, transaccion_id, estado) VALUES (?, ?, ?, ?, 'PENDIENTE')")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(idCliente, total, metodoPago, transaccionID)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error al obtener el ID del pedido insertado", err)
		return 0, err
	}

	log.Println("Pedido creado exitosamente con ID:", id)
	return int(id), nil
}

func CreateDetallePedido(idPedido, idProducto, cantidad int, precioUnitario float64) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("INSERT INTO detalles_pedido (id_pedido, id_producto, cantidad, precio_unitario) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(idPedido, idProducto, cantidad, precioUnitario)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return err
	}
	log.Println("Detalle de pedido agregado exitosamente")
	return nil
}

func GetDetallesByPedidoID(idPedido int) ([]DetallePedido, error) {
	var detalles []DetallePedido
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return detalles, err
	}
	defer DB.Close()

	rows, err := DB.Query("SELECT id_detalle, id_pedido, id_producto, cantidad, precio_unitario, subtotal FROM detalles_pedido WHERE id_pedido = ?", idPedido)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return detalles, err
	}
	defer rows.Close()

	for rows.Next() {
		var detalle DetallePedido
		err = rows.Scan(&detalle.ID, &detalle.IDPedido, &detalle.IDProducto, &detalle.Cantidad, &detalle.PrecioUnitario, &detalle.Subtotal)
		if err != nil {
			log.Println("Error al escanear la consulta sql", err)
			return detalles, err
		}
		detalles = append(detalles, detalle)
	}
	if err = rows.Err(); err != nil {
		log.Println("Error al obtener los detalles del pedido", err)
		return detalles, err
	}
	log.Println("Detalles del pedido obtenidos", detalles)
	return detalles, nil
}

func GetPedidosByClienteID(idCliente int) ([]Pedido, error) {
	var pedidos []Pedido
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return pedidos, err
	}
	defer DB.Close()

	rows, err := DB.Query("SELECT id_pedido, id_cliente, fecha, estado, total, metodo_pago, transaccion_id FROM pedidos WHERE id_cliente = ? ORDER BY fecha DESC", idCliente)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return pedidos, err
	}
	defer rows.Close()

	for rows.Next() {
		var pedido Pedido
		var metodoPago, transaccionID sql.NullString
		err = rows.Scan(&pedido.ID, &pedido.IDCliente, &pedido.Fecha, &pedido.Estado, &pedido.Total, &metodoPago, &transaccionID)
		if err != nil {
			log.Println("Error al escanear la consulta sql", err)
			return pedidos, err
		}
		pedido.MetodoPago = metodoPago.String
		pedido.TransaccionID = transaccionID.String
		pedidos = append(pedidos, pedido)
	}
	if err = rows.Err(); err != nil {
		log.Println("Error al obtener los pedidos", err)
		return pedidos, err
	}
	return pedidos, nil
}

func UpdatePedidoStatus(id int, estado string) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("UPDATE pedidos SET estado = ? WHERE id_pedido = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(estado, id)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return err
	}
	log.Println("Estado del pedido actualizado exitosamente")
	return nil
}

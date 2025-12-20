package models

import (
	"Go-Sistemas-de-Gestion-empresarial/db"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Cliente struct {
	ID                 int
	Nombre             string
	Email              string
	PasswordHash       string
	Direccion          string
	Telefono           string
	Perfil             string
	FechaRegistro      time.Time
	FechaActualizacion time.Time
}

func GetClienteByID(id int) (Cliente, error) {
	var cliente Cliente
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return cliente, err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("SELECT id_cliente, nombre, email, password_hash, direccion, telefono, perfil, fecha_registro, fecha_actualizacion FROM clientes WHERE id_cliente = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return cliente, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	var direccion, telefono, perfil sql.NullString

	err = row.Scan(&cliente.ID, &cliente.Nombre, &cliente.Email, &cliente.PasswordHash, &direccion, &telefono, &perfil, &cliente.FechaRegistro, &cliente.FechaActualizacion)
	if err != nil {
		if err == sql.ErrNoRows {
			return cliente, fmt.Errorf("cliente no encontrado con ID: %d", id)
		}
		log.Println("Error al escanear la consulta sql", err)
		return cliente, err
	}
	cliente.Direccion = direccion.String
	cliente.Telefono = telefono.String
	cliente.Perfil = perfil.String

	log.Println("Cliente obtenido", cliente)
	return cliente, nil
}

func GetClienteByEmail(email string) (Cliente, error) {
	var cliente Cliente
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return cliente, err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("SELECT id_cliente, nombre, email, password_hash, direccion, telefono, perfil, fecha_registro, fecha_actualizacion FROM clientes WHERE email = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return cliente, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(email)
	var direccion, telefono, perfil sql.NullString

	err = row.Scan(&cliente.ID, &cliente.Nombre, &cliente.Email, &cliente.PasswordHash, &direccion, &telefono, &perfil, &cliente.FechaRegistro, &cliente.FechaActualizacion)
	if err != nil {
		if err == sql.ErrNoRows {
			return cliente, fmt.Errorf("cliente no encontrado con email: %s", email)
		}
		log.Println("Error al escanear la consulta sql", err)
		return cliente, err
	}
	cliente.Direccion = direccion.String
	cliente.Telefono = telefono.String
	cliente.Perfil = perfil.String

	return cliente, nil
}

func GetAllClientes() ([]Cliente, error) {
	var clientes []Cliente
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return clientes, err
	}
	defer DB.Close()

	rows, err := DB.Query("SELECT id_cliente, nombre, email, password_hash, direccion, telefono, perfil, fecha_registro, fecha_actualizacion FROM clientes")
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return clientes, err
	}
	defer rows.Close()

	for rows.Next() {
		var cliente Cliente
		var direccion, telefono, perfil sql.NullString
		err = rows.Scan(&cliente.ID, &cliente.Nombre, &cliente.Email, &cliente.PasswordHash, &direccion, &telefono, &perfil, &cliente.FechaRegistro, &cliente.FechaActualizacion)
		if err != nil {
			log.Println("Error al escanear la consulta sql", err)
			return clientes, err
		}
		cliente.Direccion = direccion.String
		cliente.Telefono = telefono.String
		cliente.Perfil = perfil.String
		clientes = append(clientes, cliente)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error al obtener los clientes", err)
		return clientes, err
	}

	log.Println("Clientes obtenidos", clientes)
	return clientes, nil
}

func CreateCliente(nombre, email, passwordHash, direccion, telefono string) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("INSERT INTO clientes (nombre, email, password_hash, direccion, telefono, perfil) VALUES (?, ?, ?, ?, ?, 'cliente')")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(nombre, email, passwordHash, direccion, telefono)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return err
	}
	log.Println("Cliente creado exitosamente")
	return nil
}

func UpdateCliente(id int, nombre, email, direccion, telefono string) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("UPDATE clientes SET nombre = ?, email = ?, direccion = ?, telefono = ? WHERE id_cliente = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(nombre, email, direccion, telefono, id)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return err
	}
	log.Println("Cliente actualizado exitosamente")
	return nil
}

func DeleteCliente(id int) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return err
	}
	defer DB.Close()

	stmt, err := DB.Prepare("DELETE FROM clientes WHERE id_cliente = ?")
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
	log.Println("Cliente eliminado exitosamente")
	return nil
}

func Login(email, password string) (Cliente, error) {
	cliente, err := GetClienteByEmail(email)
	if err != nil {
		return Cliente{}, err
	}

	if cliente.PasswordHash != password {
		return Cliente{}, fmt.Errorf("contraseña incorrecta")
	}

	return cliente, nil
}

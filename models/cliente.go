package models

import (
	"Go-Sistemas-de-Gestion-empresarial/db"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Cliente representa a un usuario registrado en el sistema.
type Cliente struct {
	ID                 int       // Identificador único del cliente
	Nombre             string    // Nombre completo del cliente
	Email              string    // Correo electrónico único
	PasswordHash       string    // Hash de la contraseña (encapsulada)
	Direccion          string    // Dirección física del cliente
	Telefono           string    // Número de teléfono de contacto
	Perfil             string    // Perfil del usuario (ej. "cliente", "admin")
	FechaRegistro      time.Time // Fecha en que se registró el cliente
	FechaActualizacion time.Time // Fecha de la última actualización de datos
}

// VerifyPassword verifica si la contraseña proporcionada coincide con el hash almacenado.
// Esto encapsula la lógica de verificación de contraseñas.
func (c *Cliente) VerifyPassword(password string) bool {
	// En un escenario real, aquí usaríamos bcrypt.CompareHashAndPassword
	return c.PasswordHash == password
}

// GetClienteByID obtiene un cliente por su ID desde la base de datos.
func GetClienteByID(id int) (Cliente, error) {
	var cliente Cliente
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return cliente, fmt.Errorf("error de conexión: %w", err)
	}
	defer DB.Close()

	stmt, err := DB.Prepare("SELECT id_cliente, nombre, email, password_hash, direccion, telefono, perfil, fecha_registro, fecha_actualizacion FROM clientes WHERE id_cliente = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return cliente, fmt.Errorf("error preparando consulta: %w", err)
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
		return cliente, fmt.Errorf("error al leer datos: %w", err)
	}
	cliente.Direccion = direccion.String
	cliente.Telefono = telefono.String
	cliente.Perfil = perfil.String

	log.Println("Cliente obtenido", cliente)
	return cliente, nil
}

// GetClienteByEmail recupera un cliente usando su dirección de correo electrónico.
func GetClienteByEmail(email string) (Cliente, error) {
	var cliente Cliente
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return cliente, fmt.Errorf("error de conexión: %w", err)
	}
	defer DB.Close()

	stmt, err := DB.Prepare("SELECT id_cliente, nombre, email, password_hash, direccion, telefono, perfil, fecha_registro, fecha_actualizacion FROM clientes WHERE email = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return cliente, fmt.Errorf("error preparando consulta: %w", err)
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
		return cliente, fmt.Errorf("error al leer datos: %w", err)
	}
	cliente.Direccion = direccion.String
	cliente.Telefono = telefono.String
	cliente.Perfil = perfil.String

	return cliente, nil
}

// GetAllClientes devuelve una lista de todos los clientes registrados.
func GetAllClientes() ([]Cliente, error) {
	var clientes []Cliente
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return clientes, fmt.Errorf("error de conexión: %w", err)
	}
	defer DB.Close()

	rows, err := DB.Query("SELECT id_cliente, nombre, email, password_hash, direccion, telefono, perfil, fecha_registro, fecha_actualizacion FROM clientes")
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return clientes, fmt.Errorf("error ejecutando consulta: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cliente Cliente
		var direccion, telefono, perfil sql.NullString
		err = rows.Scan(&cliente.ID, &cliente.Nombre, &cliente.Email, &cliente.PasswordHash, &direccion, &telefono, &perfil, &cliente.FechaRegistro, &cliente.FechaActualizacion)
		if err != nil {
			log.Println("Error al escanear la consulta sql", err)
			return clientes, fmt.Errorf("error escaneando fila: %w", err)
		}
		cliente.Direccion = direccion.String
		cliente.Telefono = telefono.String
		cliente.Perfil = perfil.String
		clientes = append(clientes, cliente)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error al obtener los clientes", err)
		return clientes, fmt.Errorf("error iterando filas: %w", err)
	}

	log.Println("Clientes obtenidos", clientes)
	return clientes, nil
}

// CreateCliente registra un nuevo cliente en la base de datos.
func CreateCliente(nombre, email, passwordHash, direccion, telefono string) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return fmt.Errorf("error de conexión: %w", err)
	}
	defer DB.Close()

	stmt, err := DB.Prepare("INSERT INTO clientes (nombre, email, password_hash, direccion, telefono, perfil) VALUES (?, ?, ?, ?, ?, 'cliente')")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return fmt.Errorf("error preparando consulta: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(nombre, email, passwordHash, direccion, telefono)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return fmt.Errorf("error ejecutando inserción: %w", err)
	}
	log.Println("Cliente creado exitosamente")
	return nil
}

// UpdateCliente actualiza los datos de un cliente existente.
func UpdateCliente(id int, nombre, email, direccion, telefono string) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return fmt.Errorf("error de conexión: %w", err)
	}
	defer DB.Close()

	stmt, err := DB.Prepare("UPDATE clientes SET nombre = ?, email = ?, direccion = ?, telefono = ? WHERE id_cliente = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return fmt.Errorf("error preparando consulta: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(nombre, email, direccion, telefono, id)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return fmt.Errorf("error ejecutando actualización: %w", err)
	}
	log.Println("Cliente actualizado exitosamente")
	return nil
}

// DeleteCliente elimina un cliente de la base de datos por su ID.
func DeleteCliente(id int) error {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Error al conectar con la base de datos", err)
		return fmt.Errorf("error de conexión: %w", err)
	}
	defer DB.Close()

	stmt, err := DB.Prepare("DELETE FROM clientes WHERE id_cliente = ?")
	if err != nil {
		log.Println("Error al preparar la consulta sql", err)
		return fmt.Errorf("error preparando consulta: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Println("Error al ejecutar la consulta sql", err)
		return fmt.Errorf("error ejecutando eliminación: %w", err)
	}
	log.Println("Cliente eliminado exitosamente")
	return nil
}

// Login autentica a un usuario verificando su email y contraseña.
func Login(email, password string) (Cliente, error) {
	cliente, err := GetClienteByEmail(email)
	if err != nil {
		return Cliente{}, err
	}

	// Usamos el método encapsulado para verificar la contraseña
	if !cliente.VerifyPassword(password) {
		return Cliente{}, fmt.Errorf("contraseña incorrecta")
	}

	return cliente, nil
}

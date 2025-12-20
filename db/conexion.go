// este archivo va a guardar las configuraciones de la base de datos.
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Connect crea y devuelve una conexión a la base de datos MySQL usando
// variables de entorno (`DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_NAME`).
// Devuelve error si la conexión o el ping fallan.
func Connect() (*sql.DB, error) {
	// cargar directamente los datos desde el env
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	// loq ue hacer es configurar la de conexion con mysql
	fmt.Sprintf("root:@tcp(localhost:3306)/pooenlinea")
	dns := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// crear la conexion
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}

	// probar que conexion esa correcta
	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Conectado exitosamente con la base de datos")

	return db, nil

}

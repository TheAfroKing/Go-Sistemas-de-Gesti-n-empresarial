
# Go E-commerce (Sistemas de Gestión Empresarial)

Proyecto de ejemplo de un sistema e-commerce desarrollado en Go como parte
del trabajo práctico para la materia. Incluye rutas de cliente y administración,
manejo de carrito, pedidos, autenticación básica con cookies y persistencia en MySQL.

**Autor**: Fabián Paredes

## Características
- Listado y detalle de productos
- Carrito de compras con gestión de items
- Proceso de checkout (simulado)
- Panel de administración para productos, pedidos y clientes
- Persistencia en MySQL

## Requisitos
- Go 1.18+ instalado
- MySQL (o servidor compatible) accesible
- Variables de entorno configuradas (ver sección "Configuración")

## Instalación
1. Clona el repositorio:

```bash
git clone https://github.com/TheAfroKing/Go-Sistemas-de-Gestion-empresarial.git
cd Go-Sistemas-de-Gestion-empresarial
```

2. Descarga dependencias (mod tidy):

```bash
go mod tidy
```

## Configuración
Crea un archivo `.env` en la raíz del proyecto con las siguientes variables:

```env
DB_USER=root
DB_PASSWORD=tu_password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=nombre_basedatos
PORT=8080
```

Nota: El proyecto usa `github.com/joho/godotenv` para cargar variables de entorno
en desarrollo. En producción preferir variables de entorno del sistema.

## Estructura del proyecto
- `eCommerce.go` : punto de entrada y registro de rutas
- `db/` : conexión a la base de datos (`conexion.go`)
- `handlers/` : controladores HTTP para cliente y admin
- `models/` : lógica y acceso a datos (productos, clientes, carrito, pedidos)
- `templates/` : vistas HTML
- `static/` : archivos estáticos (CSS, JS, imágenes)

Vistas principales de `templates/`:
- `templates/base.html` - layout principal
- `templates/cliente/` - vistas cliente (carrito, checkout, productos, perfil)
- `templates/admin/` - vistas de administración (productos, ordenes, clientes)

## Ejecutar
Con el `.env` configurado, ejecuta:

```bash
go run ./ecommerce.go
```

O compila el binario:

```bash
go build -o ecommerce .
./ecommerce
```

El servidor por defecto escucha en el puerto definido por `PORT` (8080 por defecto).

## Base de datos
Se asume un esquema MySQL con tablas para `clientes`, `productos`, `carritos`,
`items_carrito`, `pedidos` y `detalles_pedido`. 

## Diagramas
- Diagrama de Clases

![Diagrana de Clases](https://github.com/TheAfroKing/Go-Sistemas-de-Gestion-empresarial/blob/master/Diagrama%20de%20clases%20-%20final.png)
- Diagrama de Arquitectura

![Diagrana de Clases](https://github.com/TheAfroKing/Go-Sistemas-de-Gestion-empresarial/blob/master/Diagrama%20de%20arquitectura.png)



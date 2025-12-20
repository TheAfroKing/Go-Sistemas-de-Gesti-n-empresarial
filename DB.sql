CREATE DATABASE `ecommerce_db` 

USE `ecommerce_db` 

CREATE TABLE `carritos` (
  `id_carrito` int NOT NULL AUTO_INCREMENT,
  `id_cliente` int NOT NULL,
  `fecha_creacion` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_carrito`),
  KEY `id_cliente` (`id_cliente`),
  CONSTRAINT `carritos_ibfk_1` FOREIGN KEY (`id_cliente`) REFERENCES `clientes` (`id_cliente`)
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `categorias` (
  `id_categoria` int NOT NULL AUTO_INCREMENT,
  `nombre` varchar(50) NOT NULL,
  `descripcion` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id_categoria`),
  UNIQUE KEY `nombre` (`nombre`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `clientes` (
  `id_cliente` int NOT NULL AUTO_INCREMENT,
  `nombre` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `direccion` text,
  `telefono` varchar(20) DEFAULT NULL,
  `fecha_registro` datetime DEFAULT CURRENT_TIMESTAMP,
  `fecha_actualizacion` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `perfil` varchar(20) DEFAULT 'cliente',
  PRIMARY KEY (`id_cliente`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `detalles_pedido` (
  `id_detalle` int NOT NULL AUTO_INCREMENT,
  `id_pedido` int NOT NULL,
  `id_producto` int NOT NULL,
  `cantidad` int NOT NULL,
  `precio_unitario` decimal(10,2) NOT NULL,
  `subtotal` decimal(10,2) GENERATED ALWAYS AS ((`cantidad` * `precio_unitario`)) STORED,
  PRIMARY KEY (`id_detalle`),
  KEY `id_pedido` (`id_pedido`),
  KEY `id_producto` (`id_producto`),
  CONSTRAINT `detalles_pedido_ibfk_1` FOREIGN KEY (`id_pedido`) REFERENCES `pedidos` (`id_pedido`) ON DELETE CASCADE,
  CONSTRAINT `detalles_pedido_ibfk_2` FOREIGN KEY (`id_producto`) REFERENCES `productos` (`id_producto`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `items_carrito` (
  `id_item` int NOT NULL AUTO_INCREMENT,
  `id_carrito` int NOT NULL,
  `id_producto` int NOT NULL,
  `cantidad` int NOT NULL DEFAULT '1',
  PRIMARY KEY (`id_item`),
  KEY `id_carrito` (`id_carrito`),
  KEY `id_producto` (`id_producto`),
  CONSTRAINT `items_carrito_ibfk_1` FOREIGN KEY (`id_carrito`) REFERENCES `carritos` (`id_carrito`) ON DELETE CASCADE,
  CONSTRAINT `items_carrito_ibfk_2` FOREIGN KEY (`id_producto`) REFERENCES `productos` (`id_producto`)
) ENGINE=InnoDB AUTO_INCREMENT=55 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `pedidos` (
  `id_pedido` int NOT NULL AUTO_INCREMENT,
  `id_cliente` int NOT NULL,
  `fecha` datetime DEFAULT CURRENT_TIMESTAMP,
  `estado` enum('PENDIENTE','PAGADO','ENVIADO','ENTREGADO','CANCELADO') DEFAULT 'PENDIENTE',
  `total` decimal(10,2) NOT NULL,
  `metodo_pago` varchar(50) DEFAULT NULL,
  `transaccion_id` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id_pedido`),
  KEY `id_cliente` (`id_cliente`),
  CONSTRAINT `pedidos_ibfk_1` FOREIGN KEY (`id_cliente`) REFERENCES `clientes` (`id_cliente`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `producto_categorias` (
  `id_producto` int NOT NULL,
  `id_categoria` int NOT NULL,
  PRIMARY KEY (`id_producto`,`id_categoria`),
  KEY `id_categoria` (`id_categoria`),
  CONSTRAINT `producto_categorias_ibfk_1` FOREIGN KEY (`id_producto`) REFERENCES `productos` (`id_producto`) ON DELETE CASCADE,
  CONSTRAINT `producto_categorias_ibfk_2` FOREIGN KEY (`id_categoria`) REFERENCES `categorias` (`id_categoria`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `productos` (
  `id_producto` int NOT NULL AUTO_INCREMENT,
  `nombre` varchar(150) NOT NULL,
  `descripcion` text,
  `precio` decimal(10,2) NOT NULL,
  `stock` int NOT NULL DEFAULT '0',
  `sku` varchar(50) DEFAULT NULL,
  `activo` tinyint(1) DEFAULT '1',
  `fecha_creacion` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_producto`),
  UNIQUE KEY `sku` (`sku`),
  CONSTRAINT `productos_chk_1` CHECK ((`precio` >= 0)),
  CONSTRAINT `productos_chk_2` CHECK ((`stock` >= 0))
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
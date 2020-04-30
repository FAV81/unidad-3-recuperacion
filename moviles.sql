CREATE DATABASE moviles;

USE moviles;

CREATE TABLE celulares(
	id INT AUTO_INCREMENT PRIMARY KEY,
    precio DECIMAL(16,2),
    descripcion VARCHAR(450),
    marca VARCHAR(200),
    modelo VARCHAR(200),
    lanzamiento DATE,
    creado TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    

);

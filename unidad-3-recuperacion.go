package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Variable conection assigned.
var DB *sql.DB
var err error 

// Celular Struct model
type Celular struct {
	Id           int    `json:"id"`
	Precio       string `json:"precio"`
	Descripcion  string `json:"descripcion"`
	Marca        string `json:"marca"`
	Modelo       string `json:"modelo"`
	Lanzamiento  string `json:"lanzamiento"`
	Creado       string `json:"creado"`
}

func main() {
	router := gin.Default()

    DB, err = sql.Open("mysql", "root@/moviles")   
	if err != nil {
		panic(err.Error())
	}
	defer DB.Close()

	err = DB.Ping()
	if err != nil {
		panic(err.Error())
	}

	router.GET("/celulares/mostrartodos", mostrarTodos)
	router.GET("/celulares/mostraruno/:id", mostrarUno)
	router.DELETE("/celulares/borrar/:id", borrar)
	router.POST("/celulares/agregar", agregar)
	router.PUT("/celulares/modificar/:id", modificar)

	router.Run(":8000") // listen and serve on 0.0.0.0:8069 (for windows "localhost:8069")
}

func mostrarTodos(c *gin.Context) {
	rows, err := DB.Query("SELECT * FROM celulares")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} 
	defer rows.Close()

	var celulares []Celular 
	for rows.Next() {
		var celular Celular 
		rows.Scan(&celular.Id, &celular.Precio, &celular.Descripcion, &celular.Marca, &celular.Modelo, &celular.Lanzamiento, &celular.Creado)
		celulares = append(celulares, celular) 
	}
	c.JSON(200, celulares)
}

func mostrarUno(c *gin.Context) {
	id := c.Param("id")
    var celular Celular

	err := DB.QueryRow("SELECT * FROM celulares WHERE id = ?", id).Scan(&celular.Id, &celular.Precio, &celular.Descripcion, &celular.Marca, &celular.Modelo, &celular.Lanzamiento, &celular.Creado)
	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"warning": "no existing data"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, celular)
}

func borrar(c *gin.Context) {
	id := c.Param("id")     
	stmt, err := DB.Query("DELETE FROM celulares WHERE id = ?", id)
	if err != nil {
		c.JSON(500, gin.H{"message": "Row not deleted"})
		return
	}else{
		c.JSON(200, gin.H{"message": "Deleted successuflly"})
	}	
	defer stmt.Close()
}


func agregar(c *gin.Context) {
	precio := c.PostForm("precio")
	descripcion := c.PostForm("descripcion")
	marca := c.PostForm("marca")
	modelo := c.PostForm("modelo")
	lanzamiento := c.PostForm("lanzamiento")

	stmt, err := DB.Prepare("INSERT INTO celulares (`precio`, `descripcion`, `marca`, `modelo`, `lanzamiento`) VALUES (?,?,?,?,?)")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	res, err := stmt.Exec(precio, descripcion, marca, modelo, lanzamiento)
	lid, err := res.LastInsertId()
	defer stmt.Close()

	c.JSON(201, gin.H{"status" : 201, "message" : "Product created successfully!", "Id": lid})
}



func modificar(c *gin.Context) {
	id := c.Param("id")
	precio := c.PostForm("precio")
	descripcion := c.PostForm("descripcion")
	marca := c.PostForm("marca")
	modelo := c.PostForm("modelo")
	lanzamiento := c.PostForm("lanzamiento")

	stmt, err := DB.Prepare("UPDATE celulares SET `precio` = ?, `descripcion` = ?, `marca` = ?, `modelo` = ?, `lanzamiento` = ? WHERE id = ?;")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	_, err = stmt.Exec(precio, descripcion, marca, modelo, lanzamiento, id)

	defer stmt.Close()

	c.JSON(200, gin.H{"message": "Successfully updated"})
}
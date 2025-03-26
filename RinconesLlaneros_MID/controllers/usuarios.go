package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/golang-jwt/jwt/v4"
	servicios "github.com/sena_2824182/RinconesLlaneros_MID/RinconesLlaneros_MID/Services"
	"golang.org/x/crypto/bcrypt"
)

// UsuariosController operations for Usuarios
type UsuariosController struct {
	beego.Controller
}


type Claims struct {
	Id int `json:"Id"`
	jwt.StandardClaims
}

// URLMapping ...
func (c *UsuariosController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Put", c.PutContraseña)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Usuarios
// @Param	body		body 	models.Usuarios	true		"body for Usuarios content"
// @Success 201 {object} models.Usuarios
// @Failure 403 body is empty
// @router / [post]
func (c *UsuariosController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Usuarios by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Usuarios
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UsuariosController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Usuarios
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Usuarios
// @Failure 403
// @router / [get]
func (c *UsuariosController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Usuarios
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Usuarios	true		"body for Usuarios content"
// @Success 200 {object} models.Usuarios
// @Failure 403 :id is not int
// @router /:id [put]
// Put reenvía la solicitud PUT a otro API (backend)
func (c *UsuariosController) Put() {
	// Obtén el ID del usuario desde la URL
	id := c.Ctx.Input.Param(":id")

	// Definir la URL del backend
	apiURL := "http://localhost:8082/v1/Usuarios/" + id // Cambia esta URL al endpoint correcto de tu backend

	// Estructura para recibir los datos actualizados del usuario
	var usuario map[string]interface{}

	// Deserializamos el cuerpo de la solicitud JSON en la variable 'usuario'
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &usuario); err != nil {
		fmt.Println("Error al parsear el JSON:", err)
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al parsear el JSON",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Crear una nueva solicitud PUT al backend
	client := &http.Client{Timeout: 10 * time.Second} // Establecer un tiempo de espera
	usuarioJSON, _ := json.Marshal(usuario)           // Convertimos los datos a JSON

	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(usuarioJSON))
	if err != nil {
		fmt.Println("Error al crear la solicitud:", err)
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al crear la solicitud",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Establecer encabezados de la solicitud
	req.Header.Set("Content-Type", "application/json")

	// Realizar la solicitud PUT al backend
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al realizar la solicitud al backend:", err)
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al realizar la solicitud al backend",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta del backend
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error al leer la respuesta del backend:", err)
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al leer la respuesta del backend",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Verificar el estado de la respuesta
	var apiResponse map[string]interface{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Println("Error al parsear la respuesta del backend:", err)
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al parsear la respuesta del backend",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Si la respuesta del backend es exitosa, retornamos un mensaje de éxito
	if resp.StatusCode == 200 {
		c.Data["json"] = map[string]interface{}{
			"Message":             "Usuario actualizado correctamente",
			"status":              200,
			"success":             true,
			"usuario actualizado": apiResponse, // Aquí se puede devolver la respuesta del backend
		}
	} else {
		// En caso de error en el backend
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al actualizar el usuario en el backend",
			"status":  500,
			"success": false,
		}
	}

	// Servir la respuesta JSON al frontend
	c.ServeJSON()
}

// PutContraseña ...
// @Title PutContraseña
// @Description Cambiar la contraseña de un usuario
// @Param	id		path 	string	true		"El ID del usuario"
// @Param	body		body 	models.CambioContraseña	true		"El cuerpo de la solicitud con la contraseña actual y la nueva contraseña"
// @Success 200 {object} map[string]interface{}
// @Failure 403 body is empty
// @router /cambiar-contrasena/:id [put]
func (c *UsuariosController) PutContraseña() {
	var datos struct {
		Token           string `json:"token"`
		NuevaContraseña string `json:"nueva_contraseña"`
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &datos); err != nil {
		c.Data["json"] = map[string]interface{}{"success": false, "message": "Datos inválidos"}
		c.ServeJSON()
		return
	}

	claims, err := servicios.ValidarTokenRecuperacion(datos.Token)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"success": false, "message": "Token inválido o expirado"}
		c.ServeJSON()
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(datos.NuevaContraseña), bcrypt.DefaultCost)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"success": false, "message": "Error al encriptar contraseña"}
		c.ServeJSON()
		return
	}

	err = servicios.ActualizarContraseñaUsuario(claims.Id, string(hashedPassword))
	if err != nil {
		c.Data["json"] = map[string]interface{}{"success": false, "message": "Error al actualizar contraseña"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"success": true, "message": "Contraseña actualizada correctamente"}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Usuarios
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UsuariosController) Delete() {

}

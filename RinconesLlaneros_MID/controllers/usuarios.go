package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/astaxie/beego"
)

// UsuariosController operations for Usuarios
type UsuariosController struct {
	beego.Controller
}

// URLMapping ...
func (c *UsuariosController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
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
func (c *UsuariosController) Put() {
	// Definir la URL del API CRUD (esto puede provenir de la configuración)
	id := c.Ctx.Input.Param(":id")                      // Obtenemos el ID del usuario desde la URL
	apiURL := "http://localhost:8082/v1/Usuarios/" + id // Asegúrate de que esta URL esté bien configurada

	// Crear la estructura que enviarás en el PUT, adaptándola a tu modelo de datos
	usuario := map[string]interface{}{
		"Nombre":         "Nuevo Nombre",
		"Correo":         "nuevo@correo.com",
		"Cedula":         "1234567890",
		"NumeroTelefono": "123456789",
	}

	// Convertir la estructura a JSON
	usuarioJSON, err := json.Marshal(usuario)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al crear el JSON de la solicitud",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Realizar la solicitud PUT al API CRUD
	client := &http.Client{Timeout: 10 * time.Second} // Tiempo de espera para la solicitud
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(usuarioJSON))
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al crear la solicitud",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Establecer el tipo de contenido como JSON
	req.Header.Set("Content-Type", "application/json")

	// Hacer la solicitud y obtener la respuesta
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al realizar la solicitud al API CRUD",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al leer la respuesta del API CRUD",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Estructura que esperamos recibir del API CRUD (ajustarla a tu formato real)
	var apiResponse map[string]interface{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al parsear la respuesta JSON",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Verificamos si la respuesta contiene un mensaje de éxito
	if resp.StatusCode == 200 {
		c.Data["json"] = map[string]interface{}{
			"Message": "Usuario actualizado exitosamente",
			"status":  200,
			"success": true,
			"usuario": apiResponse, // Aquí puedes devolver los detalles actualizados si es necesario
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al actualizar el usuario",
			"status":  500,
			"success": false,
		}
	}

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

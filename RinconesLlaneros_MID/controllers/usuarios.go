package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/sena_2824182/RinconesLlaneros_MID/RinconesLlaneros_MID/models"
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
// @Failure 400 Bad Request
// @Failure 500 Internal Server Error
// @router / [post]
func (c *UsuariosController) Post() {
	var usuario models.Usuarios

	// Decodificar el JSON recibido
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &usuario); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  400,
			"message": "Error en el formato de entrada: " + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Validar que los campos esenciales no estén vacíos
	if usuario.Nombre == "" || usuario.Correo == "" || usuario.Cedula == "" || usuario.Rol == nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  400,
			"message": "Faltan campos obligatorios: Nombre, Correo, Cedula y Rol.",
		}
		c.ServeJSON()
		return
	}

	// Convertir a JSON para enviar al CRUD local
	jsonData, err := json.Marshal(usuario)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  500,
			"message": "Error al codificar usuario: " + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Hacer la solicitud HTTP POST al CRUD local
	resp, err := http.Post("http://localhost:8081/v1/Usuarios", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  500,
			"message": "Error al comunicarse con el CRUD local: " + err.Error(),
		}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta del CRUD y devolverla tal cual
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  500,
			"message": "Error al leer la respuesta del CRUD: " + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Enviar la respuesta obtenida del CRUD local
	c.Ctx.Output.SetStatus(resp.StatusCode)
	c.Data["json"] = json.RawMessage(body)
	c.ServeJSON()
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

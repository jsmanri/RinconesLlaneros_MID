package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/sena_2824182/RinconesLlaneros_MID/RinconesLlaneros_MID/services"
)

// AdminController operations for Admin
type AdminController struct {
	beego.Controller
}

// URLMapping ...
func (c *AdminController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Admin
// @Param	body		body 	models.Admin	true		"body for Admin content"
// @Success 201 {object} models.Admin
// @Failure 403 body is empty
// @router / [post]
func (c *AdminController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Admin by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Admin
// @Failure 403 :id is empty
// @router /:id [get]
func (c *AdminController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Admin
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Admin
// @Failure 403
// @router / [get]
func (c *AdminController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Admin
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Admin	true		"body for Admin content"
// @Success 200 {object} models.Admin
// @Failure 403 :id is not int
// @router /:id [put]
func (c *AdminController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Admin
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *AdminController) Delete() {
	idUsuario, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Data["json"] = map[string]string{"error": "ID inválido"}
		c.ServeJSON()
		return
	}

	// Obtener los sitios turísticos del usuario antes de eliminarlo
	sitios, err := services.ObtenerSitiosPorUsuario(idUsuario)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Error al consultar sitios turísticos"}
		c.ServeJSON()
		return
	}

	if len(sitios) > 0 {
		// Enviar respuesta con los sitios antes de proceder
		c.Data["json"] = map[string]interface{}{
			"mensaje":      "El usuario tiene sitios turísticos registrados",
			"sitios":       sitios,
			"confirmacion": "Envíe 'confirmar' en el cuerpo para proceder con la eliminación",
		}
		c.ServeJSON()
		return
	}

	// Leer el cuerpo de la solicitud para obtener la confirmación
	var requestBody map[string]string
	json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody)

	if requestBody["confirmacion"] != "confirmar" {
		c.Data["json"] = map[string]string{"mensaje": "Eliminación cancelada"}
		c.ServeJSON()
		return
	}

	// Realizar el borrado lógico
	if err := services.DesactivarUsuario(idUsuario); err != nil {
		c.Data["json"] = map[string]string{"error": "No se pudo desactivar el usuario"}
	} else {
		c.Data["json"] = map[string]string{"mensaje": "Usuario y sitios desactivados correctamente"}
	}
	c.ServeJSON()
}
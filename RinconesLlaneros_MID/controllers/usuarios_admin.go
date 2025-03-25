package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/sena_2824182/RinconesLlaneros_MID/RinconesLlaneros_MID/services"
)

// Usuarios_adminController operations for Usuarios_admin
type Usuarios_adminController struct {
	beego.Controller
}

// URLMapping ...
func (c *Usuarios_adminController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Usuarios_admin
// @Param	body		body 	models.Usuarios_admin	true		"body for Usuarios_admin content"
// @Success 201 {object} models.Usuarios_admin
// @Failure 403 body is empty
// @router / [post]
func (c *Usuarios_adminController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Usuarios_admin by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Usuarios_admin
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Usuarios_adminController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Usuarios_admin
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Usuarios_admin
// @Failure 403
// @router / [get]
func (c *Usuarios_adminController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Usuarios_admin
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Usuarios_admin	true		"body for Usuarios_admin content"
// @Success 200 {object} models.Usuarios_admin
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Usuarios_adminController) Put() {

}

// Delete ...
// @Title DeleteUsuario
// @Description Elimina un usuario de forma lógica si es administrador
// @Param   id   path   int  true  "ID del usuario a eliminar"
// @Param   confirmar  query  string  false  "Confirmar eliminación"
// @Success 200 {object} map[string]interface{}
// @Failure 400 ID inválido
// @Failure 404 Usuario no encontrado
// @Failure 500 Error interno
// @router /usuarios_admin/:id [delete]
func (c *Usuarios_adminController) Delete() {
	idUsuario, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "ID de usuario inválido",
			"error":   err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Llamamos al servicio para desactivar el usuario y sus sitios turísticos
	err = services.DesactivarUsuario(idUsuario)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error al desactivar el usuario",
			"error":   err.Error(),
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("Usuario con ID %d y sus sitios turísticos han sido desactivados correctamente", idUsuario),
		}
	}

	c.ServeJSON()
}

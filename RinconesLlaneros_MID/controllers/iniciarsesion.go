package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
)

// IniciarSesionController operations for IniciarSesion
type IniciarSesionController struct {
	beego.Controller
}

// URLMapping ...
func (c *IniciarSesionController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create IniciarSesion
// @Param	body		body 	models.IniciarSesion	true		"body for IniciarSesion content"
// @Success 201 {object} models.IniciarSesion
// @Failure 403 body is empty
// @router / [post]
func (c *IniciarSesionController) Post() {
	// 1. Struct para recibir login
	var loginReq struct {
		Correo     string `json:"correo"`
		Contrasena string `json:"contrasena"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &loginReq); err != nil {
		c.Data["json"] = map[string]string{"error": "json_invalido"}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}

	// 2. Consultar usuario en el CRUD
	usuariosUrl := beego.AppConfig.String("Servicio_Usuarios")
	reqUrl := fmt.Sprintf("%s?query=Correo:%s", usuariosUrl, loginReq.Correo)
	resp, err := http.Get(reqUrl)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "error_conexion_crud"}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// 3. Leer la respuesta del CRUD incluyendo el campo `Activo`
	var crudResp struct {
		UsuariosConsultados []struct {
			Id     int  `json:"Id"`
			Activo bool `json:"Activo"` // 🔹 Agregamos `Activo`
			Rol    struct {
				Id int `json:"Id"`
			} `json:"Rol"`
			IdCredencialesCredenciales struct {
				Id         int    `json:"Id"`
				Contrasena string `json:"Contraseña"`
			} `json:"IdCredencialesCredenciales"`
		} `json:"usuarios consultados"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&crudResp); err != nil || len(crudResp.UsuariosConsultados) == 0 {
		c.Data["json"] = map[string]string{"error": "correo_no_encontrado"}
		c.Ctx.Output.SetStatus(404)
		c.ServeJSON()
		return
	}

	usuario := crudResp.UsuariosConsultados[0]

	// 4. Validar contraseña
	if usuario.IdCredencialesCredenciales.Contrasena != loginReq.Contrasena {
		c.Data["json"] = map[string]string{"error": "contrasena_incorrecta"}
		c.Ctx.Output.SetStatus(401)
		c.ServeJSON()
		return
	}

	// 5. Retornar la respuesta incluyendo `activo`
	respuesta := map[string]interface{}{
		"id_usuario": usuario.Id,
		"id_rol":     usuario.Rol.Id,
		"activo":     usuario.Activo, // 🔹 Enviar estado de usuario activo/inactivo
	}
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get IniciarSesion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.IniciarSesion
// @Failure 403 :id is empty
// @router /:id [get]
func (c *IniciarSesionController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get IniciarSesion
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.IniciarSesion
// @Failure 403
// @router / [get]
func (c *IniciarSesionController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the IniciarSesion
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.IniciarSesion	true		"body for IniciarSesion content"
// @Success 200 {object} models.IniciarSesion
// @Failure 403 :id is not int
// @router /:id [put]
func (c *IniciarSesionController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the IniciarSesion
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *IniciarSesionController) Delete() {

}

package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/sena_2824182/RinconesLlaneros_MID/RinconesLlaneros_MID/Services" // Ajusta el import según tu estructura
)

type TokenController struct {
	beego.Controller
}

// POST /recuperar-password
func (c *TokenController) GenerarToken() {
	var datos struct {
		Id int `json:"Id"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &datos); err != nil {
		c.Data["json"] = map[string]interface{}{"success": false, "message": "Datos inválidos"}
		c.ServeJSON()
		return
	}

	token, err := servicios.GenerarTokenRecuperacion(datos.Id)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"success": false, "message": "Error generando token"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"success": true, "token": token}
	c.ServeJSON()
}

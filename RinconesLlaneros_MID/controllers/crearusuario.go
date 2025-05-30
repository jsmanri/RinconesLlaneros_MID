package controllers

import (
	"github.com/astaxie/beego"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// CrearusuarioController operations for Crearusuario
type CrearusuarioController struct {
	beego.Controller
}

// URLMapping ...
func (c *CrearusuarioController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Crearusuario
// @Param	body		body 	models.Crearusuario	true		"body for Crearusuario content"
// @Success 201 {object} models.Crearusuario
// @Failure 403 body is empty
// @router / [post]
func (c *CrearusuarioController) Post() {

		// Parsear el cuerpo de la solicitud
	var input map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &input); err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al parsear el cuerpo de la solicitud",
			"status":  400,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Validar que el campo "contrasena" exista
	contrasena, exists := input["contrasena"].(string)
	if !exists {
		c.Data["json"] = map[string]interface{}{
			"Message": "Falta el campo 'contrasena' en el cuerpo de la solicitud",
			"status":  400,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Crear el JSON para enviar al servicio de Credenciales
	jsonContraseña := map[string]string{
		"Contraseña": contrasena,
	}

	// Obtener la URL del servicio de credenciales desde app.conf
	servicioCredenciales := beego.AppConfig.String("Servicio_Credenciales")

	// Realizar la solicitud POST al servicio de credenciales
	client := &http.Client{Timeout: 10 * time.Second}
	reqBody := encodeToJSON(jsonContraseña)
	resp, err := client.Post(servicioCredenciales, "application/json", reqBody)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al realizar la solicitud al servicio de Credenciales",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Leer y parsear la respuesta del servicio de Credenciales
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al leer la respuesta del servicio de Credenciales",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	var credencialesResponse map[string]interface{}
	if err := json.Unmarshal(body, &credencialesResponse); err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al parsear la respuesta del servicio de Credenciales",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Validar que la creación de la credencial fue exitosa
	if !credencialesResponse["success"].(bool) {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al crear la credencial",
			"status":  400,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Obtener el ID de la credencial creada
	credencialCreada := credencialesResponse["credencial creada"].(map[string]interface{})
	idCredencial := credencialCreada["Id"].(float64)

	// Crear el campo FotoPerfil como texto que contiene un JSON
	fotoPerfil := map[string]string{
		"imagen_base64": input["fotoPerfil"].(string), // Asegúrate de que input["fotoPerfil"] sea la cadena base64 enviada desde Angular
	}
	fotoPerfilJSON, err := json.Marshal(fotoPerfil)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al procesar el campo 'FotoPerfil'",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Crear el JSON para el servicio de Usuarios
	jsonUsuario := map[string]interface{}{
		"Nombre":           input["nombre"],
		"Rol":              map[string]interface{}{"Id": input["rol"]},
		"Correo":           input["correo"],
		"Cedula":           input["cedula"],
		"NumeroTelefono":   input["telefono"],
		"FotoPerfil":       string(fotoPerfilJSON), // Aquí enviamos el JSON como texto
		"IdCredencialesCredenciales": map[string]interface{}{
			"Id": idCredencial,
		},
	}

	// Obtener la URL del servicio de usuarios desde app.conf
	servicioUsuarios := beego.AppConfig.String("Servicio_Usuarios")

	// Realizar la solicitud POST al servicio de usuarios
	reqBody = encodeToJSON(jsonUsuario)
	respUsuario, err := client.Post(servicioUsuarios, "application/json", reqBody)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al realizar la solicitud al servicio de Usuarios",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}
	defer respUsuario.Body.Close()

	// Leer y validar la respuesta del servicio de Usuarios
	if respUsuario.StatusCode != http.StatusCreated {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al crear el usuario en el servicio de Usuarios",
			"status":  respUsuario.StatusCode,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Devolver respuesta exitosa
	c.Data["json"] = map[string]interface{}{
		"Message": "Usuario creado exitosamente",
		"status":  201,
		"success": true,
	}
	c.ServeJSON()
}

// Helper function to encode a Go map into a JSON payload
func encodeToJSON(data interface{}) *bytes.Buffer {
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(data)
	return buffer
}



// GetOne ...
// @Title GetOne
// @Description get Crearusuario by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Crearusuario
// @Failure 403 :id is empty
// @router /:id [get]
func (c *CrearusuarioController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Crearusuario
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Crearusuario
// @Failure 403
// @router / [get]
func (c *CrearusuarioController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Crearusuario
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Crearusuario	true		"body for Crearusuario content"
// @Success 200 {object} models.Crearusuario
// @Failure 403 :id is not int
// @router /:id [put]
func (c *CrearusuarioController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Crearusuario
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *CrearusuarioController) Delete() {

}

package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"fmt"
	"bytes"

)

// ActulizarContraseñaController operations for ActulizarContraseña
type ActulizarContraseñaController struct {
	beego.Controller
}

// URLMapping ...
func (c *ActulizarContraseñaController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create ActulizarContraseña
// @Param	body		body 	models.ActulizarContraseña	true		"body for ActulizarContraseña content"
// @Success 201 {object} models.ActulizarContraseña
// @Failure 403 body is empty
// @router / [post]
func (c *ActulizarContraseñaController) Post() {

	var requestBody struct {
		IdUsuario       int    `json:"idUsuario"`
		IdCredenciales  int    `json:"idCredenciales"`
		NuevaContrasena string `json:"nuevaContrasena"`
	}

	// Parse JSON body from Angular
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody); err != nil {
		fmt.Println("Error parsing request body:", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error parsing request body",
			"error":   err.Error(),
		}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}
	fmt.Println("Request recibido:", requestBody)

	client := &http.Client{}

	// 1. Desactivar la credencial antigua
	updateData := map[string]interface{}{
		"Activo": false,
	}
	credencialesURL := beego.AppConfig.String("Servicio_Credenciales")
	putURL := fmt.Sprintf("%s/%d", credencialesURL, requestBody.IdCredenciales)
	updateBody, err := json.Marshal(updateData)
	if err != nil {
		fmt.Println("Error marshaling update data:", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error marshaling update data",
			"error":   err.Error(),
		}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}
	fmt.Println("PUT URL:", putURL)
	fmt.Println("Payload para desactivar credencial:", string(updateBody))

	req, err := http.NewRequest("PUT", putURL, bytes.NewBuffer(updateBody))
	if err != nil {
		fmt.Println("Error creando PUT request:", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error creando request a Credenciales",
			"error":   err.Error(),
		}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Println("Error en respuesta al desactivar credencial:", err, "Status:", resp.StatusCode)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error desactivando credencial anterior",
			"error":   err,
			"status":  resp.StatusCode,
		}
		c.Ctx.Output.SetStatus(502)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// 2. Crear nueva credencial
	createData := map[string]interface{}{
		"Contraseña": requestBody.NuevaContrasena,
	}
	createBody, err := json.Marshal(createData)
	if err != nil {
		fmt.Println("Error marshaling new credential:", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error marshaling new credential data",
			"error":   err.Error(),
		}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}
	fmt.Println("POST URL:", credencialesURL)
	fmt.Println("Payload para nueva credencial:", string(createBody))

	reqCreate, err := http.NewRequest("POST", credencialesURL, bytes.NewBuffer(createBody))
	if err != nil {
		fmt.Println("Error creando POST request:", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error creando request para nueva credencial",
			"error":   err.Error(),
		}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}
	reqCreate.Header.Set("Content-Type", "application/json")
	respCreate, err := client.Do(reqCreate)
	if err != nil || respCreate.StatusCode < 200 || respCreate.StatusCode > 299 {
		fmt.Println("Error en respuesta al crear nueva credencial:", err, "Status:", respCreate.StatusCode)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error creando nueva credencial",
			"error":   err,
			"status":  respCreate.StatusCode,
		}
		c.Ctx.Output.SetStatus(502)
		c.ServeJSON()
		return
	}
	defer respCreate.Body.Close()

	respCreateBody, err := ioutil.ReadAll(respCreate.Body)
	if err != nil {
		fmt.Println("Error leyendo respuesta de crear credencial:", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error leyendo la respuesta del servicio Credenciales",
			"error":   err.Error(),
		}
		c.Ctx.Output.SetStatus(502)
		c.ServeJSON()
		return
	}
	fmt.Println("Respuesta crear credencial:", string(respCreateBody))

	var credencialNuevaResp map[string]interface{}
	if err := json.Unmarshal(respCreateBody, &credencialNuevaResp); err != nil {
		fmt.Println("Error parseando respuesta nueva credencial:", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error parseando la respuesta del servicio Credenciales",
			"error":   err.Error(),
		}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}

	var idNuevaCredencial interface{}
	if cred, ok := credencialNuevaResp["credencial creada"].(map[string]interface{}); ok {
		idNuevaCredencial = cred["Id"]
		fmt.Println("Nueva credencial ID:", idNuevaCredencial)
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "No se pudo obtener el ID de la nueva credencial",
		}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}

	// 3. Consultar el usuario
	servicioUsuariosURL := beego.AppConfig.String("Servicio_Usuarios")
	getUsuarioURL := fmt.Sprintf("%s/%d", servicioUsuariosURL, requestBody.IdUsuario)
	fmt.Println("Consultando usuario en:", getUsuarioURL)
	respUsuario, err := http.Get(getUsuarioURL)
	if err != nil || respUsuario.StatusCode < 200 || respUsuario.StatusCode > 299 {
		fmt.Println("Error consultando usuario:", err, "Status:", respUsuario.StatusCode)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error consultando usuario",
			"error":   err,
			"status":  respUsuario.StatusCode,
		}
		c.Ctx.Output.SetStatus(502)
		c.ServeJSON()
		return
	}
	defer respUsuario.Body.Close()
	respUsuarioBody, err := ioutil.ReadAll(respUsuario.Body)
	if err != nil {
		fmt.Println("Error leyendo respuesta de usuario:", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error leyendo la respuesta del servicio de usuarios",
			"error":   err.Error(),
		}
		c.Ctx.Output.SetStatus(502)
		c.ServeJSON()
		return
	}
	fmt.Println("Respuesta usuario:", string(respUsuarioBody))

	// Aquí parseamos la respuesta con la estructura que retorna el CRUD
	var usuarioApiResponse map[string]interface{}
	if err := json.Unmarshal(respUsuarioBody, &usuarioApiResponse); err != nil {
		fmt.Println("Error parseando usuario:", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error parseando el usuario (response api)",
			"error":   err.Error(),
		}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}

	// Extraemos el usuario real bajo la clave "usuario consultados"
	usuarioData, ok := usuarioApiResponse["usuario consultados"].(map[string]interface{})
	if !ok {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "No se encontró la clave 'usuario consultados' en la respuesta del API de usuarios",
		}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}

	// 4. Modificar credencial del usuario (solo el campo Id)
	usuarioData["IdCredencialesCredenciales"] = map[string]interface{}{
		"Id": idNuevaCredencial,
	}
	fmt.Println("Usuario actualizado con nueva credencial:", usuarioData)

	// 5. Actualizar usuario (enviar SOLO el objeto usuario plano)
	updateUsuarioBody, err := json.Marshal(usuarioData)
	if err != nil {
		fmt.Println("Error serializando usuario actualizado:", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error serializando usuario actualizado",
			"error":   err.Error(),
		}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}
	putUsuarioURL := getUsuarioURL
	reqUsuarioUpdate, err := http.NewRequest("PUT", putUsuarioURL, bytes.NewBuffer(updateUsuarioBody))
	if err != nil {
		fmt.Println("Error creando PUT para actualizar usuario:", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error creando petición PUT al servicio de usuarios",
			"error":   err.Error(),
		}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}
	reqUsuarioUpdate.Header.Set("Content-Type", "application/json")
	respUsuarioUpdate, err := client.Do(reqUsuarioUpdate)
	if err != nil || respUsuarioUpdate.StatusCode < 200 || respUsuarioUpdate.StatusCode > 299 {
		fmt.Println("Error actualizando usuario:", err, "Status:", respUsuarioUpdate.StatusCode)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Error actualizando usuario con la nueva credencial",
			"error":   err,
			"status":  respUsuarioUpdate.StatusCode,
		}
		c.Ctx.Output.SetStatus(502)
		c.ServeJSON()
		return
	}
	defer respUsuarioUpdate.Body.Close()

	// Éxito total
	fmt.Println("Contraseña actualizada correctamente")
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "Contraseña actualizada correctamente",
	}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get ActulizarContraseña by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ActulizarContraseña
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ActulizarContraseñaController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get ActulizarContraseña
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.ActulizarContraseña
// @Failure 403
// @router / [get]
func (c *ActulizarContraseñaController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the ActulizarContraseña
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.ActulizarContraseña	true		"body for ActulizarContraseña content"
// @Success 200 {object} models.ActulizarContraseña
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ActulizarContraseñaController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the ActulizarContraseña
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ActulizarContraseñaController) Delete() {

}

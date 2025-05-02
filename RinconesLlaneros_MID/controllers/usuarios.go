package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/sena_2824182/RinconesLlaneros_MID/RinconesLlaneros_MID/services"
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
// @Param    body        body     models.Usuarios true        "body for Usuarios content"
// @Success 201 {object} models.Usuarios
// @Failure 400 Bad Request
// @Failure 500 Internal Server Error
// @router / [post]
// API_CRUD_Usuarios es la URL base de la API CRUD de usuarios.
// API_CRUD_Usuarios es la URL base de la API CRUD.
var API_CRUD_Usuarios string

func init() {
	API_CRUD_Usuarios = beego.AppConfig.String("API_CRUD_Usuarios")
	if API_CRUD_Usuarios == "" {
		log.Fatal("API_CRUD_Usuarios no configurado")
	}
	log.Printf("API_CRUD_Usuarios: %s", API_CRUD_Usuarios)
}

// Post maneja la creación de un nuevo usuario.
func (c *UsuariosController) Post() {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &jsonData); err != nil {
		c.respond(400, "JSON inválido: "+err.Error())
		return
	}

	// Validación de los campos requeridos
	if err := ValidateRequiredFields(jsonData, []string{"Nombre", "Correo", "Cedula", "Rol", "IdCredencialesCredenciales"}); err != nil {
		c.respond(400, "Faltan campos: "+err.Error())
		return
	}

	// Validar que "Rol" y "IdCredencialesCredenciales" contienen un ID válido
	rol, ok := jsonData["Rol"].(map[string]interface{})
	if !ok || rol["Id"] == nil {
		c.respond(400, "Rol o Id del rol no válido")
		return
	}
	rolID := int(rol["Id"].(float64))

	credenciales, ok := jsonData["IdCredencialesCredenciales"].(map[string]interface{})
	if !ok || credenciales["Id"] == nil {
		c.respond(400, "Credenciales o Id de credenciales no válido")
		return
	}
	credencialesID := int(credenciales["Id"].(float64))

	// Verificar los demás campos
	nombre := jsonData["Nombre"].(string)
	correo := jsonData["Correo"].(string)
	cedula := jsonData["Cedula"].(string)
	numeroTelefono := jsonData["NumeroTelefono"].(string)
	fotoPerfil := jsonData["FotoPerfil"].(string)

	// Crear el nuevo usuario
	nuevoUsuario := map[string]interface{}{
		"IdCredencialesCredenciales": map[string]interface{}{"Id": credencialesID},
		"Rol":                        map[string]interface{}{"Id": rolID},
		"Nombre":                     nombre,
		"Correo":                     correo,
		"Cedula":                     cedula,
		"NumeroTelefono":             numeroTelefono,
		"FotoPerfil":                 fotoPerfil,
		"Activo":                     true,
	}

	// Llamar a la función para agregar el usuario a CRUD
	usuarioID, err := c.addToCRUD("/usuarios", nuevoUsuario)
	if err != nil {
		c.respond(500, "Error creando usuario: "+err.Error())
		return
	}

	c.respondWithID(201, "Usuario creado", usuarioID)
}

// createCredenciales maneja la creación de credenciales.
func (c *UsuariosController) createCredenciales(data map[string]interface{}) (int, error) {
	password := data["Contraseña"].(string)
	correo := data["Correo"].(string)

	hashedPassword, err := services.HashContraseña(password)
	if err != nil {
		return 0, fmt.Errorf("error hash contraseña: %v", err)
	}

	credenciales := map[string]interface{}{
		"Contraseña": hashedPassword,
		"Activo":     true,
	}

	id, err := c.addToCRUD("/credenciales", credenciales)
	if err != nil {
		return 0, err
	}

	// Enviar correo de verificación en segundo plano
	go func() {
		token := "verificacion_" + strconv.Itoa(id)
		if err := services.EnviarCorreo(correo, token); err != nil {
			log.Printf("Error enviando correo de verificación: %v", err)
		}
	}()

	return id, nil
}

// addToCRUD envía un objeto al CRUD y devuelve el ID creado.
func (c *UsuariosController) addToCRUD(endpoint string, data map[string]interface{}) (int, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("error serializando JSON: %v", err)
	}

	resp, err := http.Post(API_CRUD_Usuarios+endpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return 0, fmt.Errorf("error enviando petición: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error leyendo respuesta: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("error API (%d): %s", resp.StatusCode, string(body))
	}

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		return 0, fmt.Errorf("error deserializando respuesta: %v", err)
	}

	idFloat, ok := res["Id"].(float64)
	if !ok {
		return 0, fmt.Errorf("Id inválido en respuesta")
	}

	return int(idFloat), nil
}

// GetRolIDFromCRUD obtiene el ID de un rol desde el CRUD.
func GetRolIDFromCRUD(rolNombre string) (int, error) {
	url := fmt.Sprintf("%s/roles?query=Nombre:%s", API_CRUD_Usuarios, rolNombre)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error conexión API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error leyendo respuesta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error API (%d): %s", resp.StatusCode, string(body))
	}

	var roles []map[string]interface{}
	if err := json.Unmarshal(body, &roles); err != nil {
		return 0, fmt.Errorf("error deserializando roles: %v", err)
	}

	if len(roles) == 0 {
		return 0, fmt.Errorf("rol no encontrado: %s", rolNombre)
	}

	idFloat, ok := roles[0]["Id"].(float64)
	if !ok {
		return 0, fmt.Errorf("Id de rol inválido")
	}

	return int(idFloat), nil
}

// ValidateRequiredFields verifica que existan campos requeridos.
func ValidateRequiredFields(data map[string]interface{}, fields []string) error {
	for _, field := range fields {
		if _, ok := data[field]; !ok {
			return fmt.Errorf("campo faltante: %s", field)
		}
	}
	return nil
}

// respond responde con un mensaje simple.
func (c *UsuariosController) respond(status int, message string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]interface{}{"success": status < 400, "status": status, "message": message}
	c.ServeJSON()
}

// respondWithID responde con un mensaje y un ID.
func (c *UsuariosController) respondWithID(status int, message string, id int) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]interface{}{"success": true, "status": status, "message": message, "id": id}
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
func (c *UsuariosController) DeleteUsuario() {
}

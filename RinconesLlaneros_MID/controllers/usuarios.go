package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (c *UsuariosController) Post() {
	bodyBytes := c.Ctx.Input.RequestBody

	if len(bodyBytes) == 0 {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"success": false, "status": 400, "message": "Cuerpo de la solicitud vacío"}
		c.ServeJSON()
		return
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"success": false, "status": 400, "message": "JSON inválido: " + err.Error()}
		c.ServeJSON()
		return
	}

	if err := ValidateRequiredFields(jsonData, []string{"Nombre", "Correo", "Cedula", "Rol", "Contraseña"}); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"success": false, "status": 400, "message": "Faltan campos: " + err.Error()}
		c.ServeJSON()
		return
	}

	credencialesID, err := c.PostCredenciales(jsonData)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"success": false, "status": 500, "message": "Error credenciales: " + err.Error()}
		c.ServeJSON()
		return
	}

	var rolID int
	if rolNombre, ok := jsonData["Rol"].(string); ok {
		rolID, err = GetRolIDFromCRUD(rolNombre)
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]interface{}{"success": false, "status": 500, "message": "Error Rol: " + err.Error()}
			c.ServeJSON()
			return
		}
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"success": false, "status": 400, "message": "Rol inválido"}
		c.ServeJSON()
		return
	}

	nuevoUsuario := map[string]interface{}{
		"IdCredencialesCredenciales": map[string]interface{}{"Id": credencialesID},
		"Rol":                         map[string]interface{}{"Id": rolID},
		"Nombre":                      jsonData["Nombre"],
		"Correo":                      jsonData["Correo"],
		"Cedula":                      jsonData["Cedula"],
		"NumeroTelefono":              jsonData["NumeroTelefono"],
		"FotoPerfil":                  jsonData["FotoPerfil"],
		"Activo":                      true,
	}

	usuarioID, err := c.AddUsuarioToCRUD(nuevoUsuario)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"success": false, "status": 500, "message": "Error al crear usuario: " + err.Error()}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = map[string]interface{}{"success": true, "status": 201, "message": "Usuario creado", "id": usuarioID}
	c.ServeJSON()
}

func (c *UsuariosController) PostCredenciales(jsonData map[string]interface{}) (int, error) {
	password, ok := jsonData["Contraseña"].(string)
	if !ok {
		return 0, fmt.Errorf("contraseña inválida o ausente")
	}
	correo, ok := jsonData["Correo"].(string)
	if !ok {
		return 0, fmt.Errorf("correo inválido o ausente")
	}

	hashedPassword, err := services.HashContraseña(password)
	if err != nil {
		return 0, fmt.Errorf("error al hashear contraseña: %v", err)
	}
	credenciales := map[string]interface{}{
		"Contraseña": hashedPassword,
		"Activo":     true,
	}

	credencialesID, err := c.AddCredencialesToCRUD(credenciales)
	if err != nil {
		return 0, fmt.Errorf("error al guardar credenciales: %v", err)
	}

	token := "verificacion_" + strconv.Itoa(credencialesID)
	if err := services.EnviarCorreo(correo, token); err != nil {
		log.Printf("Error al enviar correo: %v", err)
	}

	return credencialesID, nil
}

func GetRolIDFromCRUD(rolNombre string) (int, error) {
	url := fmt.Sprintf("%s/roles?query=Nombre:%s", API_CRUD_Usuarios, rolNombre)
	log.Printf("URL RolID: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error conexión API CRUD: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error leyendo respuesta API: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API CRUD error (%d): %s", resp.StatusCode, string(body))
	}

	var roles []map[string]interface{}
	if err := json.Unmarshal(body, &roles); err != nil {
		return 0, fmt.Errorf("error deserializando respuesta: %v", err)
	}

	if len(roles) == 0 {
		return 0, fmt.Errorf("rol no encontrado: %s", rolNombre)
	}

	idFloat, ok := roles[0]["Id"].(float64)
	if !ok {
		return 0, fmt.Errorf("Id de rol no es numérico")
	}

	return int(idFloat), nil
}

func ValidateRequiredFields(data map[string]interface{}, fields []string) error {
	for _, field := range fields {
		if _, ok := data[field]; !ok {
			return fmt.Errorf("%s", field)
		}
	}
	return nil
}

func (c *UsuariosController) AddUsuarioToCRUD(usuario map[string]interface{}) (int, error) {
	usuarioJSON, err := json.Marshal(usuario)
	if err != nil {
		return 0, fmt.Errorf("error serializando usuario: %v", err)
	}

	resp, err := services.Metodo_post(API_CRUD_Usuarios, "/usuarios", usuarioJSON)
	if err != nil {
		return 0, fmt.Errorf("error petición API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("error API (%d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error leyendo respuesta: %v", err)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		return 0, fmt.Errorf("error deserializando respuesta: %v", err)
	}

	idFloat, ok := res["Id"].(float64)
	if !ok {
		return 0, fmt.Errorf("Id inválido")
	}

	return int(idFloat), nil
}

func (c *UsuariosController) AddCredencialesToCRUD(credenciales map[string]interface{}) (int, error) {
	credencialesJSON, err := json.Marshal(credenciales)
	if err != nil {
		return 0, fmt.Errorf("error serializando credenciales: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/credenciales", API_CRUD_Usuarios), bytes.NewBuffer(credencialesJSON))
	if err != nil {
		return 0, fmt.Errorf("error creando petición: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error ejecutando petición: %v", err)
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
		return 0, fmt.Errorf("Id inválido")
	}

	return int(idFloat), nil
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

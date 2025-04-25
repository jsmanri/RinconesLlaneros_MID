package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
// @Failure 400 Bad Request
// @Failure 500 Internal Server Error
// @router / [post]
func (c *UsuariosController) Post() {
	// Leer el cuerpo de la petición directamente como bytes
	bodyBytes := c.Ctx.Input.RequestBody

	// Validar que el cuerpo de la petición no esté vacío
	if len(bodyBytes) == 0 {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  400,
			"message": "El cuerpo de la petición está vacío.",
		}
		c.ServeJSON()
		return
	}

	// Intentar decodificar el JSON para una validación básica
	var jsonData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  400,
			"message": "Error en el formato de entrada JSON: " + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Validar campos requeridos
	requiredFields := []string{"Nombre", "Correo", "Cedula", "Rol", "Contraseña"}
	if err := ValidateRequiredFields(jsonData, requiredFields); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  400,
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Deserializar el JSON en un struct Usuarios para facilitar el manejo de datos
	var nuevoUsuario Usuarios
	err = json.Unmarshal(bodyBytes, &nuevoUsuario)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  400,
			"message": "Error al deserializar JSON: " + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 1. Guardar Rol
	rolID := 0
	rolNombre, ok := jsonData["Rol"].(string) // Obtener el Rol como string del JSON
	if ok {
		// Llamar a la función para obtener el ID del rol desde la API CRUD
		rolID, err = GetRolIDFromCRUD(rolNombre) // Implementar esta función
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]interface{}{
				"success": false,
				"status":  500,
				"message": "Error al obtener el ID del rol desde el CRUD: " + err.Error(),
			}
			c.ServeJSON()
			return
		}
		nuevoUsuario.Rol = &Roles{Id: rolID} // Asignar el ID del rol a nuevoUsuario
	} else if nuevoUsuario.Rol != nil && nuevoUsuario.Rol.Id != 0 {
		rolID = nuevoUsuario.Rol.Id
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  400,
			"message": "El Rol debe ser un ID existente o un nombre de Rol válido.",
		}
		c.ServeJSON()
		return
	}

	// 2. Guardar Credenciales (Contraseña)
	// En un sistema real, NUNCA guardes la contraseña en texto plano
	hashedPassword, err := HashContraseña(jsonData["Contraseña"].(string)) // Obtener la contraseña del jsonData
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  500,
			"message": "Error al hashear la contraseña: " + err.Error(),
		}
		c.ServeJSON()
		return
	}
	credenciales := Credenciales{Contraseña: hashedPassword, Activo: true}

	credencialesID64, err := AddCredenciales(&credenciales) // Usa la función para agregar credenciales
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  500,
			"message": "Error al guardar credenciales: " + err.Error(),
		}
		c.ServeJSON()
		return
	}
	credencialesID := int(credencialesID64)
	nuevoUsuario.IdCredencialesCredenciales = &Credenciales{Id: credencialesID} // Asigna el ID de credenciales

	nuevoUsuario.Rol = &Roles{Id: rolID}

	// Guardar el usuario en la base de datos
	nuevoUsuario.Activo = true
	usuarioID64, err := AddUsuarios(&nuevoUsuario)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  500,
			"message": "Error al guardar el usuario: " + err.Error(),
		}
		c.ServeJSON()
		return
	}
	usuarioID := int(usuarioID64)

	// Responder al cliente con el ID del nuevo usuario
	c.Ctx.Output.SetStatus(201) // 201 Created
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"status":  201,
		"message": "Usuario creado exitosamente",
		"id":      usuarioID,
	}
	c.ServeJSON()
}

// GetRolIDFromCRUD obtiene el ID del rol desde la API CRUD.
//
//	Debes implementar esta función para hacer la llamada HTTP a tu API CRUD.
func GetRolIDFromCRUD(rolNombre string) (rolID int, err error) {
	// Construir la URL de la API CRUD para obtener el rol por nombre.
	//  Asegúrate de que la URL y los parámetros sean correctos para tu API.
	url := fmt.Sprintf("http://localhost:8081/v1/roles?nombre=%s", rolNombre) // Ejemplo

	// Hacer la petición GET a la API CRUD
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error al hacer la petición a la API CRUD: %v", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta de la API CRUD
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error al leer la respuesta de la API CRUD: %v", err)
	}

	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error de la API CRUD, código de estado: %d, respuesta: %s", resp.StatusCode, string(body))
	}

	// Deserializar la respuesta JSON de la API CRUD en una estructura adecuada.
	//  Aquí se asume que la API CRUD devuelve un JSON con un campo "Id".
	var respuestaRol struct {
		Id int `json:"Id"`
	}
	err = json.Unmarshal(body, &respuestaRol)
	if err != nil {
		return 0, fmt.Errorf("error al deserializar la respuesta de la API CRUD: %v", err)
	}

	return respuestaRol.Id, nil
}

// ValidateRequiredFields Valida que los campos requeridos estén presentes en el JSON.
func ValidateRequiredFields(data map[string]interface{}, fields []string) error {
	for _, field := range fields {
		if _, ok := data[field]; !ok {
			return fmt.Errorf("falta el campo obligatorio: %s", field)
		}
	}
	return nil
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

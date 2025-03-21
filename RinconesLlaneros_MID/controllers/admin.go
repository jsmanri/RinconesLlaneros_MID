package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/astaxie/beego"
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

	// Definir la URL del API CRUD (esto puede provenir de la configuración)
	apiURL := "http://localhost:8082/v1/Usuarios" // Asegúrate de que esta URL esté bien configurada

	// Realizar la solicitud GET al API CRUD
	client := &http.Client{Timeout: 10 * time.Second} // Tiempo de espera para la solicitud
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al crear la solicitud",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Hacer la solicitud y obtener la respuesta
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al realizar la solicitud al API CRUD",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al leer la respuesta del API CRUD",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Estructura que esperamos recibir del API CRUD
	var apiResponse map[string]interface{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error al parsear la respuesta JSON",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Extraer solo los datos de los usuarios desde la respuesta
	usuarios, ok := apiResponse["usuarios consultados"].([]interface{})
	if !ok {
		c.Data["json"] = map[string]interface{}{
			"Message": "Error en el formato de los usuarios en la respuesta",
			"status":  500,
			"success": false,
		}
		c.ServeJSON()
		return
	}

	// Crear un mapa para contar los usuarios registrados por año y mes
	usuariosPorFecha := make(map[string]map[string]int)

	// Crear una lista nueva solo con los campos necesarios (Nombre y Activo)
	var filteredUsuarios []map[string]interface{}
	for _, user := range usuarios {
		userData, ok := user.(map[string]interface{})
		if ok {
			// Crear un mapa con solo los campos "Nombre" y "Activo"
			filteredUser := map[string]interface{}{
				"Nombre": userData["Nombre"],
				"Activo": userData["Activo"],
			}
			filteredUsuarios = append(filteredUsuarios, filteredUser)

			// Obtener la fecha de creación del usuario
			fechaCreacion, ok := userData["FechaCreacion"].(string)
			if ok {
				// Parsear la fecha y extraer el año y el mes
				fecha, err := time.Parse(time.RFC3339, fechaCreacion)
				if err != nil {
					// Si ocurre un error al parsear la fecha, se ignora el usuario
					continue
				}

				// Obtener el año y el mes en formato "YYYY" y "MM"
				anno := fecha.Format("2006")
				mes := fecha.Format("01")

				// Inicializar el mapa para ese año si aún no existe
				if usuariosPorFecha[anno] == nil {
					usuariosPorFecha[anno] = make(map[string]int)
				}

				// Incrementar el contador para ese mes y año
				usuariosPorFecha[anno][mes]++
			}
		}
	}

	// Agregar el conteo de usuarios registrados
	usuarioCount := len(filteredUsuarios)

	// Crear el resultado con los datos de los usuarios por mes y año
	var resultado []map[string]interface{}
	for anno, meses := range usuariosPorFecha {
		for mes, count := range meses {
			resultado = append(resultado, map[string]interface{}{
				"Año":                 anno,
				"Mes":                 mes,
				"UsuariosRegistrados": count,
			})
		}
	}

	// Devolver la respuesta con los usuarios filtrados y el conteo de usuarios por año y mes
	c.Data["json"] = map[string]interface{}{
		"Usuarios":         filteredUsuarios,
		"TotalUsuarios":    usuarioCount,
		"UsuariosPorFecha": resultado,
	}
	c.ServeJSON()
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

}

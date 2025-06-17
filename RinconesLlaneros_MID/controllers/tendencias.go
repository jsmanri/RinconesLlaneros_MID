package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/sena_2824182/RinconesLlaneros_MID/RinconesLlaneros_MID/services"
)

// TendenciasController operations for Tendencias
type TendenciasController struct {
	beego.Controller
}

// URLMapping ...
func (c *TendenciasController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Tendencias
// @Param	body		body 	models.Tendencias	true		"body for Tendencias content"
// @Success 201 {object} models.Tendencias
// @Failure 403 body is empty
// @router / [post]
func (c *TendenciasController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Tendencias by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Tendencias
// @Failure 403 :id is empty
// @router /:id [get]
func (c *TendenciasController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Tendencias
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Tendencias
// @Failure 403
// @router / [get]
func (c *TendenciasController) GetAll() {

	// Obtener JSON de comentarios en una sola consulta
	jsonComentariosStr, err := services.Metodo_get_all("host_api", "Comentarios?limit=0")
	if err != nil {
		fmt.Println("Error al obtener comentarios desde el CRUD:", err)
		c.CustomAbort(500, "Error al obtener los comentarios desde el CRUD")
		return
	}

	// Procesar el JSON de comentarios
	jsonComentarios, err := services.ProcesarJson(jsonComentariosStr)
	if err != nil {
		fmt.Println("Error al procesar los comentarios:", err)
		c.CustomAbort(500, "Error al procesar el JSON de comentarios")
		return
	}

	resultados_parcial := jsonComentarios["comentarios consultados"]

	Arreglo_comentarios, _ := services.ConvertToSliceOfMaps(resultados_parcial)

	groupedItems := services.GroupByID(Arreglo_comentarios)

	var resultado_final []map[string]interface{}

	for idsitio, grupositios := range groupedItems {
		idstring_sitio := fmt.Sprintf("%v", idsitio)
		url := "Sitios_Turisticos/" + idstring_sitio
		sitio, err := services.Metodo_get_all("host_api", url)
		if err != nil {
			fmt.Println("Error al obtener comentarios desde el CRUD:", err)
			c.CustomAbort(500, "Error al obtener los comentarios desde el CRUD")
			return

		}
		jsonsitio, err := services.ProcesarJson(sitio)
		if err != nil {
			fmt.Println("Error al procesar los comentarios:", err)
			c.CustomAbort(500, "Error al procesar el JSON de comentarios")
			return
		}

		jsonsitio_resumido := map[string]interface{}{
			"Id_Sitio":    jsonsitio["sitio consultado"].(map[string]interface{})["Id"],
			"Nombre":      jsonsitio["sitio consultado"].(map[string]interface{})["NombreSitioTuristico"],
			"Descripcion": jsonsitio["sitio consultado"].(map[string]interface{})["DescripcionSitioTuristico"],
			"Fotositio":   jsonsitio["sitio consultado"].(map[string]interface{})["FotoSitio"],
		}

		var ponderacion interface{}
		var ponderacion_total []interface{}
		for _, comentario := range grupositios {

			ponderacion = comentario["Calificacion"]
			ponderacion_total = append(ponderacion_total, ponderacion)
		}

		var suma float64
		for _, ponderacion := range ponderacion_total {

			ponderacion_int, _ := strconv.Atoi(ponderacion.(string))
			suma += float64(ponderacion_int)
		}
		total_ponderacion := suma / float64(len(ponderacion_total))


		jsonsitio_resumido["Ponderacion"] = total_ponderacion
		jsonsitio_resumido["Cantidad_comentarios"] = len(grupositios)

		resultado_final = append(resultado_final, jsonsitio_resumido)

	}
	c.Data["json"] = map[string]interface{}{
		"status":    200,
		"message":   "Consulta realizada correctamente",
		"resultado": resultado_final,
	}

	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Tendencias
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Tendencias	true		"body for Tendencias content"
// @Success 200 {object} models.Tendencias
// @Failure 403 :id is not int
// @router /:id [put]
func (c *TendenciasController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Tendencias
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *TendenciasController) Delete() {

}

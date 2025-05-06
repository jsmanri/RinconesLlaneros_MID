package controllers

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/sena_2824182/RinconesLlaneros_MID/RinconesLlaneros_MID/services"
)

// Sitios_turisticosController operations for Sitios_turisticos
type Sitios_turisticosController struct {
	beego.Controller
}

// URLMapping ...
func (c *Sitios_turisticosController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Sitios_turisticos
// @Param	body		body 	models.Sitios_turisticos	true		"body for Sitios_turisticos content"
// @Success 201 {object} models.Sitios_turisticos
// @Failure 403 body is empty
// @router / [post]
func (c *Sitios_turisticosController) Post() {
	fmt.Println("Metodo Post")

}

// GetOne ...
// @Title GetOne
// @Description get Sitios_turisticos by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Sitios_turisticos
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Sitios_turisticosController) GetOne() {
	fmt.Println("MEtodo GetbyID")

}

// GetAll ...
// @Title GetAll
// @Description get Sitios_turisticos
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Sitios_turisticos
// @Failure 403
// @router / [get]
func (c *Sitios_turisticosController) GetAll() {

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
			"Nombre":      jsonsitio["sitio consultado"].(map[string]interface{})["NombreSitioTuristico"],
			"Descripcion": jsonsitio["sitio consultado"].(map[string]interface{})["DescripcionSitioTuristico"],
			"Fotositio":   jsonsitio["sitio consultado"].(map[string]interface{})["FotoSitio"],
		}

		fmt.Printf("ID %v:\n", idsitio)
		var ponderacion interface{}
		var ponderacion_total []interface{}
		for _, comentario := range grupositios {
			fmt.Printf("  %v\n", comentario)
			fmt.Println("Cantidad comentarios", len(grupositios))


			ponderacion = comentario["Calificacion"]
			fmt.Println("Ponderacion", ponderacion)
			ponderacion_total = append(ponderacion_total, ponderacion)
		}

		fmt.Println("Ponderacion total arreglo", ponderacion_total)
		var suma float64
		for _, ponderacion := range ponderacion_total {
			fmt.Println("tipo dato", reflect.TypeOf(ponderacion))
			fmt.Println("Ponderacion for ", ponderacion)
			ponderacion_int, _ := strconv.Atoi(ponderacion.(string))
			suma += float64(ponderacion_int)
		}
		total_ponderacion := suma/float64(len(ponderacion_total))

		fmt.Println("Ponderacion media", suma)
		fmt.Println("Ponderacion total", total_ponderacion)

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
// @Description update the Sitios_turisticos
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Sitios_turisticos	true		"body for Sitios_turisticos content"
// @Success 200 {object} models.Sitios_turisticos
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Sitios_turisticosController) Put() {
	fmt.Println("Metodo Put")

}

// Delete ...
// @Title Delete
// @Description delete the Sitios_turisticos
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Sitios_turisticosController) Delete() {
	fmt.Println("Metodo Delete")

}

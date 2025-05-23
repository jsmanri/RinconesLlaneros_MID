package controllers

import (
	"fmt"
	"reflect"
	"sort"
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
    idSitio := c.GetString(":id")
    if idSitio == "" {
        c.CustomAbort(400, "ID del sitio no proporcionado")
        return
    }

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

    // Agrupar comentarios por ID de sitio
    groupedItems := services.GroupByID(Arreglo_comentarios)

    // Obtener solo los comentarios del sitio solicitado
    grupositios, existe := groupedItems[idSitio]
    var comentariosDetallados []map[string]interface{}
    if existe {
        for _, comentario := range grupositios {
            autor := comentario["IdUsuario"].(map[string]interface{})["Nombre"]
            texto := comentario["Comentario"].(string)
            calificacionStr := comentario["Calificacion"].(string)
            calificacion, _ := strconv.Atoi(calificacionStr)

            comentarioObj := map[string]interface{}{
                "Autor":        autor,
                "texto":        texto,
                "calificacion": calificacion,
            }
            comentariosDetallados = append(comentariosDetallados, comentarioObj)
        }
    }

    // Obtener datos del sitio turístico específico
    urlSitio := "Sitios_Turisticos/" + idSitio
    sitioStr, err := services.Metodo_get_all("host_api", urlSitio)
    if err != nil {
        fmt.Println("Error al obtener el sitio desde el CRUD:", err)
        c.CustomAbort(500, "Error al obtener el sitio desde el CRUD")
        return
    }

    jsonSitio, err := services.ProcesarJson(sitioStr)
    if err != nil {
        fmt.Println("Error al procesar el JSON del sitio:", err)
        c.CustomAbort(500, "Error al procesar el JSON del sitio")
        return
    }

    // Construir respuesta solo con el sitio solicitado
    jsonsitio_resumido := map[string]interface{}{
        "Id_Sitio":    jsonSitio["sitio consultado"].(map[string]interface{})["Id"],
        "Nombre":      jsonSitio["sitio consultado"].(map[string]interface{})["NombreSitioTuristico"],
        "Descripcion": jsonSitio["sitio consultado"].(map[string]interface{})["DescripcionSitioTuristico"],
        "Fotositio":   jsonSitio["sitio consultado"].(map[string]interface{})["FotoSitio"],
        "Ubicacion":   jsonSitio["sitio consultado"].(map[string]interface{})["Ubicacion"],
        "Telefono":    jsonSitio["sitio consultado"].(map[string]interface{})["IdUsuario"].(map[string]interface{})["NumeroTelefono"],
        "Horario":     jsonSitio["sitio consultado"].(map[string]interface{})["Horario"],
        "Comentarios": comentariosDetallados,
    }

    // Calcular ponderación si hay comentarios
    if len(comentariosDetallados) > 0 {
        var suma float64
        for _, comentario := range comentariosDetallados {
            suma += float64(comentario["calificacion"].(int))
        }
        jsonsitio_resumido["Ponderacion"] = suma / float64(len(comentariosDetallados))
    } else {
        jsonsitio_resumido["Ponderacion"] = 0.0
    }

    jsonsitio_resumido["Cantidad_comentarios"] = len(comentariosDetallados)

    c.Data["json"] = map[string]interface{}{
        "status":    200,
        "message":   "Consulta realizada correctamente",
        "resultado": jsonsitio_resumido,
    }

    c.ServeJSON()
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
		var comentariosDetallados []map[string]interface{}

		for _, comentario := range grupositios {
			autor := comentario["IdUsuario"].(map[string]interface{})["Nombre"]
			texto := comentario["Comentario"].(string)
			calificacionStr := comentario["Calificacion"].(string)
			calificacion, _ := strconv.Atoi(calificacionStr)

			comentarioObj := map[string]interface{}{
				"Autor":        autor,
				"texto":        texto,
				"calificacion": calificacion,
			}
			comentariosDetallados = append(comentariosDetallados, comentarioObj)
		}

		jsonsitio_resumido := map[string]interface{}{
			"Id_Sitio":    jsonsitio["sitio consultado"].(map[string]interface{})["Id"],
			"Nombre":      jsonsitio["sitio consultado"].(map[string]interface{})["NombreSitioTuristico"],
			"Descripcion": jsonsitio["sitio consultado"].(map[string]interface{})["DescripcionSitioTuristico"],
			"Fotositio":   jsonsitio["sitio consultado"].(map[string]interface{})["FotoSitio"],
			"Ubicacion":   jsonsitio["sitio consultado"].(map[string]interface{})["Ubicacion"],
			"Telefono":    jsonsitio["sitio consultado"].(map[string]interface{})["IdUsuario"].(map[string]interface{})["NumeroTelefono"],
			"Horario":     jsonsitio["sitio consultado"].(map[string]interface{})["Horario"],
			"Comentarios": comentariosDetallados,
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
		total_ponderacion := suma / float64(len(ponderacion_total))

		fmt.Println("Ponderacion media", suma)
		fmt.Println("Ponderacion total", total_ponderacion)

		jsonsitio_resumido["Ponderacion"] = total_ponderacion
		jsonsitio_resumido["Cantidad_comentarios"] = len(grupositios)

		resultado_final = append(resultado_final, jsonsitio_resumido)
	}

	sort.Slice(resultado_final, func(i, j int) bool {
		// Comparar la ponderación de cada sitio
		return resultado_final[i]["Ponderacion"].(float64) > resultado_final[j]["Ponderacion"].(float64)
	})
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

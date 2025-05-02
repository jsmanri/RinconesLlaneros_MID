package controllers

import (
	"fmt"
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
    // Obtener JSON de sitios turísticos
    jsonSitiosStr, err := services.Metodo_get_all("host_api", "Sitios_Turisticos")
    if err != nil {
        c.CustomAbort(500, "Error al obtener los sitios turísticos desde el CRUD")
        return
    }

    // Procesar el JSON de sitios turísticos
    jsonSitios, err := services.ProcesarJson(jsonSitiosStr)
    if err != nil {
        c.CustomAbort(500, "Error al procesar el JSON de sitios turísticos")
        return
    }

    // Verificar si "sitios consultados" existe y tiene datos
    sitiosDataInterface, ok := jsonSitios["sitios consultados"]
    if !ok || sitiosDataInterface == nil {
        c.CustomAbort(500, "Error: No se encontraron sitios turísticos")
        return
    }

    // Convertir a slice de interfaces
    sitiosData, ok := sitiosDataInterface.([]interface{})
    if !ok {
        c.CustomAbort(500, "Error: Formato incorrecto de sitios turísticos en el JSON")
        return
    }

    // Obtener JSON de comentarios en una sola consulta
    jsonComentariosStr, err := services.Metodo_get_all("host_api", "Comentarios")
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

    // Verificar si existen comentarios en la API
    comentariosDataInterface, ok := jsonComentarios["comentarios consultados"]
    if !ok || comentariosDataInterface == nil {
        fmt.Println("Advertencia: 'comentarios consultados' no tiene datos. Asignando array vacío.")
        comentariosDataInterface = []interface{}{} // Evita errores si la API devuelve null
    }

    // Convertir a slice de interfaces
    comentariosData, ok := comentariosDataInterface.([]interface{})
    if !ok {
        c.CustomAbort(500, "Error: Formato incorrecto de comentarios en el JSON")
        return
    }

    // Preparar el resultado final
    var resultadoTotal []map[string]interface{}

    // Recorremos cada sitio turístico
    for _, sitio := range sitiosData {
        sitioMap, ok := sitio.(map[string]interface{})
        if !ok {
            continue
        }

        // Obtener ID del sitio turístico
        idSitioInterface, ok := sitioMap["Id"]
        if !ok || idSitioInterface == nil {
            continue
        }

        idSitio := fmt.Sprintf("%v", idSitioInterface) // Convertir a string

        // Variables para acumular la información de los comentarios
        var sumaPuntuaciones float64
        var cantidadComentarios int

        // Filtrar los comentarios por ID de sitio turístico
        for _, comentario := range comentariosData {
            comentarioMap, ok := comentario.(map[string]interface{})
            if !ok {
                continue
            }

            // Obtener el ID del sitio turístico asociado al comentario
            idSitioComentarioInterface, existe := comentarioMap["IdSitiosTuristicos"]
            if !existe || idSitioComentarioInterface == nil {
                continue
            }

            // Extraer correctamente el ID
            idSitioComentarioMap, ok := idSitioComentarioInterface.(map[string]interface{})
            if !ok {
                continue
            }

            idSitioComentario, ok := idSitioComentarioMap["Id"]
            if !ok {
                continue
            }

            // Convertir ambos IDs a string antes de comparar
            if fmt.Sprintf("%v", idSitioComentario) == idSitio {
                // Obtener la calificación y validarla antes de sumarla
                calificacionStr, ok := comentarioMap["Calificacion"].(string)
                if ok {
                    calificacionFloat, err := strconv.ParseFloat(calificacionStr, 64)
                    if err == nil && calificacionFloat > 0 { // Solo incluir valores válidos
                        sumaPuntuaciones += calificacionFloat
                        cantidadComentarios++

                        // Depuración: Mostrar calificación válida procesada
                        fmt.Printf("Calificación procesada para sitio %s: %v\n", idSitio, calificacionFloat)
                    } else {
                        fmt.Printf("Advertencia: Calificación inválida para sitio ID %v\n", idSitio)
                    }
                }
            }
        }

        // **Calcular la media de puntuación con valores seguros**
        var mediaPuntuacion float64
        if cantidadComentarios > 0 {
            mediaPuntuacion = sumaPuntuaciones / float64(cantidadComentarios)
        } else {
            mediaPuntuacion = 0.0 // Evitar valores incorrectos si no hay comentarios
        }

        // Crear el resultado parcial para el sitio turístico
        resultadoParcial := map[string]interface{}{
            "Id":                        sitioMap["Id"],
            "NombreSitioTuristico":      sitioMap["NombreSitioTuristico"],
            "DescripcionSitioTuristico": sitioMap["DescripcionSitioTuristico"],
            "Ubicacion":                 sitioMap["Ubicacion"],
            "Horario":                   sitioMap["Horario"],
            "Dueño":                     sitioMap["IdUsuario"].(map[string]interface{})["Nombre"],
            "Categoria":                 sitioMap["IdCategoria"].(map[string]interface{})["Nombre"],
            "FotoSitio":                 sitioMap["FotoSitio"],
            "Comentarios": map[string]interface{}{
                "Cantidad":        cantidadComentarios,
                "MediaPuntuacion": mediaPuntuacion,
            },
        }

        // Agregar el resultado parcial al resultado total
        resultadoTotal = append(resultadoTotal, resultadoParcial)
    }

    // **Ordenar los sitios turísticos según la media de puntuación**
    sort.Slice(resultadoTotal, func(i, j int) bool {
        return resultadoTotal[i]["Comentarios"].(map[string]interface{})["MediaPuntuacion"].(float64) >
            resultadoTotal[j]["Comentarios"].(map[string]interface{})["MediaPuntuacion"].(float64)
    })

    // **Ya no se limita el número de sitios, el cliente decide cuántos mostrar**

    // Enviar la respuesta final
    c.Data["json"] = map[string]interface{}{
        "success":   len(resultadoTotal) > 0,
        "status":    200,
        "message":   "Consulta realizada correctamente",
        "resultado": resultadoTotal,
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

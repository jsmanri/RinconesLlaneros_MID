package main

import (
	_ "github.com/sena_2824182/RinconesLlaneros_MID/RinconesLlaneros_MID/routers" // Rutas de la aplicación

	"github.com/astaxie/beego"            // Framework Beego
	"github.com/astaxie/beego/plugins/cors" // Plugin CORS para Beego

)

func main() {
	// Configuración de CORS
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		// Permitimos solicitudes de todos los orígenes
		AllowAllOrigins: true,
		// Permitimos los métodos que serán necesarios
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// Permitimos estos encabezados
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		// Permitir credenciales si es necesario
		AllowCredentials: true,
	}))

	// Configuración en modo desarrollo
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// Inicia el servidor
	beego.Run()
}
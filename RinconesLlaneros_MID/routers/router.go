// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/sena_2824182/RinconesLlaneros_MID/RinconesLlaneros_MID/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/creacionusu",
			beego.NSInclude(
				&controllers.CrearusuarioController{},
			),
		),
		beego.NSNamespace("/admin",
			beego.NSInclude(
				&controllers.AdminController{},
			),
		),
		beego.NSNamespace("/newpassword",
			beego.NSInclude(
				&controllers.ActulizarContraseñaController{},
			),
		),
		beego.NSNamespace("/Sesion",
		    beego.NSInclude(
				&controllers.IniciarSesionController{},
			),
		),
	)
	beego.AddNamespace(ns)
}

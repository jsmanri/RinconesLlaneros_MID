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
	beego.Router("/v1/Usuarios", &controllers.UsuariosController{}, "post:Post")
    beego.Router("/v1/admin/usuarios/:id", &controllers.Usuarios_adminController{}, "delete:Delete")
	beego.Router("/v1/usuarios/:id", &controllers.UsuariosController{}, "delete:Delete")
}
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Roles struct {
	Id     int    `orm:"column(Id);pk"`
	Nombre string `orm:"column(Nombre)"`
}

type Credenciales struct {
	Id       int    `orm:"column(Id);pk"`
	Usuario  string `orm:"column(Usuario)"`
	Password string `orm:"column(Password)"`
}

type Usuarios struct {
	Id                         int           `orm:"column(Id_Usuario);pk;auto"`
	Nombre                     string        `orm:"column(Nombre)"`
	Rol                        *Roles        `orm:"column(Rol);rel(fk)"`
	Correo                     string        `orm:"column(Correo)"`
	Cedula                     string        `orm:"column(Cedula)"`
	NumeroTelefono             string        `orm:"column(Numero_Telefono)"`
	FotoPerfil                 string        `orm:"column(Foto_Perfil);type(text);null"`
	FechaCreacion              time.Time     `orm:"column(Fecha_Creacion);type(timestamp with time zone);auto_now_add"`
	FechaModificacion          time.Time     `orm:"column(Fecha_Modificacion);type(timestamp with time zone);auto_now"`
	Activo                     bool          `orm:"column(Activo)"`
	IdCredencialesCredenciales *Credenciales `orm:"column(Id_Credenciales);rel(fk)"`
}

// Registrar el modelo en Beego ORM
func init() {
	orm.RegisterModel(new(Usuarios), new(Roles), new(Credenciales))
}

package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type Categorias struct {
	Id                int       `orm:"column(Id_categoria);pk;auto"`
	Nombre            string    `orm:"column(Nombre)"`
	Descripcion       string    `orm:"column(Descripcion)"`
	FechaCreacion     time.Time `orm:"column(Fecha_creacion);type(timestamp with time zone);auto_now_add"`
	FechaModificacion time.Time `orm:"column(Fecha_modificacion);type(timestamp with time zone);auto_now"`
	Activo            bool      `orm:"column(Activo)"`
}

func init() {
	orm.RegisterModel(new(Categorias))
}

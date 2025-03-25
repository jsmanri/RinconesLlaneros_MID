package models

import (
	"time"
)

type SitiosTuristicos struct {
	Id                        int         `orm:"column(Id_sitio_turistico);pk;auto"`
	NombreSitioTuristico      string      `orm:"column(Nombre_Sitio_turistico)"`
	DescripcionSitioTuristico string      `orm:"column(Descripcion_Sitio_turistico)"`
	Ubicacion                 string      `orm:"column(Ubicacion)"`
	Horario                   string      `orm:"column(Horario)"`
	IdUsuario                 *Usuarios   `orm:"column(Id_Usuario);rel(fk)"`
	IdCategoria               *Categorias `orm:"column(Id_Categoria);rel(fk)"`
	FotoSitio                 string      `orm:"column(Foto_Sitio);type(text)"`
	FechaCreacion             time.Time   `orm:"column(Fecha_creacion);type(timestamp with time zone);auto_now_add"`
	FechaModifcacion          time.Time   `orm:"column(Fecha_modifcacion);type(timestamp with time zone);auto_now"`
	Activo                    bool        `orm:"column(Activo)"`
}

package models

import (
	"errors"

	"github.com/astaxie/beego/orm"
)

type SitiosTuristicos struct {
	Id        int    `orm:"column(id);pk;auto"`
	Nombre    string `orm:"column(nombre)"`
	UsuarioID int    `orm:"column(usuario_id)"`
}

func ObtenerSitiosPorUsuario(userID int) ([]SitiosTuristicos, error) {
	o := orm.NewOrm()
	var sitios []SitiosTuristicos

	_, err := o.Raw("SELECT * FROM sitios WHERE usuario_id = ?", userID).QueryRows(&sitios)
	if err != nil {
		return nil, errors.New("error al obtener los sitios turísticos del usuario")
	}

	return sitios, nil
}

type SitioTuristico struct {
	Id        int    `orm:"auto;pk"`
	Nombre    string `orm:"size(100)"`
	Ubicacion string `orm:"size(200)"`
	UsuarioId int    `orm:"index"`
}

func init() {
	orm.RegisterModel(new(SitioTuristico))
}
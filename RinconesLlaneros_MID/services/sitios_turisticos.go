package services

import (
	"github.com/astaxie/beego/orm"
	"github.com/sena_2824182/RinconesLlaneros_MID/RinconesLlaneros_MID/models"
)

// Obtener sitios turísticos por usuario
func ObtenerSitiosPorUsuario(idUsuario int) ([]models.SitiosTuristicos, error) {
	o := orm.NewOrm()
	var sitios []models.SitiosTuristicos
	_, err := o.QueryTable(new(models.SitiosTuristicos)).Filter("IdUsuario", idUsuario).All(&sitios)
	return sitios, err
}

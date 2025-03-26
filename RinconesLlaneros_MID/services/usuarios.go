package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/sena_2824182/RinconesLlaneros_MID/RinconesLlaneros_MID/models"
)

// Desactivar usuario y sus sitios turísticos
var servicioCRUD = beego.AppConfig.String("Servicio_CRUD")

// DesactivarUsuario realiza el borrado lógico del usuario y sus sitios
func DesactivarUsuario(idUsuario int) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return err
	}

	// Desactivar sitios turísticos
	_, err = o.QueryTable(new(models.SitiosTuristicos)).Filter("IdUsuario", idUsuario).Update(orm.Params{"Activo": false})
	if err != nil {
		o.Rollback()
		return err
	}

	// Desactivar usuario
	_, err = o.QueryTable(new(models.Usuarios)).Filter("Id", idUsuario).Update(orm.Params{"Activo": false})
	if err != nil {
		o.Rollback()
		return err
	}

	return o.Commit()
}

// ObtenerSitiosPorUsuario devuelve los sitios turísticos de un usuario desde el CRUD local
func ObtenerSitiosPorUsuario(idUsuario int) ([]string, error) {
	o := orm.NewOrm()
	var sitios []string

	query := fmt.Sprintf("%s/sitios_turisticos?usuario_id=%d", servicioCRUD, idUsuario)

	// Simulación de consulta al CRUD local
	_, err := o.Raw(query).QueryRows(&sitios)
	if err != nil {
		return nil, errors.New("error al obtener los sitios turísticos")
	}

	return sitios, nil
}

type Usuario struct {
	Id        int       `orm:"auto;pk"`
	Nombre    string    `orm:"size(100)"`
	Email     string    `orm:"size(100);unique"`
	Eliminado bool      `orm:"default(false)"` // Borrado lógico
	FechaBaja time.Time `orm:"null"`
}

func init() {
	// Registrar modelo
	orm.RegisterModel(new(Usuario))
}

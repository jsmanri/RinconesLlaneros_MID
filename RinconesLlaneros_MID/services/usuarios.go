package services

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/sena_2824182/RinconesLlaneros_MID/RinconesLlaneros_MID/models"
)

// Desactivar usuario y sus sitios turísticos
func DesactivarUsuario(idUsuario int) error {
	o := orm.NewOrm()

	// Iniciar transacción
	err := o.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar la transacción: %v", err)
	}

	// Desactivar los sitios turísticos asociados al usuario
	_, err = o.QueryTable(new(models.SitiosTuristicos)).Filter("IdUsuario", idUsuario).Update(orm.Params{"Activo": false})
	if err != nil {
		o.Rollback()
		return fmt.Errorf("error al desactivar sitios turísticos del usuario %d: %v", idUsuario, err)
	}

	// Desactivar el usuario
	_, err = o.QueryTable(new(models.Usuarios)).Filter("Id", idUsuario).Update(orm.Params{"Activo": false})
	if err != nil {
		o.Rollback()
		return fmt.Errorf("error al desactivar usuario %d: %v", idUsuario, err)
	}

	// Confirmar transacción
	if err := o.Commit(); err != nil {
		return fmt.Errorf("error al confirmar la transacción: %v", err)
	}

	return nil
}

// AutoEliminarUsuario - Permite que un usuario elimine su cuenta y sus sitios turísticos
func AutoEliminarUsuario(idUsuario int) error {
	return DesactivarUsuario(idUsuario) // Reutilizamos la función
}

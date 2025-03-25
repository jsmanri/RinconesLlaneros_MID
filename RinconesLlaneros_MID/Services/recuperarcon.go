package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"
)

// URL del servicio usuario (Evita repetir la URL en el código)
const ServicioUsuarioUP = "http://localhost:8082/v1/Usuarios"

// Generar código aleatorio de 6 dígitos
func GenerarCodigo() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// Enviar código por correo
func EnviarCodigoPorCorreo(codigo, correo string) error {
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"
	from := "tucorreo@gmail.com"
	password := "tucontraseña"

	subject := "Código de verificación"
	body := fmt.Sprintf("Tu código de verificación es: %s", codigo)

	msg := []byte("To: " + correo + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpServer)

	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{correo}, msg)
	return err
}

// Obtener usuario por ID desde el servicio de usuarios
func ObtenerUsuarioPorID(id int) (*map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%d", ServicioUsuarioUP, id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error al obtener usuario, estado HTTP: %d", resp.StatusCode)
	}

	var usuario map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &usuario)
	if err != nil {
		return nil, err
	}

	return &usuario, nil
}

// Guardar código de recuperación en la base de datos
func GuardarCodigoEnDB(id int, codigo string) error {
	url := fmt.Sprintf("%s/%d", ServicioUsuarioUP, id)

	data := map[string]interface{}{
		"CodigoRecuperacion": codigo,
		"CodigoExpira":       time.Now().Add(10 * time.Minute).Format(time.RFC3339),
	}

	usuarioJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(usuarioJson))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error al guardar código en DB, estado HTTP: %d", resp.StatusCode)
	}

	return nil
}

// Obtener código guardado en la base de datos
func ObtenerCodigoGuardado(id int) (string, error) {
	url := fmt.Sprintf("%s/%d", ServicioUsuarioUP, id)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error al obtener código, estado HTTP: %d", resp.StatusCode)
	}

	var respuesta struct {
		CodigoRecuperacion string `json:"CodigoRecuperacion"`
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		return "", err
	}

	return respuesta.CodigoRecuperacion, nil
}

// Actualizar la contraseña del usuario en la base de datos
func ActualizarContraseñaUsuario(id int, nuevaContraseña string) error {
	url := fmt.Sprintf("%s/%d", ServicioUsuarioUP, id)

	// Aquí deberías tener tu función para hashear la contraseña
	hashedPassword := nuevaContraseña // Cambia esto por la función real

	data := map[string]interface{}{
		"IdCredencialesCredenciales": map[string]interface{}{
			"Contraseña": hashedPassword,
			"Activo":     true,
		},
	}

	usuarioJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(usuarioJson))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error al actualizar contraseña, estado HTTP: %d", resp.StatusCode)
	}

	return nil
}

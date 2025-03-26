package servicios

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Id int `json:"Id"`
	jwt.RegisteredClaims
}

// Definir la URL base del servicio de usuarios
const ServicioUsuarioUP = "http://localhost:8082/usuarios"

// ActualizarContraseñaUsuario actualiza la contraseña de un usuario
func ActualizarContraseñaUsuario(id int, nuevaContraseña string) error {
	// Crear el JSON con la nueva contraseña
	data := map[string]interface{}{
		"IdCredencialesCredenciales": map[string]interface{}{
			"Contraseña": nuevaContraseña,
			"Activo":     true,
		},
	}

	usuarioJson, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error al serializar JSON: %w", err)
	}

	url := fmt.Sprintf("%s/%d", ServicioUsuarioUP, id)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(usuarioJson))
	if err != nil {
		return fmt.Errorf("error al crear solicitud HTTP: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al enviar solicitud HTTP: %w", err)
	}
	defer resp.Body.Close()

	// Leer respuesta en caso de error
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error al actualizar contraseña, estado HTTP: %d, respuesta: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Clave secreta (debería ser segura y venir de variables de entorno)
var jwtSecret = []byte("clave_secreta")

// Generar un token de recuperación de contraseña
func GenerarTokenRecuperacion(Id int) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Id: Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Validar token de recuperación
func ValidarTokenRecuperacion(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, errors.New("error al analizar el token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	// Verificar expiración
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expirado")
	}

	return claims, nil
}

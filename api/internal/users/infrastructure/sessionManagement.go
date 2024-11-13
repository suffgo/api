package infrastructure

import (
	"log"
	"net/http"
	"strconv"
	sv "suffgo/internal/shared/domain/valueObjects"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func createSession(userID sv.ID, c echo.Context) error {
	// Crear la sesión
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 2, // Duración de la sesión en segundos (2 días en este caso)
		HttpOnly: true,
		Secure:   false, // Establecer en 'false' si estás usando 'http'
		SameSite: http.SameSiteLaxMode,
	}
	// Convertir el userID a string antes de almacenarlo
	sess.Values["user_id"] = strconv.FormatUint(uint64(userID.Id), 10)
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		log.Printf("Error al guardar la sesión: %v", err)
		return err
	}
	log.Println("Sesión guardada exitosamente")

	return nil
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			log.Printf("Error al obtener la sesión: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error al obtener la sesión"})
		}

		log.Printf("Valores de la sesión: %#v", sess.Values)

		userID, ok := sess.Values["user_id"].(string)
		if !ok || userID == "" {
			log.Println("Usuario no autenticado: user_id no encontrado en la sesión")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "usuario no autenticado"})
		}

		c.Set("user_id", userID)

		return next(c)
	}
}

func Logout(c echo.Context) error {
	
	sess, err := session.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al obtener la sesión"})
	}

	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	
	return nil
}

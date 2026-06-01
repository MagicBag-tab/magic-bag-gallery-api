package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const RoleKey contextKey = "role"
const TipoEmpleadoKey contextKey = "tipo_empleado"
const NombreKey contextKey = "nombre"

const sessionName = "magic_bag_session"

var store *sessions.CookieStore

func InitSessionStore() {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = os.Getenv("JWT_SECRET")
	}
	if secret == "" {
		secret = "magic-bag-dev-session-secret"
	}

	store = sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	}
}

func getStore() *sessions.CookieStore {
	if store == nil {
		InitSessionStore()
	}
	return store
}

func Session(r *http.Request) (*sessions.Session, error) {
	return getStore().Get(r, sessionName)
}

func SaveSessionUser(w http.ResponseWriter, r *http.Request, idUsuario int, role, nombre, tipoEmpleado string) error {
	session, err := Session(r)
	if err != nil {
		return err
	}
	session.Values["user_id"] = idUsuario
	session.Values["role"] = role
	session.Values["nombre"] = nombre
	session.Values["tipo_empleado"] = tipoEmpleado
	return session.Save(r, w)
}

func ClearSession(w http.ResponseWriter, r *http.Request) error {
	session, err := Session(r)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Session(r)
		if err != nil {
			http.Error(w, "Sesion invalida", http.StatusUnauthorized)
			return
		}

		idUsuario, ok := session.Values["user_id"].(int)
		if !ok || idUsuario == 0 {
			http.Error(w, "Sesion requerida", http.StatusUnauthorized)
			return
		}

		role, ok := session.Values["role"].(string)
		if !ok || role == "" {
			http.Error(w, "Sesion requerida", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, idUsuario)
		ctx = context.WithValue(ctx, RoleKey, role)
		ctx = context.WithValue(ctx, NombreKey, session.Values["nombre"])
		ctx = context.WithValue(ctx, TipoEmpleadoKey, session.Values["tipo_empleado"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(RoleKey).(string)
			if !ok || userRole != role {
				http.Error(w, "Acceso denegado", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

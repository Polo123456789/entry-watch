# PLAN.md - Sistema de Autenticación entry-watch

## Visión General

Implementar el sistema de autenticación completo en `internal/http/auth/` siguiendo los patrones del proyecto de referencia.

## Especificaciones Técnicas

- **Sesiones**: gorilla/sessions con CookieStore
- **Passwords**: bcrypt de golang.org/x/crypto
- **Base de datos**: SQLC para queries generados
- **Bootstrap**: Automático al arranque con valores por defecto
- **Login**: Por email, redirige según rol
- **Rutas**: `/auth/login`, `/auth/logout`

## Decisiones de Diseño

1. **Login único**: Todos los roles usan `/auth/login`, el handler redirige según el rol del usuario a:
   - `/super/` → RoleSuperAdmin
   - `/admin/` → RoleAdmin
   - `/guard/` → RoleGuard
   - `/neighbor/` → RoleUser

2. **Superadmin inicial**:
   - Email: `admin@localhost`
   - Password: `changeme`
   - Se crea automáticamente al arranque si no existe ningún superadmin habilitado

3. **Estructura centralizada**: Todo el auth vive en `internal/http/auth/` (no distribuido por roles como en el reference)

4. **Sesiones**: Cookie HTTP-only, secure, 12 horas de duración

---

# Current Progress

[INSTRUCCIONES PARA AGENTES]

- Trabaja en UNA sola tarea a la vez (secuencial, nunca paralelo)
- Eres libre de elegir la tarea más importante en tu sesión
- Crea sub-tareas adicionales si es necesario
- Documenta blockers claramente para el siguiente agente
- Haz commits regulares con mensajes descriptivos
- Trabaja en la rama actual (no crees nuevas ramas)
- Cuando no queden tareas pendientes, crea un archivo `DONE` en la raíz

## Current

*No hay tareas en progreso actualmente*

## Pending

*No hay tareas pendientes - el sistema de autenticación está completo*

## Completed

* ✅ Crear queries SQLC para usuarios en db/sqlc/users.sql
    * Subtasks:
        * ✅ Crear archivo `db/sqlc/users.sql` con queries necesarias:
            * `GetUserByEmail` - Buscar usuario por email
            * `GetUserByID` - Buscar usuario por ID
            * `CreateUser` - Crear nuevo usuario
            * `CountSuperAdmins` - Contar superadmins habilitados
            * `UpdateUserPassword` - Actualizar password
        * ✅ Verificar que el schema de la tabla users tiene todos los campos necesarios
        * ✅ Asegurar que los tipos coincidan con el domain model

* ✅ Generar código SQLC con make sqlc
    * Subtasks:
        * ✅ Ejecutar `make sqlc` para generar código
        * ✅ Verificar que se generó `internal/sqlc/users.sql.go`
        * ✅ Revisar que los tipos generados son correctos

* ✅ Crear UserStore en internal/sqlc/user-store.go
    * Subtasks:
        * ✅ Crear struct `UserStore` que envuelva las queries SQLC generadas
        * ✅ Implementar método `GetByEmail(ctx, email) (entry.UserWithPassword, bool, error)`
        * ✅ Implementar método `GetByID(ctx, id) (*entry.User, bool, error)`
        * ✅ Implementar método `CreateUser(ctx, email, firstName, lastName, user, passwordHash) error`
        * ✅ Implementar método `CountSuperAdmins(ctx) (int64, error)`
        * ✅ Crear función constructor `NewUserStore(db) *UserStore`

* ✅ Crear internal/http/auth/store.go con interface UserStore
    * Subtasks:
        * ✅ Definir interface `UserStore` con los métodos necesarios
        * ✅ Asegurar que la interface esté en el paquete auth (no en sqlc)
        * ✅ Documentar la interface
        * ✅ Usar `entry.UserWithPassword` (movido a entry package para evitar ciclos de import)

* ✅ Crear internal/http/auth/handlers.go (login/logout)
    * Subtasks:
        * ✅ Crear `hGetLogin(session, logger)` - Muestra formulario
        * ✅ Crear `hPostLogin(session, store, logger)` - Valida credenciales
            * Usar `bcrypt.CompareHashAndPassword` para validar
            * Crear sesión con `gorilla/sessions`
            * Redirigir según rol del usuario
        * ✅ Crear `hGetLogout(session)` - Invalida sesión
        * ✅ Implementar helper `attemptLogin()` para validación de credenciales
        * ✅ Crear helper `setCurrentUser()` para guardar en sesión
        * ✅ Crear helper `getRedirectForRole()` para redirección según rol
        * ✅ Usar `entry.UserSafeError` para mensajes de error (movido a entry package)

* ✅ Crear internal/http/auth/middleware.go (AuthMiddleware)
            * Recuperar user ID de la sesión
            * Buscar usuario en DB
            * Inyectar usuario en contexto con `entry.WithUser()`
            * Si no hay sesión, redirigir a `/auth/login`
        * [ ] Crear helper `CurrentUser(session, r) (*entry.User, bool)`
        * [ ] Crear helper `RedirectIfAuthenticated` (para página de login)
        * [ ] Definir constante `authSessionKey = "entry-watch-auth"`

* Crear internal/http/auth/bootstrap.go (EnsureSuperAdminExists)
    * Subtasks:
        * [ ] Crear función `EnsureSuperAdminExists(store UserStore, logger) error`
        * [ ] Verificar si existe algún superadmin habilitado (`CountSuperAdmins`)
        * [ ] Si no existe, crear superadmin con:
            * Email: `admin@localhost`
            * Password: `changeme` (hash con bcrypt)
            * FirstName: `Super`
            * LastName: `Admin`
            * Role: `RoleSuperAdmin`
            * Enabled: `true`
            * CondominiumID: `nil/0` (superadmin no pertenece a un condominio específico)
        * [ ] Loggear cuando se crea el superadmin inicial (WARN nivel)

* Crear internal/http/auth/routes.go
    * Subtasks:
        * [ ] Crear función `Handle(app *entry.App, logger *slog.Logger, session sessions.Store) http.Handler`
        * [ ] Configurar rutas:
            * `/auth/login` GET → hGetLogin
            * `/auth/login` POST → hPostLogin
            * `/auth/logout` GET → hGetLogout
        * [ ] Crear userStore con `sqlc.NewUserStore(app.Store.DB())`
        * [ ] Aplicar middleware `RedirectIfAuthenticated` solo a login GET

* Crear template login.templ en internal/templates/auth/
    * Subtasks:
        * [ ] Crear `internal/templates/auth/login.templ`
        * [ ] Diseño consistente con templates existentes (usar base layout)
        * [ ] Formulario con campos: email, password
        * [ ] Mostrar mensaje de error si existe (query param `?error=1`)
        * [ ] Link/botón de submit
        * [ ] Generar con `make templates`

* Actualizar internal/http/routes.go para incluir /auth/
    * Subtasks:
        * [ ] Descomentar línea `// TODO: mux.Handle("/auth/", auth.Handle(app, logger))`
        * [ ] Pasar session store a auth.Handle()
        * [ ] Importar paquete auth

* Actualizar internal/http/middleware.go (reemplazar mock por AuthMiddleware)
    * Subtasks:
        * [ ] Reemplazar el mock de usuario (líneas 28-34) con llamada a `auth.AuthMiddleware`
        * [ ] El middleware debe recibir session store
        * [ ] Mantener `CanonicalLoggerMiddleware` para logging pero sin mock
        * [ ] Asegurar que el orden es correcto: AuthMiddleware → CanonicalLoggerMiddleware → handler

* Actualizar internal/entry/app.go (agregar UserStore a interface Store)
    * Subtasks:
        * [ ] Agregar `UserStore` a la interface `Store`
        * [ ] Verificar que `sqlc.Store` implementa todos los métodos necesarios
        * [ ] Agregar método `DB() *sql.DB` a la interface si es necesario para crear UserStore

* Actualizar cmd/entry-watch/main.go (llamar EnsureSuperAdminExists)
    * Subtasks:
        * [ ] Importar `github.com/Polo123456789/entry-watch/internal/http/auth`
        * [ ] Crear session store (gorilla/sessions CookieStore)
        * [ ] Llamar `auth.EnsureSuperAdminExists(store, logger)` antes de iniciar servidor
        * [ ] Pasar session store a `http.NewServer()` para que esté disponible en routes
        * [ ] Configurar session store con:
            * Secret: hardcoded para MVP ("entry-watch-secret-change-in-prod")
            * MaxAge: 12 horas
            * HttpOnly: true
            * Secure: false (true en producción, detectar via DEBUG flag)

* Ejecutar make templates y make sqlc
    * Subtasks:
        * [ ] `make sqlc` - Generar queries SQLC
        * [ ] `make templates` - Generar templates auth
        * [ ] Verificar archivos generados

* Ejecutar make tidy y make audit
    * Subtasks:
        * [ ] `make tidy` - Formatear y ordenar imports
        * [ ] `make audit` - Ejecutar lint y tests
        * [ ] Corregir cualquier error encontrado

## Completed

*Aún no hay tareas completadas*

---

# Lessons Learned

[ESPACIO PARA DOCUMENTAR APRENDIZAJES]

## Blockers Resueltos

*Documentar aquí los bloqueos encontrados y sus soluciones*

## Patrones del Proyecto de Referencia

- **Password hashing**: Usar `golang.org/x/crypto/bcrypt` con `bcrypt.DefaultCost`
- **Sessions**: `github.com/gorilla/sessions` CookieStore, no JWT
- **Store pattern**: Interface en `http/auth/` + implementación en `internal/sqlc/`
- **Error handling**: Crear errores user-safe para mensajes en español
- **Middleware order**: AuthMiddleware → Logger → Recover → Handler
- **User injection**: Usar `entry.WithUser()` y `entry.UserFromCtx()`

## Convenciones del Proyecto

- Agrupar imports: stdlib, terceros, internos (separados por línea en blanco)
- Handlers devuelven `error` y usan `util.Handler()` wrapper
- No usar emojis en código ni comentarios (a menos que se pida explícitamente)
- Ejecutar `make tidy` antes de commits
- Ejecutar `make audit` para validar calidad

## Estructura de Tipos

- Domain types: `internal/entry/` (User, Visit, Condominium)
- Store interfaces: `internal/http/auth/` (UserStore)
- Store implementations: `internal/sqlc/` (user-store.go)
- HTTP handlers: `internal/http/auth/` (handlers.go, middleware.go)
- Templates: `internal/templates/auth/` (login.templ)

## Librerías Necesarias

Agregar a go.mod:
- `github.com/gorilla/sessions v1.4.0` (sessions)
- `golang.org/x/crypto v0.41.0` (bcrypt)

Verificar si ya están en go.mod del reference-proyect.

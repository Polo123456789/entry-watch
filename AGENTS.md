
AGENTS.md — Guía para agentes
--------------------------------

Estas instrucciones están pensadas para agentes autónomos que trabajen en este repositorio (entry-watch). Incluyen comandos de build/lint/test, convenciones de estilo y reglas específicas del proyecto. Sigue las normas de seguridad y control de cambios del repositorio.

1) Comandos principales (build / lint / test)
------------------------------------------------
- Formateo y tidy: `make tidy` (ejecuta `go fmt ./...` y `go mod tidy`).
- Generar plantillas: `make templates` (usa `go tool templ generate -lazy`).
- Generar cliente SQL (sqlc): `make sqlc` (ejecuta `go tool sqlc generate`).
- Compilar: `make build` (genera binario en `/tmp/bin/entry-watch`).
- Ejecutar: `make run` (construye y ejecuta el binario).
- Ejecutar con recarga en caliente: `make run/live` (requiere `air`).
- Migraciones:
  - Crear: `make migration/create` (interactivo).
  - Aplicar: `make migration/up`.
  - Revertir: `make migration/down`.
- Calidad: `make audit` (ejecuta `go vet`, `golangci-lint`, `govulncheck` y tests según Makefile).
- Tests: `make test` (ejecuta `go test -tags assert -race -buildvcs ./...`).

Ejecutar un solo test
----------------------
- En el paquete actual:
  - `go test -run '^TestMyCase$' -v`.
- Desde la raíz para un paquete concreto:
  - `go test ./internal/http -run '^TestDashboard_Post$' -v`.
- Con tags usados por el repo (Makefile usa `-tags assert`):
  - `go test -tags assert ./internal/http -run '^TestDashboard_Post$' -v`

2) Generación y artefactos
---------------------------
- Templates: fuentes en `internal/templates/**/*.templ` — archivos generados terminan en `_templ.go`. No editar manualmente los `.templ` generados por `go tool templ`.
- SQLC: fuentes en `db/sqlc/*.sql` y configuración en `sqlc.yaml`. Código generado en `internal/sqlc` — NO editar archivos generados. Cualquier implementación manual deberá colocarse junto a ellos en archivos con nombres distintos (por ejemplo `internal/sqlc/store.go`, `internal/sqlc/user_store.go`).

3) Dónde viven los tipos y responsabilidades (concreción)
-------------------------------------------------------
- Tipos de dominio y lógica de negocio: `internal/entry` (ej.: `internal/entry/visit.go`, `internal/entry/user.go`, `internal/entry/app.go`).
- Autenticación y gestión de usuarios (CRUD, login/logout, middleware): `internal/http/auth` (handlers y middleware). Si no existe todavía, ese es el lugar recomendado.
-- Handlers y vistas por tipo de usuario (UI): cada rol tiene su propio paquete bajo `internal/http/` — por ejemplo `internal/http/user`, `internal/http/superadmin`, `internal/http/admin`, `internal/http/guard`. Cada paquete contiene las vistas/handlers específicos de ese rol y debe consumir servicios de dominio o `internal/http/auth` para la lógica compartida.
- Código sqlc generado: `internal/sqlc` — mantener generado y no editar; implementaciones manuales y wrappers deben vivir "junto a" los generados.

4) Convenciones de estilo y prácticas (específicas)
--------------------------------------------------
Formato e imports
- Usa `gofmt`/`go fmt` siempre. Ejecuta `make tidy` antes de crear commits.
- Agrupa imports en 3 bloques separados por línea en blanco: stdlib, terceros, imports internos (ej. `github.com/Polo123456789/entry-watch/...`).

Nombres y tipos
- Exportados: UpperCamelCase (`NewServer`, `VisitStore`). No exportes cosas innecesarias.
- No exportados: lowerCamelCase o snake style según Go (`userStore`, `hashPassword`).
- Interfaces: usar nombres descriptivos (`VisitStore`, `CondominiumStore`). Evitar sufijos `-er` si no aporta claridad.
- Structs: usar punteros como receptor cuando el método muta estado; usar valores para tipos pequeños e inmutables.

Funciones y contexto
- Para funciones que operan con request-scoped data, el primer parámetro debe ser `ctx context.Context`.
- Constructores: `NewXxx(...)` y devolver puntero cuando la instancia es mutable o costosa.

Errores
- Siempre devolver `error` como último valor. Wrappea errores con `fmt.Errorf("...: %w", err)` para preservar cadena.
- Usa tipos de error concretos para distinguir casos (p. ej. `UnauthorizedError`, `ForbiddenError`) y `errors.Is` / `errors.As` para comprobarlos.
- Evita exponer mensajes internos al cliente; registra el error completo y devuelve mensajes genéricos HTTP cuando corresponda.

HTTP / handlers
- Sigue el patrón local: handlers que devuelven `error` y son adaptados por `internal/http/util.Handler(logger, func(w,r) error)`.
- Usa `context` para inyectar `entry.User` con `entry.WithUser(ctx, user)` y los helpers `entry.RequireRole(...)`.
- Middleware de autenticación debe validar sesión/token y poner el usuario en el contexto; no inyectes usuarios mock (remplazar `CanonicalLoggerMiddleware` mock en el futuro).

Pruebas
- Prefiere pruebas table-driven y testea la lógica de negocio con unit tests.
- Usa `t.Parallel()` en tests independientes.
- Para integraciones con sqlite usa DB temporales o `:memory:` y aplicar migraciones en setup.

5) Regla SQLC y archivos generados
---------------------------------
- NUNCA editar archivos generados por sqlc en `internal/sqlc`.
- Todas las implementaciones manuales (wrappers, adaptadores, store implementations) deben vivir "junto a" los archivos generados en `internal/sqlc` pero en archivos con nombres distintos (`*_impl.go`, `store.go`, etc.).

6) Reglas de Git y operaciones seguras
------------------------------------
- Nunca ejecutar comandos destructivos sin confirmación (ej. `git reset --hard`).
- No cambiar head o hacer force push a ramas protegidas sin permiso explícito.
- No crear commits automáticos a menos que el usuario lo solicite; si se te pide, sigue el flujo: `make tidy`, `make audit`, añadir tests y luego crear commit con mensaje claro.

7) Cursor / Copilot rules
-------------------------
- No se han encontrado reglas Cursor (`.cursor/rules/`) ni archivos de instrucciones para Copilot (`.github/copilot-instructions.md`). Si aparecen posteriormente, deben seguirse.

8) Checklist rápido para agentes antes de crear un PR
----------------------------------------------------
- Ejecutar: `make templates` y `make sqlc` si cambias templates o SQL.
- Ejecutar: `make tidy` y `make audit`.
- Ejecutar: `go test ./...` y arreglar fallos.
- Verificar que no se modificaron archivos generados por sqlc.

9) Línea corta de referencia (para no perderse)
------------------------------------------------
- Tipos de dominio: `internal/entry`.
- Auth (CRUD + login/logout): `internal/http/auth`.
- Código sqlc generado: `internal/sqlc` — NO editar; implementaciones manuales junto a los generados.

Si necesitas más detalle (ej. ejemplos de código, comandos de single-test con tags, o reglas de estilo más estrictas) pídelo y amplio este archivo.

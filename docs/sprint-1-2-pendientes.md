# Pendientes Sprint 1 y 2

Este documento agrupa deudas tecnicas, validaciones pendientes y mejoras posibles detectadas durante el avance de Sprint 1 y Sprint 2.

## Sprint 1

### Autenticacion y seguridad

- Revisar politica final de expiracion de JWT para ambientes productivos.
- Definir estrategia de refresh tokens si el MVP requiere sesiones largas.
- Agregar rotacion de secretos JWT en una etapa posterior.
- Agregar rate limiting para `/api/v1/auth/login`.
- Definir bloqueo temporal por intentos fallidos de login.
- Revisar si `super_admin` requiere endpoints separados o si se mantiene fuera del MVP inicial.

### Usuarios y companias

- Crear use cases de companies/users si se necesitan endpoints administrativos internos.
- Definir flujo seguro para crear usuarios desde un panel administrativo.
- Agregar validaciones de normalizacion de email en capa application, no solo en base de datos/repositorio.
- Definir comportamiento de usuarios inactivos asociados a drivers.

### Testing

- Agregar tests de integracion de login contra base de datos local con `TEST_DATABASE_URL`.
- Agregar casos de expiracion real de token JWT usando tiempos controlados.
- Agregar pruebas de permisos cruzados entre tenants para usuarios.

## Sprint 2

### Vehicles

- Definir si `DELETE /api/v1/vehicles/{id}` debe quedar como inactivacion permanente o si existira flujo de reactivacion.
- Evaluar filtros de listado por `status`, `vehicle_type` y busqueda por placa.
- Agregar paginacion cuando el volumen de flota crezca.
- Definir catalogo controlado para `vehicle_type` o mantener texto libre durante MVP.

### Drivers

- Definir flujo formal para vincular/desvincular `drivers.user_id`.
- Definir si un driver suspendido puede ser reactivado por `company_admin`.
- Evaluar filtros por `status`, documento, licencia y email.
- Validar formato de telefono/email en application cuando el MVP lo requiera.

### Routes

- Definir si una ruta archivada puede reactivarse mediante flujo explicito.
- Agregar filtros por ciudad origen, ciudad destino y status.
- Revisar si `base_price` necesita manejo decimal estricto en Go antes de incorporar pagos/tarifas avanzadas.
- Definir reglas de negocio para rutas con mismo origen/destino pero distinto servicio.

### Route Stops

- Exponer endpoint de reorder si el admin-web lo necesita.
- Definir si al eliminar un stop se debe compactar automaticamente `stop_order`.
- Agregar pruebas de integracion para reorder con base de datos real.
- Evaluar si route stops deben usar PostGIS `geometry(Point, 4326)` en una iteracion posterior.

### Autorizacion y multi-tenancy

- Mantener tests de aislamiento por tenant para cada nuevo modulo.
- Agregar pruebas end-to-end de driver role contra todos los catalogos operativos.
- Revisar si `super_admin` tendra lectura cross-tenant o si queda fuera del MVP.

### API y DX

- Consolidar helpers JSON/error para reducir duplicacion entre handlers.
- Definir formato estandar de errores con `code`, `message` y campos invalidos si el frontend lo necesita.
- Agregar paginacion comun para listados.
- Agregar OpenAPI cuando el contrato de Sprint 2 se estabilice.

### Observabilidad

- Agregar request IDs.
- Agregar logging estructurado de errores sin exponer datos sensibles.
- Agregar metricas basicas de HTTP si el despliegue MVP lo requiere.

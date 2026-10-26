# High-Performance Modular HTTP Server in Go

Servidor HTTP universal desarrollado desde cero en Go, inspirado en la arquitectura modular de Apache 2.4. Diseñado con un enfoque en la seguridad, limpieza de código, alto rendimiento y apto para despliegues en producción mediante contenedores.

---

## Arquitectura y Módulos

El servidor implementa una cadena estricta de middlewares organizados secuencialmente para garantizar la máxima seguridad antes de procesar cualquier petición:

1. **`mod_remoteip`**: Extrae y valida la IP real del cliente, asegurando la compatibilidad con proxies inversos (Cloudflare, Nginx).
2. **`mod_ratelimit`**: Control de tráfico por IP para mitigar ataques de denegación de servicio (DoS) y fuerza bruta.
3. **`mod_bodylimit`**: Restringe estrictamente el tamaño máximo del cuerpo de las peticiones HTTP (`POST`, `PUT`) para evitar saturación de memoria.
4. **`mod_headers`**: Inyecta cabeceras defensivas de seguridad global (HSTS, CSP, X-Frame-Options, X-Content-Type-Options).
5. **`mod_log_config`**: Auditoría detallada de accesos con búfer circular en memoria para métricas en tiempo real.
6. **`mod_authz_core`**: Control de acceso basado en tokens Bearer para proteger rutas privadas.
7. **`mod_api_path`**: Enrutador de API REST con soporte de CORS dinámico basado en listas blancas leídas desde variables de entorno.
8. **`mod_dir`**: Servidor de archivos estáticos optimizado con prevención estricta contra ataques de *Path Traversal*.

---

## 📱 Compatibilidad y Entornos Probados

Gracias a su arquitectura ligera desarrollada desde cero en Go y sin dependencias pesadas, **go-modular-server** garantiza una portabilidad excepcional y un consumo mínimo de recursos. Ha sido validado con éxito en los siguientes entornos:

- **Linux / Ubuntu:** Probado en servidores locales e integrado de forma continua mediante GitHub Actions.
- **Docker:** Contenedorizado mediante imágenes *multi-stage* optimizadas y orquestadas con Docker Compose.
- **Android (Termux):** Probado y verificado en entornos móviles locales a través de Termux, demostrando un rendimiento sobresaliente y una respuesta inmediata en arquitecturas ARM.

### 🚀 Despliegue rápido en Termux (Prueba local)
Si deseas levantar el servidor localmente en un entorno Android con Termux, puedes hacerlo ejecutando los siguientes comandos:

```bash
# 1. Actualizar paquetes e instalar Go y Git
pkg update && pkg install golang git -y

# 2. Clonar el repositorio
git clone [https://github.com/bruno2603588/go-modular-server.git](https://github.com/bruno2603588/go-modular-server.git)
cd go-modular-server

# 3. Ejecutar el servidor
go run main.go

## Estructura del Proyecto

```text
HTTP-SERVER/
├── dist/                # Archivos estáticos del frontend (React / Vite)
├── modules/
│   ├── mod_api_path.go  # Enrutador y controlador de API y CORS
│   ├── mod_authz_core.go# Autenticación y autorización por tokens
│   ├── mod_bodylimit.go # Límite de tamaño para el cuerpo de peticiones
│   ├── mod_dir.go       # Manejador seguro de archivos estáticos
│   ├── mod_headers.go   # Inyección de cabeceras de seguridad
│   ├── mod_log_config.go# Auditoría y almacenamiento de logs en memoria
│   ├── mod_ratelimit.go # Limitador de tasa de peticiones por IP
│   ├── mod_remoteip.go  # Extracción de IP real del cliente
│   └── mod_test.go      # Pruebas unitarias automatizadas
├── .env                 # Variables de entorno locales (Privado)
├── Dockerfile           # Imagen multi-stage optimizada para producción
├── docker-compose.yml   # Orquestador de contenedores
├── go.mod               # Definición de módulos de Go
├── main.go              # Punto de entrada y ensamblaje de la cadena
└── README.md            # Documentación del proyecto

# Feed de Posts con Go, Rust y Python

## 🚀 Visión General del Proyecto
Este es un proyecto de un feed de posts que demuestra una arquitectura de microservicios robusta y escalable. La aplicación está diseñada para manejar el tráfico de manera eficiente al distribuir las responsabilidades entre tres servicios principales, cada uno optimizado para una tarea específica.

* **Go** maneja la API principal y el tráfico de la web.
* **Rust** gestiona las notificaciones en tiempo real a través de WebSockets.
* **Python** se encarga de las tareas de procesamiento de datos en segundo plano.

## 🛠️ Tecnologías Utilizadas

### Backend - API Principal (Go)
El servicio principal está desarrollado en **Go** con el framework **Gin**. Su rol es manejar las peticiones HTTP y la lógica de la API REST, sirviendo como el "cerebro" central de la aplicación.
* **Go**: Para la API REST de alto rendimiento.
* **Gin**: Framework web para la gestión de rutas y middleware.
* **SQLite3**: Base de datos para el almacenamiento persistente de usuarios y publicaciones.
* **gRPC**: Para la comunicación eficiente con el servicio de Python.
* **CORS**: Para permitir peticiones desde el frontend.

### Backend - Servicio de Notificaciones (Rust)
El servicio de notificaciones está construido en **Rust** con el framework **Actix**. Su propósito es manejar las conexiones en tiempo real de los usuarios a través de WebSockets, permitiendo enviar notificaciones instantáneas (por ejemplo, cuando un usuario publica un nuevo post) sin afectar el rendimiento de la API principal.
* **Rust**: Para un rendimiento y concurrencia sin igual.
* **Actix**: Framework para construir servicios de WebSockets y actores.
* **WebSockets**: Para la comunicación en tiempo real con el frontend.

### Backend - Procesamiento de Datos (Python)
Este servicio, desarrollado en **Python**, se encarga de las tareas de procesamiento que no requieren respuestas en tiempo real. Por ejemplo, podría procesar los datos de las publicaciones para un análisis posterior, moderar contenido o realizar tareas programadas. La comunicación con este servicio se hace a través de **gRPC**.
* **Python**: Para tareas de procesamiento y análisis de datos.
* **gRPC**: Para la comunicación de alto rendimiento entre servicios.

### Frontend
* **JavaScript (Vanilla)**: Para la lógica interactiva del lado del cliente.
* **HTML5 & CSS3**: Estructura y estilos de la interfaz de usuario.

## 📂 Estructura del Proyecto

mi-red-social/
├── api-go/                # Servicio de la API principal
│   ├── main.go
│   ├── ...
├── rust-notifications/    # Servicio de notificaciones en tiempo real
│   ├── src/
│   │   └── main.rs
│   ├── Cargo.toml
│   └── ...
├── python-processor/      # Servicio de procesamiento de datos
│   ├── main.py
│   ├── requirements.txt
│   └── ...
├── frontend/              # Frontend (HTML, CSS, JS)
│   ├── index.html
│   ├── ...
└── README.md

## ⚙️ Cómo Ejecutar el Proyecto

1.  **Clonar el repositorio.**

2.  **Ejecutar el servicio de la API de Go:**
    `cd api-go && go run .`

3.  **Ejecutar el servicio de Rust:**
    `cd rust-notifications && cargo run`

4.  **Ejecutar el servicio de Python:**
    `cd python-processor && python main.py`

5.  **Ejecutar el Frontend:**
    Sirve la carpeta `frontend` con un servidor de archivos estáticos. Puedes usar el servidor de Go simple o cualquier otro.

6.  **Abrir la aplicación:**
    Ve a tu navegador y visita la URL de tu frontend (por lo general, `http://localhost:8000`).

🚀 Seeder (Poblador de Datos)
Para facilitar el desarrollo y las pruebas, este proyecto incluye un script de seeder que genera datos de prueba de forma automática.

El seeder está ubicado en el archivo seeder.go en la raíz del proyecto. Este script se encarga de:

Crear 100 usuarios con nombres y correos electrónicos generados aleatoriamente.

Crear 100 publicaciones, asignando una publicación a cada uno de los usuarios creados.

Cómo Usar el Seeder
El seeder se ejecuta automáticamente cada vez que inicias el servidor principal de Go. No necesitas hacer nada más que seguir los pasos de la sección "Cómo Ejecutar el Proyecto".

Nota: Si ya tienes datos en tu base de datos, el seeder los eliminará y los reemplazará con los nuevos datos de prueba.

## 🤝 Contribuciones
¡Las contribuciones son bienvenidas!

## 📄 Licencia
Este proyecto está bajo la Licencia MIT.
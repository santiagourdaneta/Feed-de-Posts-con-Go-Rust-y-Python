package main

import (
    "net/http" // Para manejar las respuestas del servidor
    "github.com/gin-gonic/gin" // El framework Gin
    "database/sql"
    _ "github.com/mattn/go-sqlite3" // La biblioteca para hablar con SQLite
    "github.com/didip/tollbooth/v7"
    "github.com/didip/tollbooth_gin"
     "log"
    "golang.org/x/crypto/bcrypt"
    "strconv"
    "github.com/gin-contrib/cors" 
)

// Este es el molde para las galletas (las publicaciones)
type Post struct {
    ID        int    `json:"id"`
    UserID    int    `json:"user_id"`
    Content   string `json:"content"`
    CreatedAt string `json:"created_at"`
    Username  string `json:"username"` 
}

// Una variable para que todo tu programa sepa dónde está el cuaderno de notas
var db *sql.DB

// Una función para crear las mesas en el cuaderno de notas
func createTables() {
    // La orden para crear la mesa de usuarios
    sqlStmt := `
    create table if not exists users (id integer not null primary key autoincrement, username text, email text, password_hash text);
    `
    _, err := db.Exec(sqlStmt)
    if err != nil {
        panic(err)
    }

    // La orden para crear la mesa de publicaciones
    sqlStmt = `
    create table if not exists posts (id integer not null primary key autoincrement, user_id integer, content text, created_at text);
    `
    _, err = db.Exec(sqlStmt)
    if err != nil {
        panic(err)
    }
}

func main() {

     // Aquí es donde abrimos el cuaderno de notas, si no existe lo crea
    var err error
    db, err = sql.Open("sqlite3", "./mi_red_social.db")
    if err != nil {
        // Si no se puede abrir el cuaderno, el programa se detiene y nos avisa
        panic("Falla al abrir la base de datos: " + err.Error())
    }
    defer db.Close() // Esto asegura que el cuaderno se cierre cuando el programa termine

    // Crea las "mesas" si aún no existen
   // createTables()
    // Ejecuta la función del seeder para poblar la base de datos
   // SeedDatabase(db)

   // Index for user_id in the posts table
// This makes the JOIN on posts.user_id extremely fast.
_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);`)
if err != nil {
    log.Fatal("Failed to create user_id index:", err)
}

// Index for created_at in the posts table
// This speeds up the ORDER BY clause for getting the latest posts.
_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);`)
if err != nil {
    log.Fatal("Failed to create created_at index:", err)
}


    // 1. Crea el motor de nuestro servidor, como si fuera el motor de un coche.
    router := gin.Default()

    // Add CORS middleware to allow requests from your frontend server
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://127.0.0.1:8000", "http://localhost:8000"},
        AllowMethods:     []string{"GET", "POST"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    //Contrata al guardia y dile las reglas (10 peticiones por minuto)
    lmt := tollbooth.NewLimiter(10, nil)
    lmt.SetMessage("Demasiadas peticiones. Inténtalo de nuevo más tarde.")
    //Pon al guardia a revisar todas las peticiones
    router.Use(tollbooth_gin.LimitHandler(lmt))


    // 2. Dibuja los caminos (las rutas de la API)

    router.POST("/api/users/register", func(c *gin.Context) {
        // 1. Aquí está nuestro molde para el paquete
        var user struct {
            Username string `json:"username" binding:"required"`
            Email    string `json:"email" binding:"required,email"`
            Password string `json:"password" binding:"required"`
        }

        // 2. El inspector revisa el paquete
        if err := c.ShouldBindJSON(&user); err != nil {
            // Si el paquete está mal, le decimos "¡Paquete rechazado!"
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // 3. Si todo está bien, podemos guardar al usuario
// 3.1. Convertimos la contraseña a un código secreto y seguro (hash)
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo procesar la contraseña."})
    return
}

// 3.2. Le decimos a Go que guarde los datos en la mesa de usuarios
// Usamos "db.Exec" para ejecutar una orden en el cuaderno de notas.
_, err = db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", user.Username, user.Email, string(hashedPassword))
if err != nil {
    // Si ya existe un usuario con ese nombre o correo, avisa del error
    // En un proyecto real, manejarías errores específicos de SQLite.
    c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar al usuario. Es posible que el nombre de usuario o correo ya estén en uso."})
    return
}

c.JSON(http.StatusCreated, gin.H{"message": "¡Usuario registrado con éxito!"})
    })

    router.GET("/api/users/:username", func(c *gin.Context) {
        // 1. Le decimos a Go que tome el nombre del usuario de la dirección
        username := c.Param("username")

        // 2. Le preguntamos al cuaderno de notas si tiene ese nombre
        var id int
        var email string
        var passwordHash string
        
        // Aquí está la magia: la orden para buscar en el cuaderno
        row := db.QueryRow("SELECT id, email, password_hash FROM users WHERE username = ?", username)
        err := row.Scan(&id, &email, &passwordHash)

        if err != nil {
            // 3. Si no lo encuentra (porque la fila está vacía), le dices que el perfil no existe
            if err == sql.ErrNoRows {
                c.JSON(http.StatusNotFound, gin.H{"error": "El perfil no fue encontrado."})
                return
            }
            // Si hay otro problema, le dices que algo salió mal
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor."})
            return
        }

        // 4. Si lo encuentra, tomas los datos y se los devuelves
        c.JSON(http.StatusOK, gin.H{
            "id":       id,
            "username": username,
            "email":    email,
        })
    })

    // Un camino para crear una publicación nueva.
    // Usamos POST porque queremos ENVIAR nueva información.
    router.POST("/api/posts/create", func(c *gin.Context) {
    // 1. Le pides a Go que tome el texto de la publicación y el ID del usuario
    var newPost struct {
        UserID  int    `json:"user_id"`
        Content string `json:"content"`
    }
    if err := c.ShouldBindJSON(&newPost); err != nil {
        // Si lo que nos envían no tiene los datos correctos, les dices que está mal
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos."})
        return
    }

    // 2. Le ordenas a Go que guarde la información en la mesa de publicaciones
    // El "INSERT INTO" es la orden para guardar
    result, err := db.Exec("INSERT INTO posts (user_id, content, created_at) VALUES (?, ?, datetime('now'))", newPost.UserID, newPost.Content)
    if err != nil {
        // Si no se puede guardar, le avisas a la persona que hubo un problema
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la publicación."})
        return
    }

    // 3. Si todo salió bien, le dices que la publicación fue un éxito
    id, _ := result.LastInsertId()
    c.JSON(http.StatusCreated, gin.H{
        "message":      "Publicación creada con éxito.",
        "post_id":      id,
        "user_id":      newPost.UserID,
        "content":      newPost.Content,
    })
})

    // Un camino para ver el muro de publicaciones.
    // Usamos GET para OBTENER las publicaciones.
   router.GET("/api/posts/feed", func(c *gin.Context) {
    // Leer los parámetros de la URL
    page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
    if err != nil || page < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
        return
    }
    limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
    if err != nil || limit < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
        return
    }

    // Calcula el desplazamiento (offset)
    offset := (page - 1) * limit

    // Usa el LIMIT y OFFSET en tu consulta
    rows, err := db.Query(`
        SELECT p.id, p.user_id, p.content, p.created_at, u.username
        FROM posts p
        JOIN users u ON p.user_id = u.id
        ORDER BY p.created_at DESC
        LIMIT ? OFFSET ?
    `, limit, offset)
    if err != nil {
        log.Println("Error fetching posts:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var p Post
        if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.CreatedAt, &p.Username); err != nil {
            log.Println("Error scanning post row:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
            return
        }
        posts = append(posts, p)
    }

    c.JSON(http.StatusOK, posts)
})

    // 3. ¡Enciende el motor y haz que el servidor escuche!
    // Esto hará que tu red social empiece a funcionar en el puerto 8080.
   router.Run(":8082")
}
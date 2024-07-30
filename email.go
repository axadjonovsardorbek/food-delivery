package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "gopkg.in/gomail.v2"
    "github.com/dgrijalva/jwt-go"
    _ "github.com/lib/pq"
    "context"
)

// Email xabar strukturasi
type EmailMessage struct {
    ID      int    `json:"id"`
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

var (
    redisClient *redis.Client
    db          *sql.DB
    ctx         = context.Background()
)

func main() {
    var err error

    // Redis bilan bog'lanish
    redisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // PostgreSQL bilan bog'lanish
    db, err = sql.Open("postgres", "user=postgres dbname=emailservice sslmode=disable")
    if err != nil {
        log.Fatalf("PostgreSQL ulanishda xato: %s", err)
    }
    defer db.Close()

    // Gin router
    r := gin.Default()

    // JWT middleware
    r.Use(jwtMiddleware())

    // Email yuborish uchun endpoint
    r.POST("/send-email", sendEmail)

    // Redis-dan xabarlarni qabul qilish uchun goroutine
    go consumeEmails()

    // Server ishga tushirish
    r.Run(":8050")
}

func jwtMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte("your-secret-key"), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        c.Next()
    }
}

func sendEmail(c *gin.Context) {
    var email EmailMessage
    if err := c.ShouldBindJSON(&email); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Email ma'lumotlarini PostgreSQL-ga saqlash
    _, err := db.Exec("INSERT INTO emails (to_address, subject, body) VALUES ($1, $2, $3)",
        email.To, email.Subject, email.Body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Email saqlashda xato"})
        return
    }

    // Email-ni Redis-ga yuborish
    emailJSON, _ := json.Marshal(email)
    err = redisClient.RPush(ctx, "email_queue", emailJSON).Err()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis-ga yuborishda xato"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Email qabul qilindi va navbatga qo'yildi"})
}

func consumeEmails() {
    for {
        // Redis-dan xabarni olish
        result, err := redisClient.BLPop(ctx, 0, "email_queue").Result()
        if err != nil {
            log.Printf("Redis-dan xabar olishda xato: %v", err)
            continue
        }

        var email EmailMessage
        err = json.Unmarshal([]byte(result[1]), &email)
        if err != nil {
            log.Printf("JSON parse qilishda xato: %v", err)
            continue
        }

        sendEmailToSMTP(email)
    }
}

func sendEmailToSMTP(email EmailMessage) {
    m := gomail.NewMessage()
    m.SetHeader("From", "from@example.com")
    m.SetHeader("To", email.To)
    m.SetHeader("Subject", email.Subject)
    m.SetBody("text/plain", email.Body)

    d := gomail.NewDialer("smtp.example.com", 587, "user", "password")

    if err := d.DialAndSend(m); err != nil {
        log.Printf("Email yuborishda xato: %v", err)
    } else {
        log.Printf("Email muvaffaqiyatli yuborildi: %s", email.To)
    }
}
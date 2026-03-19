package intializer

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Id       int64  `gorm:"primaryKey"`
	UserName string `gorm:"size:255"`
	Email    string `gorm:"size:255"`
}
type Auth struct {
	Id       uint64 `gorm:"primarykey"`
	Name     string `json:"name"`
	Gmail    string `json:"email"`
	Password string `json:"password"`
}
type ContactMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
type Event struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Status      string `json:"status" gorm:"default:pending"`
}

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("database connection failed ", err)
	}
	DB = db

	log.Println("database created and runed sucessfully")
}
func Checkemail(email string) bool {

	var auth Auth
	check := DB.Where("Gmail=?", email).First(&auth)
	return check.Error == nil
}

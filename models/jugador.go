package models

type Jugador struct {
    ID       uint   `gorm:"primaryKey;column:id" json:"id"`
    Name     string `gorm:"column:Name" json:"name"`
    Email    string `gorm:"column:Email" json:"email"`
    Password string `gorm:"column:Password" json:"-"`
}
func (Jugador) TableName() string {
	return "jugador"
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}


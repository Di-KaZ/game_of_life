package database

type User struct {
	ID       uint `gorm:"primarykey"`
	Name     string
	Password string
	Width    int
	Height   int
	Alive    int
}

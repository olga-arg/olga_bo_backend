package domain

type UserTeam struct {
	UserID string `gorm:"primaryKey;foreignKey:UserID"`
	TeamID string `gorm:"primaryKey;foreignKey:TeamID"`
}

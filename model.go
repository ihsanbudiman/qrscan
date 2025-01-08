package main

import (
	"time"

	"github.com/guregu/null/v5"
)

type User struct {
	Id        uint      `json:"id" gorm:"column:id"`
	Username  string    `json:"username" gorm:"column:username"`
	Password  string    `json:"password" gorm:"column:password"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt null.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

type Qrscan struct {
	Id         uint      `json:"id" gorm:"column:id"`
	Uuid       string    `json:"uuid" gorm:"column:uuid"`
	UserId     null.Int  `json:"user_id" gorm:"column:user_id"`
	IsValid    bool      `json:"is_valid" gorm:"column:is_valid"`
	ValidUntil string    `json:"valid_until" gorm:"column:valid_until"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt  null.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

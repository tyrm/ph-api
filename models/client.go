package models

import (
	"fmt"
	"time"
)

type Client struct {
	ID     uint          `json:"id" gorm:"primary_key"`
	Name   string        `json:"name"`

	ClientID string      `json:"access" gorm:"not null;unique"`
	Secret   string      `json:"secret"`
	Domain   string      `json:"domain"`
	User     User        `json:"user"`
	UserID   int         `json:"-"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

// GetID client id
func (c *Client) GetID() string {
	return c.ClientID
}

// GetSecret client domain
func (c *Client) GetSecret() string {
	return c.Secret
}

// GetDomain client domain
func (c *Client) GetDomain() string {
	return c.Domain
}

// GetUserID user id
func (c *Client) GetUserID() string {
	return fmt.Sprint(c.User.ID)
}

func GetClient(id string) (cli Client, err error) {
	err = db.Where("client_id=?", id).First(&cli).Error
	if err != nil {
		logger.Errorf("GetClient(%s) Error: %s", id, err)
	}

	return
}

func SetClient(cli *Client) (err error) {
	err = db.Create(&cli).Error
	if err != nil {
		logger.Errorf("SetClient(%s) Error: %s", cli.ClientID, err)
	}

	return
}
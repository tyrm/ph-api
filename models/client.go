package models

import (
	"fmt"
	"time"

	"github.com/google/jsonapi"
)

type Client struct {
	ID     uint          `gorm:"primary_key"`
	Name   string        `jsonapi:"attr,name"`

	ClientID string      `jsonapi:"primary,client" gorm:"not null;unique"`
	Secret   string      `jsonapi:"attr,secret"`
	Domain   string      `jsonapi:"attr,domain"`
	User     *User        `jsonapi:"relation,user"`
	UserID   int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
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

func GetClientPage(count int, page int) (clients []Client, err error) {
	offset := count * page;
	err = db.Preload("User").Limit(count).Offset(offset).Find(&clients).Error

	return
}

func (c *Client) JSONAPIMeta() *jsonapi.Meta {
	return &jsonapi.Meta{
		"created_at": c.CreatedAt,
		"updated_at": c.UpdatedAt,
	}
}

func SetClient(cli *Client) (err error) {
	err = db.Create(&cli).Error
	if err != nil {
		logger.Errorf("SetClient(%s) Error: %s", cli.ClientID, err)
	}

	return
}
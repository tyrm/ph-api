package models

import "time"

type Client struct {
	ID     string `gorm:"primary_key",json:"id"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
	Domain string `json:"domain"`
	UserID string `json:"user_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// GetID client id
func (c *Client) GetID() string {
	return c.ID
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
	return c.UserID
}

func GetClient(id string) (cli Client, err error) {
	err = db.Where("id=?", id).First(&cli).Error
	if err != nil {
		logger.Errorf("GetClient(%s) Error: %s", id, err)
	}

	return
}

func SetClient(cli *Client) (err error) {
	err = db.Create(&cli).Error
	if err != nil {
		logger.Errorf("SetClient(%s) Error: %s", cli.ID, err)
	}

	return
}
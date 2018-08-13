package models

import "github.com/jinzhu/gorm"

type Client struct {
	gorm.Model

	CID    string `json:"id"`
	Secret string `json:"decret"`
	Domain string `json:"domain"`
	UserID string `json:"user_id"`
}

// GetID client id
func (c *Client) GetID() string {
	return c.CID
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
	err = db.Where("c_id=?", id).First(&cli).Error
	if err != nil {
		logger.Errorf("GetClient(%s) Error: %s", id, err)
	}

	return
}

func SetClient(id string, secret string, domain string, userID string) (cli Client, err error) {
	newClient := Client{
		CID:    id,
		Secret: secret,
		Domain: domain,
		UserID: userID,
	}

	err = db.Create(&newClient).Error
	if err != nil {
		logger.Errorf("SetClient(%s) Error: %s", newClient.CID, err)
	}

	return
}
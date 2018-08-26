package oauth

import (
	"strconv"

	"../models"

	"gopkg.in/oauth2.v3"
)

func NewClientStore() *ClientStore {
	return &ClientStore{}
}

// ClientStore client information store
type ClientStore struct {
}

// GetByID according to the CID for the client information
func (cs *ClientStore) GetByID(id string) (cli oauth2.ClientInfo, err error) {
	client, err := models.GetClient(id)
	if err != nil {
		logger.Errorf("Error getting client: %s", err)
		return
	}

	cli = &client
	return
}

// Set set client information
func (cs *ClientStore) Set(id string, cli oauth2.ClientInfo) (err error) {
	i, err := strconv.Atoi(cli.GetUserID())
	if err != nil {
		return
	}

	usr, err := models.GetUser(i)
	if err != nil {
		return
	}
	err = models.SetClient(&models.Client{
		ClientID:    id,
		Secret: cli.GetSecret(),
		Domain: cli.GetDomain(),
		User: &usr,
	})

	return
}

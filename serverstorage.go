package gsclient

import (
	"net/http"
	"path"
)

//ServerStorageRelationList JSON struct of a list of relations between a server and storages
type ServerStorageRelationList struct {
	List []ServerStorageRelationProperties `json:"storage_relations"`
}

//ServerStorageRelationSingle JSON struct of a single relation between a server and a storage
type ServerStorageRelationSingle struct {
	Properties ServerStorageRelationProperties `json:"storage_relation"`
}

//ServerStorageRelationProperties JSON struct of properties of a relation between a server and a storage
type ServerStorageRelationProperties struct {
	ObjectUuid       string `json:"object_uuid"`
	ObjectName       string `json:"object_name"`
	Capacity         int    `json:"capacity"`
	StorageType      string `json:"storage_type"`
	Target           int    `json:"target"`
	Lun              int    `json:"lun"`
	Controller       int    `json:"controller"`
	CreateTime       string `json:"create_time"`
	BootDevice       bool   `json:"bootdevice"`
	Bus              int    `json:"bus"`
	LastUsedTemplate string `json:"last_used_template"`
	LicenseProductNo int    `json:"license_product_no"`
	ServerUuid       string `json:"server_uuid"`
}

//ServerStorageRelationCreateRequest JSON struct of a request for creating a relation between a server and a storage
type ServerStorageRelationCreateRequest struct {
	ObjectUuid string `json:"object_uuid"`
	BootDevice bool   `json:"bootdevice,omitempty"`
}

//ServerStorageRelationUpdateRequest JSON struct of a request for updating a relation between a server and a storage
type ServerStorageRelationUpdateRequest struct {
	Ordering   int      `json:"ordering,omitempty"`
	BootDevice bool     `json:"bootdevice,omitempty"`
	L3security []string `json:"l3security,omitempty"`
}

//GetServerStorageList gets a list of a specific server's storages
func (c *Client) GetServerStorageList(id string) ([]ServerStorageRelationProperties, error) {
	r := Request{
		uri:    path.Join(apiServerBase, id, "storages"),
		method: http.MethodGet,
	}
	var response ServerStorageRelationList
	err := r.execute(*c, &response)
	return response.List, err
}

//GetServerStorage gets a storage of a specific server
func (c *Client) GetServerStorage(serverId, storageId string) (ServerStorageRelationProperties, error) {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "storages", storageId),
		method: http.MethodGet,
	}
	var response ServerStorageRelationSingle
	err := r.execute(*c, &response)
	return response.Properties, err
}

//UpdateServerStorage updates a link between a storage and a server
func (c *Client) UpdateServerStorage(serverId, storageId string, body ServerStorageRelationUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "storages", storageId),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//CreateServerStorage create a link between a server and a storage
func (c *Client) CreateServerStorage(id string, body ServerStorageRelationCreateRequest) error {
	r := Request{
		uri:    path.Join(apiServerBase, id, "storages"),
		method: http.MethodPost,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteServerStorage delete a link between a storage and a server
func (c *Client) DeleteServerStorage(serverId, storageId string) error {
	r := Request{
		uri:    path.Join(apiServerBase, serverId, "storages", storageId),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//LinkStorage attaches a storage to a server
func (c *Client) LinkStorage(serverId string, storageId string, bootdevice bool) error {
	body := ServerStorageRelationCreateRequest{
		ObjectUuid: storageId,
		BootDevice: bootdevice,
	}
	return c.CreateServerStorage(serverId, body)
}

//UnlinkStorage remove a storage from a server
func (c *Client) UnlinkStorage(serverId string, storageId string) error {
	return c.DeleteServerStorage(serverId, storageId)
}

package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//ServerNetworkRelationList JSON struct of a list of relations between a server and networks
type ServerNetworkRelationList struct {
	//Array of relations between a server and networks
	List []ServerNetworkRelationProperties `json:"network_relations"`
}

//ServerNetworkRelation JSON struct of a single relation between a server and a network
type ServerNetworkRelation struct {
	//Properties of a relation between a server and a network
	Properties ServerNetworkRelationProperties `json:"network_relation"`
}

//ServerNetworkRelationProperties JSON struct of properties of a relation between a server and a network
type ServerNetworkRelationProperties struct {
	//Defines information about MAC spoofing protection (filters layer2 and ARP traffic based on MAC source).
	//It can only be (de-)activated on a private network - the public network always has l2security enabled.
	//It will be true if the network is public, and false if the network is private.
	L2security bool `json:"l2security"`

	//The UUID of the Server.
	ServerUUID string `json:"server_uuid"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//True if the network is public. If private it will be false.
	//Each private network is a secure and fully transparent 2-Layer network between servers.
	//There is no limit on how many servers can be connected to the same private network.
	PublicNet bool `json:"public_net"`

	//The UUID of firewall template.
	FirewallTemplateUUID string `json:"firewall_template_uuid"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	ObjectName string `json:"object_name"`

	//network_mac defines the MAC address of the network interface.
	Mac string `json:"mac"`

	//Defines if this object is the bootdevice. Storages, Networks and ISO-Images can have a bootdevice configured,
	//but only one bootdevice per Storage, Network or ISO-Image.
	//The boot order is as follows => Network > ISO-Image > Storage.
	BootDevice bool `json:"bootdevice"`

	//PartnerUUID
	PartnerUUID string `json:"partner_uuid"`

	//Defines the ordering of the network interfaces. Lower numbers have lower PCI-IDs.
	Ordering int `json:"ordering"`

	//Firewall that is used to this server network relation
	Firewall FirewallRules `json:"firewall"`

	//(one of network, network_high, network_insane)
	NetworkType string `json:"network_type"`

	//The UUID of the network you're requesting.
	NetworkUUID string `json:"network_uuid"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//Defines information about IP prefix spoof protection (it allows source traffic only from the IPv4/IPv4 network prefixes).
	//If empty, it allow no IPv4/IPv6 source traffic. If set to null, l3security is disabled (default).
	L3security []string `json:"l3security"`
}

//ServerNetworkRelationCreateRequest JSON struct of a request for creating a relation between a server and a network
type ServerNetworkRelationCreateRequest struct {
	//The uuid of network you wish to add. Only 7 private networks are allowed to be attached to a server
	ObjectUUID string `json:"object_uuid"`

	//The ordering of the network interfaces. Lower numbers have lower PCI-IDs. Can be empty.
	Ordering int `json:"ordering,omitempty"`

	//Whether the server boots from this network or not. Can be empty.
	BootDevice bool `json:"bootdevice,omitempty"`

	//Defines information about IP prefix spoof protection (it allows source traffic only from the IPv4/IPv4 network prefixes).
	//If empty, it allow no IPv4/IPv6 source traffic. If set to null, l3security is disabled (default).
	//Can be empty
	L3security []string `json:"l3security,omitempty"`

	//All rules of Firewall. Can be empty
	Firewall *FirewallRules `json:"firewall,omitempty"`

	//Instead of setting firewall rules manually, you can use a firewall template by setting UUID of the firewall template.
	//Can be empty.
	FirewallTemplateUUID string `json:"firewall_template_uuid,omitempty"`
}

//ServerNetworkRelationUpdateRequest JSON struct of a request for updating a relation between a server and a network
type ServerNetworkRelationUpdateRequest struct {
	//The ordering of the network interfaces. Lower numbers have lower PCI-IDs. Optional.
	Ordering int `json:"ordering,omitempty"`

	//The ordering of the network interfaces. Lower numbers have lower PCI-IDs. Optional.
	BootDevice bool `json:"bootdevice,omitempty"`

	//Defines information about IP prefix spoof protection (it allows source traffic only from the IPv4/IPv4 network prefixes).
	//If empty, it allow no IPv4/IPv6 source traffic. If set to null, l3security is disabled (default).
	//Can be empty
	L3security []string `json:"l3security,omitempty"`

	//All rules of Firewall. Optional.
	Firewall *FirewallRules `json:"firewall,omitempty"`

	//Instead of setting firewall rules manually, you can use a firewall template by setting UUID of the firewall template.
	//Optional.
	FirewallTemplateUUID string `json:"firewall_template_uuid,omitempty"`
}

//GetServerNetworkList gets a list of a specific server's networks
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getServerLinkedNetworks
func (c *Client) GetServerNetworkList(id string) ([]ServerNetworkRelationProperties, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "networks"),
		method: http.MethodGet,
	}
	var response ServerNetworkRelationList
	err := r.execute(*c, &response)
	return response.List, err
}

//GetServerNetwork gets a network of a specific server
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getServerLinkedNetwork
func (c *Client) GetServerNetwork(serverID, networkID string) (ServerNetworkRelationProperties, error) {
	if !isValidUUID(serverID) || !isValidUUID(networkID) {
		return ServerNetworkRelationProperties{}, errors.New("'serverID' or 'networksID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "networks", networkID),
		method: http.MethodGet,
	}
	var response ServerNetworkRelation
	err := r.execute(*c, &response)
	return response.Properties, err
}

//UpdateServerNetwork updates a link between a network and a server
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateServerLinkedNetwork
func (c *Client) UpdateServerNetwork(serverID, networkID string, body ServerNetworkRelationUpdateRequest) error {
	if !isValidUUID(serverID) || !isValidUUID(networkID) {
		return errors.New("'serverID' or 'networksID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "networks", networkID),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//CreateServerNetwork creates a link between a network and a storage
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/linkNetworkToServer
func (c *Client) CreateServerNetwork(id string, body ServerNetworkRelationCreateRequest) error {
	if !isValidUUID(id) || !isValidUUID(body.ObjectUUID) {
		return errors.New("'serverID' or 'network_id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "networks"),
		method: http.MethodPost,
		body:   body,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		return c.waitForServerNetworkRelCreation(id, body.ObjectUUID)
	}
	return r.execute(*c, nil)
}

//DeleteServerNetwork deletes a link between a network and a server
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/unlinkNetworkFromServer
func (c *Client) DeleteServerNetwork(serverID, networkID string) error {
	if !isValidUUID(serverID) || !isValidUUID(networkID) {
		return errors.New("'serverID' or 'networkID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "networks", networkID),
		method: http.MethodDelete,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		return c.waitForServerNetworkRelDeleted(serverID, networkID)
	}
	return r.execute(*c, nil)
}

//LinkNetwork attaches a network to a server
func (c *Client) LinkNetwork(serverID, networkID, firewallTemplate string, bootdevice bool, order int,
	l3security []string, firewall *FirewallRules) error {
	body := ServerNetworkRelationCreateRequest{
		ObjectUUID:           networkID,
		Ordering:             order,
		BootDevice:           bootdevice,
		L3security:           l3security,
		FirewallTemplateUUID: firewallTemplate,
		Firewall:             firewall,
	}
	return c.CreateServerNetwork(serverID, body)
}

//UnlinkNetwork removes the link between a network and a server
func (c *Client) UnlinkNetwork(serverID string, networkID string) error {
	return c.DeleteServerNetwork(serverID, networkID)
}

//waitForServerNetworkRelCreation allows to wait until the relation between a server and a network is created
func (c *Client) waitForServerNetworkRelCreation(serverID, networkID string) error {
	if !isValidUUID(serverID) || !isValidUUID(networkID) {
		return errors.New("'serverID' and 'networkID' are required")
	}
	uri := path.Join(apiServerBase, serverID, "networks", networkID)
	method := http.MethodGet
	return c.waitFor200Status(uri, method)
}

//waitForServerNetworkRelDeleted allows to wait until the relation between a server and a network is deleted
func (c *Client) waitForServerNetworkRelDeleted(serverID, networkID string) error {
	if !isValidUUID(serverID) || !isValidUUID(networkID) {
		return errors.New("'serverID' and 'networkID' are required")
	}
	uri := path.Join(apiServerBase, serverID, "networks", networkID)
	method := http.MethodGet
	return c.waitFor404Status(uri, method)
}

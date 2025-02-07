package gsclient

import (
	"errors"
	"fmt"
	"net/http"
	"path"
)

//NetworkList is JSON struct of a list of networks
type NetworkList struct {
	//Array of networks
	List map[string]NetworkProperties `json:"networks"`
}

//DeletedNetworkList is JSON struct of a list of deleted networks
type DeletedNetworkList struct {
	//Array of deleted networks
	List map[string]NetworkProperties `json:"deleted_networks"`
}

//Network is JSON struct of a single network
type Network struct {
	//Properties of a network
	Properties NetworkProperties `json:"network"`
}

//NetworkProperties is JSON struct of a network's properties
type NetworkProperties struct {
	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationCountry string `json:"location_country"`

	//Helps to identify which datacenter an object belongs to.
	LocationUUID string `json:"location_uuid"`

	//True if the network is public. If private it will be false.
	PublicNet bool `json:"public_net"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//One of 'network', 'network_high' or 'network_insane'.
	NetworkType string `json:"network_type"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//Status indicates the status of the object.
	Status string `json:"status"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//Defines information about MAC spoofing protection (filters layer2 and ARP traffic based on MAC source).
	//It can only be (de-)activated on a private network - the public network always has l2security enabled.
	//It will be true if the network is public, and false if the network is private.
	L2Security bool `json:"l2security"`

	//Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	//Uses IATA airport code, which works as a location identifier.
	LocationIata string `json:"location_iata"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationName string `json:"location_name"`

	//Defines if the object is administratively blocked. If true, it can not be deleted by the user.
	DeleteBlock bool `json:"delete_block"`

	//List of labels.
	Labels []string `json:"labels"`

	//The information about other object which are related to this network. the object could be servers and/or vlans
	Relations NetworkRelations `json:"relations"`
}

//NetworkRelations is JSON struct of a list of a network's relations
type NetworkRelations struct {
	//Array of object (NetworkVlan)
	Vlans []NetworkVlan `json:"vlans"`

	//Array of object (NetworkServer)
	Servers []NetworkServer `json:"servers"`
}

//NetworkVlan is JSON struct of a relation between a network and a VLAN
type NetworkVlan struct {
	//Vlan
	Vlan int `json:"vlan"`

	//Name of tenant
	TenantName string `json:"tenant_name"`

	//UUID of tenant
	TenantUUID string `json:"tenant_uuid"`
}

//NetworkServer is JSON struct of a relation between a network and a server
type NetworkServer struct {
	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//Network_mac defines the MAC address of the network interface.
	Mac string `json:"mac"`

	//Whether the server boots from this iso image or not.
	Bootdevice bool `json:"bootdevice"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//Defines information about IP prefix spoof protection (it allows source traffic only from the IPv4/IPv4 network prefixes).
	//If empty, it allow no IPv4/IPv6 source traffic. If set to null, l3security is disabled (default).
	L3security []string `json:"l3security"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	ObjectName string `json:"object_name"`

	//The UUID of the network you're requesting.
	NetworkUUID string `json:"network_uuid"`

	//The ordering of the network interfaces. Lower numbers have lower PCI-IDs.
	Ordering int `json:"ordering"`
}

//NetworkCreateRequest is JSON of a request for creating a network
type NetworkCreateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//List of labels. Can be empty.
	Labels []string `json:"labels,omitempty"`

	//Helps to identify which datacenter an object belongs to.
	LocationUUID string `json:"location_uuid"`

	//Defines information about MAC spoofing protection (filters layer2 and ARP traffic based on MAC source).
	//It can only be (de-)activated on a private network - the public network always has l2security enabled.
	//It will be true if the network is public, and false if the network is private.
	L2Security bool `json:"l2security,omitempty"`
}

//NetworkCreateResponse is JSON of a response for creating a network
type NetworkCreateResponse struct {
	//UUID of the network being created
	ObjectUUID string `json:"object_uuid"`

	//UUID of the request
	RequestUUID string `json:"request_uuid"`
}

//NetworkUpdateRequest is JSON of a request for updating a network
type NetworkUpdateRequest struct {
	//New name. Leave it if you do not want to update the name
	Name string `json:"name,omitempty"`

	//L2Security. Leave it if you do not want to update the l2 security
	L2Security bool `json:"l2security"`
}

//GetNetwork get a specific network based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getNetwork
func (c *Client) GetNetwork(id string) (Network, error) {
	if !isValidUUID(id) {
		return Network{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiNetworkBase, id),
		method: http.MethodGet,
	}
	var response Network
	err := r.execute(*c, &response)
	return response, err
}

//CreateNetwork creates a network
//
//See: https://gridscale.io/en//api-documentation/index.html#tag/network
func (c *Client) CreateNetwork(body NetworkCreateRequest) (NetworkCreateResponse, error) {
	r := Request{
		uri:    apiNetworkBase,
		method: http.MethodPost,
		body:   body,
	}
	var response NetworkCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return NetworkCreateResponse{}, err
	}
	if c.cfg.sync {
		err = c.waitForRequestCompleted(response.RequestUUID)
	}
	return response, err
}

//DeleteNetwork deletes a specific network based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteNetwork
func (c *Client) DeleteNetwork(id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiNetworkBase, id),
		method: http.MethodDelete,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForNetworkDeleted(id)
	}
	return r.execute(*c, nil)
}

//UpdateNetwork updates a specific network based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateNetwork
func (c *Client) UpdateNetwork(id string, body NetworkUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiNetworkBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForNetworkActive(id)
	}
	return r.execute(*c, nil)
}

//GetNetworkList gets a list of available networks
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getNetworks
func (c *Client) GetNetworkList() ([]Network, error) {
	r := Request{
		uri:    apiNetworkBase,
		method: http.MethodGet,
	}
	var response NetworkList
	var networks []Network
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		networks = append(networks, Network{
			Properties: properties,
		})
	}
	return networks, err
}

//GetNetworkEventList gets a list of a network's events
//
//See: https://gridscale.io/en//api-documentation/index.html#tag/network
func (c *Client) GetNetworkEventList(id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiNetworkBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var networkEvents []Event
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		networkEvents = append(networkEvents, Event{Properties: properties})
	}
	return networkEvents, err
}

//GetNetworkPublic gets public network
func (c *Client) GetNetworkPublic() (Network, error) {
	networks, err := c.GetNetworkList()
	if err != nil {
		return Network{}, err
	}
	for _, network := range networks {
		if network.Properties.PublicNet {
			return Network{Properties: network.Properties}, nil
		}
	}
	return Network{}, fmt.Errorf("Public Network not found")
}

//GetNetworksByLocation gets a list of networks by location
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedNetworks
func (c *Client) GetNetworksByLocation(id string) ([]Network, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLocationBase, id, "networks"),
		method: http.MethodGet,
	}
	var response NetworkList
	var networks []Network
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		networks = append(networks, Network{Properties: properties})
	}
	return networks, err
}

//GetDeletedNetworks gets a list of deleted networks
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedNetworks
func (c *Client) GetDeletedNetworks() ([]Network, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "networks"),
		method: http.MethodGet,
	}
	var response DeletedNetworkList
	var networks []Network
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		networks = append(networks, Network{Properties: properties})
	}
	return networks, err
}

//waitForNetworkActive allows to wait until the network's status is active
func (c *Client) waitForNetworkActive(id string) error {
	return retryWithTimeout(func() (bool, error) {
		net, err := c.GetNetwork(id)
		return net.Properties.Status != resourceActiveStatus, err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitForNetworkDeleted allows to wait until the network is deleted
func (c *Client) waitForNetworkDeleted(id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	uri := path.Join(apiNetworkBase, id)
	method := http.MethodGet
	return c.waitFor404Status(uri, method)
}

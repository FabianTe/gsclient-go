package gsclient

import (
	"errors"
	"fmt"
	"net/http"
	"path"
)

//TemplateList JSON struct of a list of templates
type TemplateList struct {
	//Array of templates
	List map[string]TemplateProperties `json:"templates"`
}

//DeletedTemplateList JSON struct of a list of deleted templates
type DeletedTemplateList struct {
	//Array of deleted templates
	List map[string]TemplateProperties `json:"deleted_templates"`
}

//Template JSON struct of a single template
type Template struct {
	//Properties of a template
	Properties TemplateProperties `json:"template"`
}

//TemplateProperties JSOn struct of properties of a template
type TemplateProperties struct {
	//Status indicates the status of the object.
	Status string `json:"status"`

	//Status indicates the status of the object.
	Ostype string `json:"ostype"`

	//Helps to identify which datacenter an object belongs to.
	LocationUUID string `json:"location_uuid"`

	//Description of the Template.
	Version string `json:"version"`

	//Description of the Template.
	LocationIata string `json:"location_iata"`

	//Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	//the object is private, the value will be true. Otherwise the value will be false.
	Private bool `json:"private"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//If a template has been used that requires a license key (e.g. Windows Servers)
	//this shows the product_no of the license (see the /prices endpoint for more details).
	LicenseProductNo int `json:"license_product_no"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//Total minutes the object has been running.
	UsageInMinutes int `json:"usage_in_minutes"`

	//The capacity of a storage/ISO-Image/template/snapshot in GB.
	Capacity int `json:"capacity"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationName string `json:"location_name"`

	//The OS distrobution that the Template contains.
	Distro string `json:"distro"`

	//Description of the Template.
	Description string `json:"description"`

	//The price for the current period since the last bill.
	CurrentPrice float64 `json:"current_price"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationCountry string `json:"location_country"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//List of labels.
	Labels []string `json:"labels"`
}

//TemplateCreateRequest JSON struct of a request for creating a template
type TemplateCreateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//snapshot uuid for template
	SnapshotUUID string `json:"snapshot_uuid"`

	//List of labels. Optional.
	Labels []string `json:"labels,omitempty"`
}

//TemplateUpdateRequest JSON struct of a request for updating a template
type TemplateUpdateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	//Optional.
	Name string `json:"name,omitempty"`

	//List of labels. Optional.
	Labels []string `json:"labels,omitempty"`
}

//GetTemplate gets a template
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getTemplate
func (c *Client) GetTemplate(id string) (Template, error) {
	if !isValidUUID(id) {
		return Template{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiTemplateBase, id),
		method: http.MethodGet,
	}
	var response Template
	err := r.execute(*c, &response)
	return response, err
}

//GetTemplateList gets a list of templates
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getTemplates
func (c *Client) GetTemplateList() ([]Template, error) {
	r := Request{
		uri:    apiTemplateBase,
		method: http.MethodGet,
	}
	var response TemplateList
	var templates []Template
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		templates = append(templates, Template{
			Properties: properties,
		})
	}
	return templates, err
}

//GetTemplateByName gets a template by its name
func (c *Client) GetTemplateByName(name string) (Template, error) {
	if name == "" {
		return Template{}, errors.New("'name' is required")
	}
	templates, err := c.GetTemplateList()
	if err != nil {
		return Template{}, err
	}
	for _, template := range templates {
		if template.Properties.Name == name {
			return Template{Properties: template.Properties}, nil
		}
	}
	return Template{}, fmt.Errorf("Template %v not found", name)
}

//CreateTemplate creates a template
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createTemplate
func (c *Client) CreateTemplate(body TemplateCreateRequest) (CreateResponse, error) {
	r := Request{
		uri:    apiTemplateBase,
		method: http.MethodPost,
		body:   body,
	}
	var response CreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return CreateResponse{}, err
	}
	if c.cfg.sync {
		err = c.waitForRequestCompleted(response.RequestUUID)
	}
	return response, err
}

//UpdateTemplate updates a template
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateTemplate
func (c *Client) UpdateTemplate(id string, body TemplateUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiTemplateBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForTemplateActive(id)
	}
	return r.execute(*c, nil)
}

//DeleteTemplate deletes a template
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteTemplate
func (c *Client) DeleteTemplate(id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiTemplateBase, id),
		method: http.MethodDelete,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForTemplateDeleted(id)
	}
	return r.execute(*c, nil)
}

//GetTemplateEventList gets a list of a template's events
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getTemplateEvents
func (c *Client) GetTemplateEventList(id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiTemplateBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var templateEvents []Event
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		templateEvents = append(templateEvents, Event{Properties: properties})
	}
	return templateEvents, err
}

//GetTemplatesByLocation gets a list of templates by location
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLocationTemplates
func (c *Client) GetTemplatesByLocation(id string) ([]Template, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLocationBase, id, "templates"),
		method: http.MethodGet,
	}
	var response TemplateList
	var templates []Template
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		templates = append(templates, Template{Properties: properties})
	}
	return templates, err
}

//GetDeletedTemplates gets a list of deleted templates
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedTemplates
func (c *Client) GetDeletedTemplates() ([]Template, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "templates"),
		method: http.MethodGet,
	}
	var response DeletedTemplateList
	var templates []Template
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		templates = append(templates, Template{Properties: properties})
	}
	return templates, err
}

//waitForTemplateActive allows to wait until the template's status is active
func (c *Client) waitForTemplateActive(id string) error {
	return retryWithTimeout(func() (bool, error) {
		template, err := c.GetTemplate(id)
		return template.Properties.Status != resourceActiveStatus, err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitForTemplateDeleted allows to wait until the template is deleted
func (c *Client) waitForTemplateDeleted(id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	uri := path.Join(apiTemplateBase, id)
	method := http.MethodGet
	return c.waitFor404Status(uri, method)
}

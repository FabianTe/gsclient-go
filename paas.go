package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//PaaSServices is the JSON struct of a list of PaaS services
type PaaSServices struct {
	//Array of PaaS services
	List map[string]PaaSServiceProperties `json:"paas_services"`
}

//DeletedPaaSServices is the JSON struct of a list of deleted PaaS services
type DeletedPaaSServices struct {
	//Array of deleted PaaS services
	List map[string]PaaSServiceProperties `json:"deleted_paas_services"`
}

//PaaSService is the JSON struct of a single PaaS service
type PaaSService struct {
	//Properties of a PaaS service
	Properties PaaSServiceProperties `json:"paas_service"`
}

//PaaSServiceProperties is the properties of a single PaaS service
type PaaSServiceProperties struct {
	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//List of labels.
	Labels []string `json:"labels"`

	//Contains the initial setup credentials for Service.
	Credentials []Credential `json:"credentials"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//Contains the IPv6 address and port that the Service will listen to,
	//you can use these details to connect internally to a service.
	ListenPorts map[string]map[string]int `json:"listen_ports"`

	//The UUID of the security zone that the service is running in.
	SecurityZoneUUID string `json:"security_zone_uuid"`

	//The template used to create the service, you can find an available list at the /service_templates endpoint.
	ServiceTemplateUUID string `json:"service_template_uuid"`

	//Total minutes the object has been running.
	UsageInMinutes int `json:"usage_in_minutes"`

	//The price for the current period since the last bill.
	CurrentPrice float64 `json:"current_price"`

	//Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	//Status indicates the status of the object.
	Status string `json:"status"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//A list of service resource limits.
	ResourceLimits []ResourceLimit `json:"resource_limits"`

	//Contains the service parameters for the service.
	Parameters map[string]interface{} `json:"parameters"`
}

//Credential is JSON struct of credential
type Credential struct {
	//The initial username to authenticate the Service.
	Username string `json:"username"`

	//The initial password to authenticate the Service.
	Password string `json:"password"`

	//The type of Service.
	Type string `json:"type"`
}

//PaaSServiceCreateRequest is JSON struct of a request for creating a PaaS service
type PaaSServiceCreateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//The template used to create the service, you can find an available list at the /service_templates endpoint.
	PaaSServiceTemplateUUID string `json:"paas_service_template_uuid"`

	//The template used to create the service, you can find an available list at the /service_templates endpoint.
	Labels []string `json:"labels,omitempty"`

	//The UUID of the security zone that the service is running in.
	PaaSSecurityZoneUUID string `json:"paas_security_zone_uuid,omitempty"`

	//A list of service resource limits.
	ResourceLimits []ResourceLimit `json:"resource_limits,omitempty"`

	//Contains the service parameters for the service.
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

//ResourceLimit is JSON struct of resource limit
type ResourceLimit struct {
	//The name of the resource you would like to cap.
	Resource string `json:"resource"`

	//The maximum number of the specific resource your service can use.
	Limit int `json:"limit"`
}

//PaaSServiceCreateResponse is JSON struct of a response for creating a PaaS service
type PaaSServiceCreateResponse struct {
	//UUID of the request
	RequestUUID string `json:"request_uuid"`

	//Contains the IPv6 address and port that the Service will listen to, you can use these details to connect internally to a service.
	ListenPorts map[string]map[string]int `json:"listen_ports"`

	//The template used to create the service, you can find an available list at the /service_templates endpoint.
	PaaSServiceUUID string `json:"paas_service_uuid"`

	//Contains the initial setup credentials for Service.
	Credentials []Credential `json:"credentials"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//A list of service resource limits.
	ResourceLimits []ResourceLimit `json:"resource_limits"`

	//Contains the service parameters for the service.
	Parameters map[string]interface{} `json:"parameters"`
}

//PaaSTemplates is JSON struct of a list of PaaS Templates
type PaaSTemplates struct {
	//Array of PaaS templates
	List map[string]PaaSTemplateProperties `json:"paas_service_templates"`
}

//PaaSTemplate is JSON struct for a single PaaS Template
type PaaSTemplate struct {
	//Properties of a PaaS template
	Properties PaaSTemplateProperties `json:"paas_service_template"`
}

//PaaSTemplateProperties is JSOn struct of properties of a PaaS template
type PaaSTemplateProperties struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//Describes the category of the service.
	Category string `json:"category"`

	//Product No
	ProductNo int `json:"product_no"`

	//List of labels.
	Labels []string `json:"labels"`

	//The amount of concurrent connections for the service.
	Resources Resource `json:"resources"`

	//Status indicates the status of the object.
	Status string `json:"status"`

	//A definition of possible service template parameters (python-cerberus compatible).
	ParametersSchema map[string]Parameter `json:"parameters_schema"`
}

//Parameter JSON of a parameter
type Parameter struct {
	//Is required
	Required bool `json:"required"`

	//Is empty
	Empty bool `json:"empty"`

	//Description of parameter
	Description string `json:"description"`

	//Maximum
	Max int `json:"max"`

	//Minimum
	Min int `json:"min"`

	//Default value
	Default interface{} `json:"default"`

	//Type of parameter
	Type string `json:"type"`

	//Allowed values
	Allowed []string `json:"allowed"`

	//Regex
	Regex string `json:"regex"`
}

//Resource JSON of a resource
type Resource struct {
	//The amount of memory required by the service, either RAM(MB) or SSD Storage(GB).
	Memory int `json:"memory"`

	//The amount of concurrent connections for the service.
	Connections int `json:"connections"`
}

//PaaSServiceUpdateRequest JSON of a request for updating a PaaS service
type PaaSServiceUpdateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	//Leave it if you do not want to update the name
	Name string `json:"name,omitempty"`

	//List of labels. Leave it if you do not want to update the list of labels
	Labels []string `json:"labels,omitempty"`

	//Contains the service parameters for the service. Leave it if you do not want to update the parameters
	Parameters map[string]interface{} `json:"parameters,omitempty"`

	//A list of service resource limits. Leave it if you do not want to update the resource limits
	ResourceLimits []ResourceLimit `json:"resource_limits,omitempty"`
}

//PaaSServiceMetrics JSON of a list of PaaS metrics
type PaaSServiceMetrics struct {
	//Array of a PaaS service's metrics
	List []PaaSMetricProperties `json:"paas_service_metrics"`
}

//PaaSServiceMetric JSON of a single PaaS metric
type PaaSServiceMetric struct {
	//Properties of a PaaS service metric
	Properties PaaSMetricProperties `json:"paas_service_metric"`
}

//PaaSMetricProperties JSON of properties of a PaaS metric
type PaaSMetricProperties struct {
	//Defines the begin of the time range.
	BeginTime GSTime `json:"begin_time"`

	//Defines the end of the time range.
	EndTime GSTime `json:"end_time"`

	//The UUID of an object is always unique, and refers to a specific object.
	PaaSServiceUUID string `json:"paas_service_uuid"`

	//CPU core usage
	CoreUsage PaaSMetricValue `json:"core_usage"`

	//Storage usage
	StorageSize PaaSMetricValue `json:"storage_size"`
}

//PaaSMetricValue JSON of a metric value
type PaaSMetricValue struct {
	//Value
	Value float64 `json:"value"`

	//Unit of the value
	Unit string `json:"unit"`
}

//PaaSSecurityZones JSON struct of a list of PaaS security zones
type PaaSSecurityZones struct {
	//Array of security zones
	List map[string]PaaSSecurityZoneProperties `json:"paas_security_zones"`
}

//PaaSSecurityZone JSON struct of a single PaaS security zone
type PaaSSecurityZone struct {
	//Properties of a security zone
	Properties PaaSSecurityZoneProperties `json:"paas_security_zone"`
}

//PaaSSecurityZoneProperties JSOn struct of properties of a PaaS security zone
type PaaSSecurityZoneProperties struct {
	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationCountry string `json:"location_country"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//Uses IATA airport code, which works as a location identifier.
	LocationIata string `json:"location_iata"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//List of labels.
	Labels []string `json:"labels"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationName string `json:"location_name"`

	//Status indicates the status of the object.
	Status string `json:"status"`

	//Helps to identify which datacenter an object belongs to.
	LocationUUID string `json:"location_uuid"`

	//Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//object (PaaSRelationService)
	Relation PaaSRelationService `json:"relation"`
}

//PaaSRelationService JSON struct of a relation between a PaaS service and a service
type PaaSRelationService struct {
	//Array of object (ServiceObject)
	Services []ServiceObject `json:"services"`
}

//ServiceObject JSON struct of a service object
type ServiceObject struct {
	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`
}

//PaaSSecurityZoneCreateRequest JSON struct of a request for creating a PaaS security zone
type PaaSSecurityZoneCreateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name,omitempty"`

	//Identifies which datacenter the object belongs to.
	LocationUUID string `json:"location_uuid,omitempty"`
}

//PaaSSecurityZoneCreateResponse JSON struct of a response for creating a PaaS security zone
type PaaSSecurityZoneCreateResponse struct {
	//UUID of the request
	RequestUUID string `json:"request_uuid"`

	//UUID of the security zone being created
	PaaSSecurityZoneUUID string `json:"paas_security_zone_uuid"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`
}

//PaaSSecurityZoneUpdateRequest JSON struct of a request for updating a PaaS security zone
type PaaSSecurityZoneUpdateRequest struct {
	//The new name you give to the security zone. Leave it if you do not want to update the name
	Name string `json:"name,omitempty"`

	//Identifies which datacenter the object belongs to. Leave it if you do not want to update the location
	LocationUUID string `json:"location_uuid,omitempty"`

	//The UUID for the security zone you would like to update. Leave it if you do not want to update the security zone
	PaaSSecurityZoneUUID string `json:"paas_security_zone_uuid,omitempty"`
}

//GetPaaSServiceList returns a list of PaaS Services
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getPaasServices
func (c *Client) GetPaaSServiceList() ([]PaaSService, error) {
	r := Request{
		uri:    path.Join(apiPaaSBase, "services"),
		method: http.MethodGet,
	}
	var response PaaSServices
	var paasServices []PaaSService
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		paasServices = append(paasServices, PaaSService{
			Properties: properties,
		})
	}
	return paasServices, err
}

//CreatePaaSService creates a new PaaS service
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createPaasService
func (c *Client) CreatePaaSService(body PaaSServiceCreateRequest) (PaaSServiceCreateResponse, error) {
	r := Request{
		uri:    path.Join(apiPaaSBase, "services"),
		method: http.MethodPost,
		body:   body,
	}
	var response PaaSServiceCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return PaaSServiceCreateResponse{}, err
	}
	if c.cfg.sync {
		err = c.waitForRequestCompleted(response.RequestUUID)
	}
	return response, err
}

//GetPaaSService returns a specific PaaS Service based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getPaasService
func (c *Client) GetPaaSService(id string) (PaaSService, error) {
	if !isValidUUID(id) {
		return PaaSService{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiPaaSBase, "services", id),
		method: http.MethodGet,
	}
	var response PaaSService
	err := r.execute(*c, &response)
	return response, err
}

//UpdatePaaSService updates a specific PaaS Service based on a given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updatePaasService
func (c *Client) UpdatePaaSService(id string, body PaaSServiceUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiPaaSBase, "services", id),
		method: http.MethodPatch,
		body:   body,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForPaaSServiceActive(id)
	}
	return r.execute(*c, nil)
}

//DeletePaaSService deletes a PaaS service
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deletePaasService
func (c *Client) DeletePaaSService(id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiPaaSBase, "services", id),
		method: http.MethodDelete,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForPaaSServiceDeleted(id)
	}
	return r.execute(*c, nil)
}

//GetPaaSServiceMetrics get a specific PaaS Service's metrics based on a given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getPaasServiceMetrics
func (c *Client) GetPaaSServiceMetrics(id string) ([]PaaSServiceMetric, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiPaaSBase, "services", id, "metrics"),
		method: http.MethodGet,
	}
	var response PaaSServiceMetrics
	var metrics []PaaSServiceMetric
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		metrics = append(metrics, PaaSServiceMetric{
			Properties: properties,
		})
	}
	return metrics, err
}

//GetPaaSTemplateList returns a list of PaaS service templates
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getPaasServiceTemplates
func (c *Client) GetPaaSTemplateList() ([]PaaSTemplate, error) {
	r := Request{
		uri:    path.Join(apiPaaSBase, "service_templates"),
		method: http.MethodGet,
	}
	var response PaaSTemplates
	var paasTemplates []PaaSTemplate
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		paasTemplate := PaaSTemplate{
			Properties: properties,
		}
		paasTemplates = append(paasTemplates, paasTemplate)
	}
	return paasTemplates, err
}

//GetPaaSSecurityZoneList get available security zones
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getPaasSecurityZones
func (c *Client) GetPaaSSecurityZoneList() ([]PaaSSecurityZone, error) {
	r := Request{
		uri:    path.Join(apiPaaSBase, "security_zones"),
		method: http.MethodGet,
	}
	var response PaaSSecurityZones
	var securityZones []PaaSSecurityZone
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		securityZones = append(securityZones, PaaSSecurityZone{
			Properties: properties,
		})
	}
	return securityZones, err
}

//CreatePaaSSecurityZone creates a new PaaS security zone
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createPaasSecurityZone
func (c *Client) CreatePaaSSecurityZone(body PaaSSecurityZoneCreateRequest) (PaaSSecurityZoneCreateResponse, error) {
	r := Request{
		uri:    path.Join(apiPaaSBase, "security_zones"),
		method: http.MethodPost,
		body:   body,
	}
	var response PaaSSecurityZoneCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return PaaSSecurityZoneCreateResponse{}, err
	}
	if c.cfg.sync {
		err = c.waitForRequestCompleted(response.RequestUUID)
	}
	return response, err
}

//GetPaaSSecurityZone get a specific PaaS Security Zone based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getPaasSecurityZone
func (c *Client) GetPaaSSecurityZone(id string) (PaaSSecurityZone, error) {
	if !isValidUUID(id) {
		return PaaSSecurityZone{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiPaaSBase, "security_zones", id),
		method: http.MethodGet,
	}
	var response PaaSSecurityZone
	err := r.execute(*c, &response)
	return response, err
}

//UpdatePaaSSecurityZone update a specific PaaS security zone based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updatePaasSecurityZone
func (c *Client) UpdatePaaSSecurityZone(id string, body PaaSSecurityZoneUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiPaaSBase, "security_zones", id),
		method: http.MethodPatch,
		body:   body,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForSecurityZoneActive(id)
	}
	return r.execute(*c, nil)
}

//DeletePaaSSecurityZone delete a specific PaaS Security Zone based on given id
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deletePaasSecurityZone
func (c *Client) DeletePaaSSecurityZone(id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiPaaSBase, "security_zones", id),
		method: http.MethodDelete,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForSecurityZoneDeleted(id)
	}
	return r.execute(*c, nil)
}

//GetDeletedPaaSServices returns a list of deleted PaaS Services
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedPaasServices
func (c *Client) GetDeletedPaaSServices() ([]PaaSService, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "paas_services"),
		method: http.MethodGet,
	}
	var response DeletedPaaSServices
	var paasServices []PaaSService
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		paasServices = append(paasServices, PaaSService{
			Properties: properties,
		})
	}
	return paasServices, err
}

//waitForPaaSServiceActive allows to wait until the PaaS service's status is active
func (c *Client) waitForPaaSServiceActive(id string) error {
	return retryWithTimeout(func() (bool, error) {
		paas, err := c.GetPaaSService(id)
		return paas.Properties.Status != resourceActiveStatus, err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitForPaaSServiceDeleted allows to wait until the PaaS service is deleted
func (c *Client) waitForPaaSServiceDeleted(id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	uri := path.Join(apiPaaSBase, "services", id)
	method := http.MethodGet
	return c.waitFor404Status(uri, method)
}

//waitForSecurityZoneActive allows to wait until the security zone's status is active
func (c *Client) waitForSecurityZoneActive(id string) error {
	return retryWithTimeout(func() (bool, error) {
		secZone, err := c.GetPaaSSecurityZone(id)
		return secZone.Properties.Status != resourceActiveStatus, err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitForSecurityZoneDeleted allows to wait until the security zone is deleted
func (c *Client) waitForSecurityZoneDeleted(id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	uri := path.Join(apiPaaSBase, "security_zones", id)
	method := http.MethodGet
	return c.waitFor404Status(uri, method)
}

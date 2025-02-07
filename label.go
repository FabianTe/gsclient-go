package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//LabelList JSON struct of a list of labels
type LabelList struct {
	//List of labels
	List map[string]LabelProperties `json:"labels"`
}

//Label JSON struct of a single label
type Label struct {
	//Properties of a label
	Properties LabelProperties `json:"label"`
}

//LabelProperties JSON struct of properties of a label
type LabelProperties struct {
	//Label's name
	Label string `json:"label"`

	//Create time of a label
	CreateTime GSTime `json:"create_time"`

	//Time of the last change of a label
	ChangeTime GSTime `json:"change_time"`

	//Relations of a label
	Relations []interface{} `json:"relations"`

	//Status indicates the status of a label.
	Status string `json:"status"`
}

//LabelCreateRequest JSON struct of a request for creating a label
type LabelCreateRequest struct {
	//Name of the new label
	Label string `json:"label"`
}

//GetLabelList gets a list of available labels
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/GetLabels
func (c *Client) GetLabelList() ([]Label, error) {
	r := Request{
		uri:    apiLabelBase,
		method: http.MethodGet,
	}
	var response LabelList
	var labels []Label
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		labels = append(labels, Label{Properties: properties})
	}
	return labels, err
}

//CreateLabel creates a new label
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/CreateLabel
func (c *Client) CreateLabel(body LabelCreateRequest) (CreateResponse, error) {
	r := Request{
		uri:    apiLabelBase,
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

//DeleteLabel deletes a label
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/DeleteLabel
func (c *Client) DeleteLabel(label string) error {
	if label == "" {
		return errors.New("'label' is required")
	}
	r := Request{
		uri:    path.Join(apiLabelBase, label),
		method: http.MethodDelete,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		return c.waitForLabelDeleted(label)
	}
	return r.execute(*c, nil)
}

//waitForLabelDeleted allows to wait until the label is deleted
func (c *Client) waitForLabelDeleted(label string) error {
	if label == "" {
		return errors.New("'label' is required")
	}
	return retryWithTimeout(func() (bool, error) {
		labels, err := c.GetLabelList()
		return isLabelInSlice(label, labels), err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//isLabelInSlice check if a label in a lice of labels
func isLabelInSlice(a string, list []Label) bool {
	for _, b := range list {
		if b.Properties.Label == a {
			return true
		}
	}
	return false
}

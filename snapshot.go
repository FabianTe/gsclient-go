package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//StorageSnapshotList is JSON structure of a list of storage snapshots
type StorageSnapshotList struct {
	//Array of snapshots
	List map[string]StorageSnapshotProperties `json:"snapshots"`
}

//DeletedStorageSnapshotList is JSON structure of a list of deleted storage snapshots
type DeletedStorageSnapshotList struct {
	//Array of deleted snapshots
	List map[string]StorageSnapshotProperties `json:"deleted_snapshots"`
}

//StorageSnapshot is JSON structure of a single storage snapshot
type StorageSnapshot struct {
	//Properties of a snapshot
	Properties StorageSnapshotProperties `json:"snapshot"`
}

//StorageSnapshotProperties JSON struct of properties of a storage snapshot
type StorageSnapshotProperties struct {
	//List of labels.
	Labels []string `json:"labels"`

	//The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	Name string `json:"name"`

	//Status indicates the status of the object.
	Status string `json:"status"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationCountry string `json:"location_country"`

	//Total minutes the object has been running.
	UsageInMinutes int `json:"usage_in_minutes"`

	//Helps to identify which datacenter an object belongs to.
	LocationUUID string `json:"location_uuid"`

	//Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	//If a template has been used that requires a license key (e.g. Windows Servers)
	//this shows the product_no of the license (see the /prices endpoint for more details).
	LicenseProductNo int `json:"license_product_no"`

	//The price for the current period since the last bill.
	CurrentPrice float64 `json:"current_price"`

	//Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	//The capacity of a storage/ISO-Image/template/snapshot in GB.
	Capacity int `json:"capacity"`

	//The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
	LocationName string `json:"location_name"`

	//Uses IATA airport code, which works as a location identifier.
	LocationIata string `json:"location_iata"`

	//Uuid of the storage used to create this snapshot
	ParentUUID string `json:"parent_uuid"`
}

//StorageSnapshotCreateRequest JSON struct of a request for creating a storage snapshot
type StorageSnapshotCreateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	//Optional
	Name string `json:"name,omitempty"`

	//List of labels. Optional.
	Labels []string `json:"labels,omitempty"`
}

//StorageSnapshotCreateResponse JSON struct of a response for creating a storage snapshot
type StorageSnapshotCreateResponse struct {
	//UUID of the request
	RequestUUID string `json:"request_uuid"`

	//UUID of the snapshot being created
	ObjectUUID string `json:"object_uuid"`
}

//StorageSnapshotUpdateRequest JSON struct of a request for updating a storage snapshot
type StorageSnapshotUpdateRequest struct {
	//The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
	//Optional
	Name string `json:"name,omitempty"`

	//List of labels. Optional.
	Labels []string `json:"labels,omitempty"`
}

//StorageRollbackRequest JSON struct of a request for rolling back
type StorageRollbackRequest struct {
	//Rollback=true => storage will be restored
	Rollback bool `json:"rollback,omitempty"`
}

//StorageSnapshotExportToS3Request JSON struct of a request for exporting a storage snapshot to S3
type StorageSnapshotExportToS3Request struct {
	S3auth struct {
		//Host of S3
		Host string `json:"host"`

		//Access key of S3
		AccessKey string `json:"access_key"`

		//Secret key of S3
		SecretKey string `json:"secret_key"`
	} `json:"s3auth"`
	S3data struct {
		//Host of S3
		Host string `json:"host"`

		//Bucket that file will be uploaded to
		Bucket string `json:"bucket"`

		//Name of the file being uploaded
		Filename string `json:"filename"`

		//Is the file private?
		Private bool `json:"private"`
	} `json:"s3data"`
}

//GetStorageSnapshotList gets a list of storage snapshots
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getSnapshots
func (c *Client) GetStorageSnapshotList(id string) ([]StorageSnapshot, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, id, "snapshots"),
		method: http.MethodGet,
	}
	var response StorageSnapshotList
	var snapshots []StorageSnapshot
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		snapshots = append(snapshots, StorageSnapshot{Properties: properties})
	}
	return snapshots, err
}

//GetStorageSnapshot gets a specific storage's snapshot based on given storage id and snapshot id.
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getSnapshot
func (c *Client) GetStorageSnapshot(storageID, snapshotID string) (StorageSnapshot, error) {
	if !isValidUUID(storageID) || !isValidUUID(snapshotID) {
		return StorageSnapshot{}, errors.New("'storageID' or 'snapshotID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, storageID, "snapshots", snapshotID),
		method: http.MethodGet,
	}
	var response StorageSnapshot
	err := r.execute(*c, &response)
	return response, err
}

//CreateStorageSnapshot creates a new storage's snapshot
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/createSnapshot
func (c *Client) CreateStorageSnapshot(id string, body StorageSnapshotCreateRequest) (StorageSnapshotCreateResponse, error) {
	if !isValidUUID(id) {
		return StorageSnapshotCreateResponse{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, id, "snapshots"),
		method: http.MethodPost,
		body:   body,
	}
	var response StorageSnapshotCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return StorageSnapshotCreateResponse{}, err
	}
	if c.cfg.sync {
		err = c.waitForRequestCompleted(response.RequestUUID)
	}
	return response, err
}

//UpdateStorageSnapshot updates a specific storage's snapshot
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/updateSnapshot
func (c *Client) UpdateStorageSnapshot(storageID, snapshotID string, body StorageSnapshotUpdateRequest) error {
	if !isValidUUID(storageID) || !isValidUUID(snapshotID) {
		return errors.New("'storageID' or 'snapshotID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, storageID, "snapshots", snapshotID),
		method: http.MethodPatch,
		body:   body,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForSnapshotActive(storageID, snapshotID)
	}
	return r.execute(*c, nil)
}

//DeleteStorageSnapshot deletes a specific storage's snapshot
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/deleteSnapshot
func (c *Client) DeleteStorageSnapshot(storageID, snapshotID string) error {
	if !isValidUUID(storageID) || !isValidUUID(snapshotID) {
		return errors.New("'storageID' or 'snapshotID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, storageID, "snapshots", snapshotID),
		method: http.MethodDelete,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForSnapshotDeleted(storageID, snapshotID)
	}
	return r.execute(*c, nil)
}

//RollbackStorage rollbacks a storage
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/StorageRollback
func (c *Client) RollbackStorage(storageID, snapshotID string, body StorageRollbackRequest) error {
	if !isValidUUID(storageID) || !isValidUUID(snapshotID) {
		return errors.New("'storageID' or 'snapshotID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, storageID, "snapshots", snapshotID, "rollback"),
		method: http.MethodPatch,
		body:   body,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForSnapshotActive(storageID, snapshotID)
	}
	return r.execute(*c, nil)
}

//ExportStorageSnapshotToS3 export a storage's snapshot to S3
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/SnapshotExportToS3
func (c *Client) ExportStorageSnapshotToS3(storageID, snapshotID string, body StorageSnapshotExportToS3Request) error {
	if !isValidUUID(storageID) || !isValidUUID(snapshotID) {
		return errors.New("'storageID' and 'snapshotID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, storageID, "snapshots", snapshotID, "export_to_s3"),
		method: http.MethodPatch,
		body:   body,
	}
	if c.cfg.sync {
		err := r.execute(*c, nil)
		if err != nil {
			return err
		}
		//Block until the request is finished
		return c.waitForSnapshotActive(storageID, snapshotID)
	}
	return r.execute(*c, nil)
}

//GetSnapshotsByLocation gets a list of storage snapshots by location
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getLocationSnapshots
func (c *Client) GetSnapshotsByLocation(id string) ([]StorageSnapshot, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLocationBase, id, "snapshots"),
		method: http.MethodGet,
	}
	var response StorageSnapshotList
	var snapshots []StorageSnapshot
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		snapshots = append(snapshots, StorageSnapshot{Properties: properties})
	}
	return snapshots, err
}

//GetDeletedSnapshots gets a list of deleted storage snapshots
//
//See: https://gridscale.io/en//api-documentation/index.html#operation/getDeletedSnapshots
func (c *Client) GetDeletedSnapshots() ([]StorageSnapshot, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "snapshots"),
		method: http.MethodGet,
	}
	var response DeletedStorageSnapshotList
	var snapshots []StorageSnapshot
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		snapshots = append(snapshots, StorageSnapshot{Properties: properties})
	}
	return snapshots, err
}

//waitForSnapshotActive allows to wait until the snapshot's status is active
func (c *Client) waitForSnapshotActive(storageID, snapshotID string) error {
	return retryWithTimeout(func() (bool, error) {
		snapshot, err := c.GetStorageSnapshot(storageID, snapshotID)
		return snapshot.Properties.Status != resourceActiveStatus, err
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitForSnapshotDeleted allows to wait until the snapshot is deleted
func (c *Client) waitForSnapshotDeleted(storageID, snapshotID string) error {
	if !isValidUUID(storageID) || !isValidUUID(snapshotID) {
		return errors.New("'storageID' or 'snapshotID' is invalid")
	}
	uri := path.Join(apiStorageBase, storageID, "snapshots", snapshotID)
	method := http.MethodGet
	return c.waitFor404Status(uri, method)
}

package endpoints

import (
	"errors"
	"strconv"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
)

// Interface describes the availability of a specific service endpoint.
type Interface string

const (
	// InterfaceAdmin makes an endpoint only available to administrators.
	InterfaceAdmin Interface = "admin"

	// InterfacePublic makes an endpoint available to everyone.
	InterfacePublic Interface = "public"

	// InterfaceInternal makes an endpoint only available within the cluster.
	InterfaceInternal Interface = "internal"
)

// EndpointOpts contains the subset of Endpoint attributes that should be used to create or update an Endpoint.
type EndpointOpts struct {
	Interface Interface
	Name      string
	Region    string
	URL       string
	ServiceID string
}

// Create inserts a new Endpoint into the service catalog.
// Within EndpointOpts, Region may be omitted by being left as "", but all other fields are required.
func Create(client *gophercloud.ServiceClient, opts EndpointOpts) (*Endpoint, error) {
	// Redefined so that Region can be re-typed as a *string, which can be omitted from the JSON output.
	type endpoint struct {
		Interface string  `json:"interface"`
		Name      string  `json:"name"`
		Region    *string `json:"region,omitempty"`
		URL       string  `json:"url"`
		ServiceID string  `json:"service_id"`
	}

	type request struct {
		Endpoint endpoint `json:"endpoint"`
	}

	type response struct {
		Endpoint Endpoint `json:"endpoint"`
	}

	// Ensure that EndpointOpts is fully populated.
	if opts.Interface == "" {
		return nil, ErrInterfaceRequired
	}
	if opts.Name == "" {
		return nil, ErrNameRequired
	}
	if opts.URL == "" {
		return nil, ErrURLRequired
	}
	if opts.ServiceID == "" {
		return nil, ErrServiceIDRequired
	}

	// Populate the request body.
	reqBody := request{
		Endpoint: endpoint{
			Interface: string(opts.Interface),
			Name:      opts.Name,
			URL:       opts.URL,
			ServiceID: opts.ServiceID,
		},
	}

	if opts.Region != "" {
		reqBody.Endpoint.Region = &opts.Region
	}

	var respBody response
	_, err := perigee.Request("POST", getListURL(client), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &respBody,
		OkCodes:     []int{201},
	})
	if err != nil {
		return nil, err
	}

	return &respBody.Endpoint, nil
}

// ListOpts allows finer control over the the endpoints returned by a List call.
// All fields are optional.
type ListOpts struct {
	Interface Interface
	ServiceID string
	Page      int
	PerPage   int
}

// List enumerates endpoints in a paginated collection, optionally filtered by ListOpts criteria.
func List(client *gophercloud.ServiceClient, opts ListOpts) (*EndpointList, error) {
	return nil, errors.New("Not implemented")
}

// Update changes an existing endpoint with new data.
func Update(client *gophercloud.ServiceClient, endpointID string, opts EndpointOpts) (*Endpoint, error) {
	return nil, errors.New("Not implemented")
}

// Delete removes an endpoint from the service catalog.
func Delete(client *gophercloud.ServiceClient, endpointID string) error {
	return errors.New("Not implemented")
}

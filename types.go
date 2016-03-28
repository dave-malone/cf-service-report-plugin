package main

import "fmt"

type pagedResponse struct {
	TotalResults int    `json:"total_results"`
	TotalPages   int    `json:"total_pages"`
	PrevURL      string `json:"prev_url"`
	NextURL      string `json:"next_url"`
}

type orgs struct {
	*pagedResponse
	Resources []orgResource `json:"resources"`
}

type services struct {
	*pagedResponse
	Resources []serviceResource `json:"resources"`
}

type servicePlans struct {
	*pagedResponse
	Resources []servicePlanResource `json:"resources"`
}

type serviceInstances struct {
	*pagedResponse
	Resources []serviceInstanceResource `json:"resources"`
}

type serviceBindings struct {
	*pagedResponse
	Resources []serviceBindingResource `json:"resources"`
}

type resourceMetadata struct {
	URL  string `json:"url"`
	GUID string `json:"guid"`
}

type orgResource struct {
	Entity   org              `json:"entity"`
	Metadata resourceMetadata `json:"metadata"`
}

type serviceResource struct {
	Entity   service          `json:"entity"`
	Metadata resourceMetadata `json:"metadata"`
}

type servicePlanResource struct {
	Entity   servicePlan      `json:"entity"`
	Metadata resourceMetadata `json:"metadata"`
}

type serviceInstanceResource struct {
	Entity   serviceInstance  `json:"entity"`
	Metadata resourceMetadata `json:"metadata"`
}

type serviceBindingResource struct {
	Entity   serviceBinding   `json:"entity"`
	Metadata resourceMetadata `json:"metadata"`
}

type org struct {
	Name string `json:"name"`
}

type service struct {
	Label           string `json:"label"`
	ServicePlansURL string `json:"service_plans_url"`
}

type servicePlan struct {
	Name                string `json:"name"`
	ServiceInstancesURL string `json:"service_instances_url"`
}

type serviceInstance struct {
	Name               string `json:"name"`
	Type               string `json:"type"`
	ServiceBindingsURL string `json:"service_bindings_url"`
}

type serviceBinding struct {
	AppGUID             string `json:"app_guid"`
	ServiceInstanceGUID string `json:"service_instance_guid"`
	AppURL              string `json:"app_url"`
	ServiceInstanceURL  string `json:"service_instance_url"`
}

func (s *serviceBinding) String() string {
	return fmt.Sprintf("Service Instance GUID: %s\nApp GUID: %s", s.ServiceInstanceGUID, s.AppGUID)
}

type serviceReport struct {
	results []serviceReportResult
}

type serviceReportResult struct {
	orgName          string
	serviceInstances string
	boundAppGuids    string
}

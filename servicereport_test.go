package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/cloudfoundry/cli/plugin/fakes"
)

func TestGetServicePlans(t *testing.T) {
	conn := new(fakes.FakeCliConnection)

	content, fileReadErr := readFile("test-data/service_plans.json")

	if fileReadErr != nil {
		panic("Failed to read file: " + fileReadErr.Error())
	}

	conn.CliCommandWithoutTerminalOutputReturns(content, nil)

	service := new(service)
	service.Label = "test-service"
	service.ServicePlansURL = "/v2/services/test-guid/service_plans"

	servicePlans, err := getServicePlans(conn, *service)

	if err != nil {
		t.Errorf("getServicePlans return an error: %v", err.Error())
	}

	if len(servicePlans.Resources) == 0 {
		t.Errorf("Expected at least one service plan from the test data, but there was %v", len(servicePlans.Resources))
	}

	for _, servicePlanResource := range servicePlans.Resources {
		if servicePlanResource.Entity.Name == "" {
			t.Error("Name was empty on ServicePlan")
		}

		if servicePlanResource.Entity.ServiceInstancesURL == "" {
			t.Error("ServiceInstancesURL was empty on ServicePlan")
		}
	}
}

func TestGetServiceBindings(t *testing.T) {
	conn := new(fakes.FakeCliConnection)

	content, fileReadErr := readFile("test-data/service_bindings.json")

	if fileReadErr != nil {
		panic("Failed to read file: " + fileReadErr.Error())
	}

	conn.CliCommandWithoutTerminalOutputReturns(content, nil)

	serviceInstance := new(serviceInstance)
	serviceInstance.Name = "test-service-instance"
	serviceInstance.ServiceBindingsURL = "/v2/service_instances/:service-instance-guid/service_bindings"

	serviceBindings, err := getServiceBindings(conn, *serviceInstance)

	if err != nil {
		t.Errorf("getServiceBindings return an error: %v", err.Error())
	}

	if len(serviceBindings.Resources) == 0 {
		t.Errorf("Expected at least one service binding from the test data, but there was %v", len(serviceBindings.Resources))
	}

	for _, serviceBindingResource := range serviceBindings.Resources {
		if serviceBindingResource.Entity.AppGUID == "" {
			t.Error("AppGUID was empty on serviceBinding")
		}

		if serviceBindingResource.Entity.AppURL == "" {
			t.Error("AppURL was empty on serviceBinding")
		}
	}
}

func TestGetOrgs(t *testing.T) {
	conn := new(fakes.FakeCliConnection)

	content, fileReadErr := readFile("test-data/orgs.json")

	if fileReadErr != nil {
		panic("Failed to read file: " + fileReadErr.Error())
	}

	conn.CliCommandWithoutTerminalOutputReturns(content, nil)

	orgs, err := getOrgs(conn)

	if err != nil {
		t.Errorf("getOrgs Returned an error: %v", err.Error())
	}

	if len(orgs.Resources) == 0 {
		t.Error("expected at least one result from getOrgs")
	}

	for _, orgResource := range orgs.Resources {
		if orgResource.Entity.Name == "" {
			t.Error("Name was null on OrgResource.Entity")
		}
	}
}

func TestGetServices(t *testing.T) {
	conn := new(fakes.FakeCliConnection)

	content, fileReadErr := readFile("test-data/org_services.json")

	if fileReadErr != nil {
		panic("Failed to read file: " + fileReadErr.Error())
	}

	conn.CliCommandWithoutTerminalOutputReturns(content, nil)

	orgResource := new(orgResource)
	org := new(org)
	org.Name = "test-org"
	resourceMetadata := new(resourceMetadata)
	resourceMetadata.GUID = "abc123"
	orgResource.Entity = *org
	orgResource.Metadata = *resourceMetadata

	services, err := getServices(conn, *orgResource)

	if err != nil {
		t.Errorf("getServices Returned an error: %v", err.Error())
	}

	if len(services.Resources) == 0 {
		t.Error("expected at least one result from getServices")
	}

	for _, serviceResource := range services.Resources {
		if serviceResource.Entity.Label == "" {
			t.Error("Label was blank on Service")
		}
	}

}

func TestGetServiceInstances(t *testing.T) {
	conn := new(fakes.FakeCliConnection)

	content, fileReadErr := readFile("test-data/service_instances.json")

	if fileReadErr != nil {
		panic("Failed to read file: " + fileReadErr.Error())
	}

	conn.CliCommandWithoutTerminalOutputReturns(content, nil)

	servicePlanResource := new(servicePlanResource)
	servicePlan := new(servicePlan)
	servicePlan.Name = "test-org"
	servicePlan.ServiceInstancesURL = "/v2/service_instances/:guid/service_bindings"
	resourceMetadata := new(resourceMetadata)
	resourceMetadata.GUID = "abc123"
	servicePlanResource.Entity = *servicePlan
	servicePlanResource.Metadata = *resourceMetadata

	serviceInstances, err := getServiceInstances(conn, *servicePlanResource)

	if err != nil {
		t.Errorf("getServiceInstances return an error: %v", err.Error())
	}

	if len(serviceInstances.Resources) == 0 {
		t.Errorf("Expected at least one service instance from the test data, but there was %v", len(serviceInstances.Resources))
	}

	for _, serviceInstanceResource := range serviceInstances.Resources {
		if serviceInstanceResource.Entity.Name == "" {
			t.Error("Name was empty on ServiceInstance")
		}

		if serviceInstanceResource.Entity.Type == "" {
			t.Error("Type was empty on ServiceInstance")
		}

		if serviceInstanceResource.Entity.ServiceBindingsURL == "" {
			t.Error("ServiceBindingsURL was empty on ServiceInstance")
		}
	}
}

func readFile(filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	lines := strings.Split(string(content), "\n")
	return lines, err
}

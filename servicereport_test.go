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

	cmd := new(ServiceReportCmd)

	servicePlans, err := cmd.getServicePlans(conn)

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

func TestGetOrgs(t *testing.T) {
	conn := new(fakes.FakeCliConnection)

	content, fileReadErr := readFile("test-data/orgs.json")

	if fileReadErr != nil {
		panic("Failed to read file: " + fileReadErr.Error())
	}

	conn.CliCommandWithoutTerminalOutputReturns(content, nil)

	cmd := new(ServiceReportCmd)

	orgs, err := cmd.getOrgs(conn)

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

	cmd := new(ServiceReportCmd)

	orgResource := new(orgResource)
	org := new(org)
	org.Name = "test-org"
	resourceMetadata := new(resourceMetadata)
	resourceMetadata.GUID = "abc123"
	orgResource.Entity = *org
	orgResource.Metadata = *resourceMetadata

	services, err := cmd.getServices(conn, *orgResource)

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

func readFile(filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	lines := strings.Split(string(content), "\n")
	return lines, err
}

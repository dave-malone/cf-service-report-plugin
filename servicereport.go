package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

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

type serviceBindings struct {
	*pagedResponse
	Resources []serviceBindingResource `json:"resources"`
}

func (s *serviceBindings) getServiceBindings() []serviceBinding {
	var serviceBindings []serviceBinding
	for _, serviceBindingsResource := range s.Resources {
		serviceBindings = append(serviceBindings, serviceBindingsResource.Entity)
	}
	return serviceBindings
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

type serviceBindingResource struct {
	Entity   serviceBinding   `json:"entity"`
	Metadata resourceMetadata `json:"metadata"`
}

type org struct {
	Name string `json:"name"`
}

type service struct {
	Label string `json:"label"`
}

type servicePlan struct {
	Name                string `json:"name"`
	ServiceInstancesURL string `json:"service_instances_url"`
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

//ServiceReportCmd the plugin
type ServiceReportCmd struct {
}

func main() {
	plugin.Start(new(ServiceReportCmd))
}

//GetMetadata returns metatada
func (cmd *ServiceReportCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "service-report",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "service-report",
				HelpText: "Services Report for all registered Service Brokers",
				UsageDetails: plugin.Usage{
					Usage: "cf service-report",
				},
			},
		},
	}
}

//Run runs the plugin
func (cmd *ServiceReportCmd) Run(cli plugin.CliConnection, args []string) {
	if nil == cli {
		fmt.Println("ERROR: CLI Connection is nil!")
		os.Exit(1)
	}

	fmt.Println("Gathering service information")

	if isLoggedIn, err := cli.IsLoggedIn(); err == nil && isLoggedIn != true {
		fmt.Println("You are not logged in. Please login using 'cf login' and try again")
		os.Exit(1)
	}

	if args[0] == "service-report" {
		cmd.printServiceUsageReport(cli)
	}
}

func (cmd *ServiceReportCmd) printServiceUsageReport(cli plugin.CliConnection) {
	servicePlans, err := cmd.getServicePlans(cli)

	if err != nil {
		fmt.Println("Failed to retreive service plans: " + err.Error())
		return
	}

	for _, servicePlanResource := range servicePlans.Resources {
		fmt.Printf("Service Plan: %s \n", servicePlanResource.Entity.Name)
	}

	orgs, err := cmd.getOrgs(cli)

	if err != nil {
		fmt.Println("Failed to retreive orgs: " + err.Error())
		return
	}

	for _, orgResource := range orgs.Resources {
		fmt.Printf("org: %s\n", orgResource.Entity.Name)

		services, err := cmd.getServices(cli, orgResource)
		if err != nil {
			fmt.Printf("Failed to retrieve services for org %s; %v", orgResource.Entity.Name, err.Error())
			return
		}

		for _, serviceResource := range services.Resources {
			fmt.Println("Service: " + serviceResource.Entity.Label)
		}
	}

}

func (cmd *ServiceReportCmd) getServicePlans(cli plugin.CliConnection) (servicePlans, error) {
	var servicePlans servicePlans

	data, err := cmd.cfcurl(cli, "/v2/service_plans")

	if nil != err {
		return servicePlans, err
	}

	err = json.Unmarshal(data, &servicePlans)

	if nil != err {
		fmt.Println("Failed to parse json: ", err.Error())
		return servicePlans, err
	}

	return servicePlans, err
}

func (cmd *ServiceReportCmd) getOrgs(cli plugin.CliConnection) (orgs, error) {
	var orgs orgs

	data, err := cmd.cfcurl(cli, "/v2/organizations?results-per-page=100")

	if nil != err {
		return orgs, err
	}

	err = json.Unmarshal(data, &orgs)

	if nil != err {
		fmt.Println("Failed to parse json: ", err.Error())
		return orgs, err
	}

	return orgs, err
}

func (cmd *ServiceReportCmd) getServices(cli plugin.CliConnection, orgResource orgResource) (services, error) {
	var services services

	data, err := cmd.cfcurl(cli, fmt.Sprintf("/v2/organizations/%s/services?results-per-page=100", orgResource.Metadata.GUID))

	if nil != err {
		return services, err
	}

	err = json.Unmarshal(data, &services)

	if nil != err {
		fmt.Println("Failed to parse json: ", err.Error())
		return services, err
	}

	return services, err
}

func (cmd *ServiceReportCmd) cfcurl(cli plugin.CliConnection, cliCommandArgs ...string) (data []byte, err error) {
	cliCommandArgs = append([]string{"curl"}, cliCommandArgs...)

	output, err := cli.CliCommandWithoutTerminalOutput(cliCommandArgs...)

	if nil != err {
		return nil, err
	}

	if nil == output || 0 == len(output) {
		return nil, errors.New("CF API returned no output")
	}

	response := strings.Join(output, " ")

	if 0 == len(response) || "" == response {
		return nil, errors.New("Failed to join output")
	}

	return []byte(response), err
}

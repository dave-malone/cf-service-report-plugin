package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

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

type serviceInstances struct {
	*pagedResponse
	Resources []serviceInstanceResource `json:"resources"`
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
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "Org\tService Instances\tBound App GUIDs")

	orgs, err := cmd.getOrgs(cli)

	if err != nil {
		fmt.Println("Failed to retreive orgs: " + err.Error())
		return
	}

	for _, orgResource := range orgs.Resources {
		var orgServiceInstances, boundApps string

		services, err := cmd.getServices(cli, orgResource)
		if err != nil {
			fmt.Printf("Failed to retrieve services for org %s; %v", orgResource.Entity.Name, err.Error())
			return
		}

		for _, serviceResource := range services.Resources {
			servicePlans, err := cmd.getServicePlans(cli, serviceResource.Entity)

			if err != nil {
				fmt.Println("Failed to retreive service plans: " + err.Error())
				return
			}

			for _, servicePlanResource := range servicePlans.Resources {
				serviceInstances, err := cmd.getServiceInstances(cli, servicePlanResource)

				if err != nil {
					fmt.Println("Failed to retreive service instances: " + err.Error())
					return
				}

				for _, serviceInstanceResource := range serviceInstances.Resources {
					if orgServiceInstances != "" {
						orgServiceInstances = fmt.Sprintf("%s, %s", orgServiceInstances, serviceInstanceResource.Entity.Name)
					} else {
						orgServiceInstances = serviceInstanceResource.Entity.Name
					}

					serviceBindings, err := cmd.getServiceBindings(cli, serviceInstanceResource.Entity)

					if err != nil {
						fmt.Println("Failed to retreive service bindings: " + err.Error())
						return
					}

					for _, serviceBindingResource := range serviceBindings.Resources {
						if boundApps != "" {
							boundApps = fmt.Sprintf("%s, %s", boundApps, serviceBindingResource.Entity.AppGUID)
						} else {
							boundApps = serviceBindingResource.Entity.AppGUID
						}
					}
				}
			}

		}

		fmt.Fprintf(w, "%v\t%v\t%v\n", orgResource.Entity.Name, orgServiceInstances, boundApps)
	}

	fmt.Fprintln(w)
	w.Flush()

}

func (cmd *ServiceReportCmd) getServicePlans(cli plugin.CliConnection, service service) (servicePlans, error) {
	var servicePlans servicePlans

	data, err := cmd.cfcurl(cli, service.ServicePlansURL)

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

func (cmd *ServiceReportCmd) getServiceInstances(cli plugin.CliConnection, servicePlanResource servicePlanResource) (serviceInstances, error) {
	var serviceInstances serviceInstances

	data, err := cmd.cfcurl(cli, servicePlanResource.Entity.ServiceInstancesURL)

	if nil != err {
		return serviceInstances, err
	}

	err = json.Unmarshal(data, &serviceInstances)

	if nil != err {
		fmt.Println("Failed to parse json: ", err.Error())
		return serviceInstances, err
	}

	return serviceInstances, err
}

func (cmd *ServiceReportCmd) getServiceBindings(cli plugin.CliConnection, serviceInstance serviceInstance) (serviceBindings, error) {
	var serviceBindings serviceBindings

	data, err := cmd.cfcurl(cli, serviceInstance.ServiceBindingsURL)

	if nil != err {
		return serviceBindings, err
	}

	err = json.Unmarshal(data, &serviceBindings)

	if nil != err {
		fmt.Println("Failed to parse json: ", err.Error())
		return serviceBindings, err
	}

	return serviceBindings, err
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

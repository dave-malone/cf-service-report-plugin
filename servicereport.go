package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

//ServiceInstances represents the structure of the GET:/v2/service_instances resposne
type ServiceInstances struct {
	TotalResults            int                      `json:"total_results"`
	TotalPages              int                      `json:"total_pages"`
	PrevUrl                 string                   `json:"prev_url"`
	NextUrl                 string                   `json:"next_url"`
	ServiceInstanceEntities []ServiceInstancesEntity `json:"resources"`
}

type ServiceInstancesEntity struct {
	Service Service `json:"entity"`
}

//Service represents the structure of the entity in the GET:/v2/service_instances response
type Service struct {
	Name               string `json:"name"`
	ServicePlanGuid    string `json:"service_plan_guid"`
	SpaceGuid          string `json:"space_guid"`
	SpaceUrl           string `json:"space_url"`
	ServiceBindingsUrl string `json:"service_bindings_url"`
	Type               string `json:"type"`
}

//ServiceReportCmd the plugin
type ServiceReportCmd struct {
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
	serviceInstances, err := cmd.getServices()
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("You are running %d total service instances.\n", serviceInstances.TotalResults)

	servicesReport := serviceInstances.getServicesReport()

	for key, value := range servicesReport {
		fmt.Println("Service:", key, "# Instances:", value)
	}
}

func (si *ServiceReportCmd) getServicesReport() map[string]int {
	m := make(map[string]int)

	for _, service := range si.getServices() {
		if i, ok := m[service.Name]; ok != true {
			m[service.Name] = 1
		} else {
			m[service.Name] = i + 1
		}
	}

	return m
}

func (si *ServiceInstances) getServices() []Service {
	var services []Service
	for _, serviceInstancesEntity := range si.ServiceInstanceEntities {
		services = append(services, serviceInstancesEntity.Service)
	}
	return services
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

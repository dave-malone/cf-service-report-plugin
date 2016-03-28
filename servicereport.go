package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/cloudfoundry/cli/plugin"
)

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

	if args[0] == "service-report" && len(args) == 1 {
		start := time.Now()
		syncReport(cli)
		end := time.Now()
		fmt.Printf("Execution Time: %v\n", end.Sub(start))
	} else if args[0] == "service-report" && args[1] == "a" {
		start := time.Now()
		asyncReport(cli)
		end := time.Now()
		fmt.Printf("Execution Time: %v\n", end.Sub(start))
	}
}

func syncReport(cli plugin.CliConnection) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "Org\tService Instances\tBound App GUIDs")

	orgs, err := getOrgs(cli)
	if err != nil {
		fmt.Println("Failed to retreive orgs: " + err.Error())
		return
	}

	for _, orgResource := range orgs.Resources {
		var orgServiceInstances, boundApps string

		services, err := getServices(cli, orgResource)
		if err != nil {
			fmt.Printf("Failed to retrieve services for org %s; %v", orgResource.Entity.Name, err.Error())
			return
		}

		for _, serviceResource := range services.Resources {
			servicePlans, err := getServicePlans(cli, serviceResource.Entity)

			if err != nil {
				fmt.Println("Failed to retreive service plans: " + err.Error())
				return
			}

			for _, servicePlanResource := range servicePlans.Resources {
				serviceInstances, err := getServiceInstances(cli, servicePlanResource)

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

					serviceBindings, err := getServiceBindings(cli, serviceInstanceResource.Entity)

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

func asyncReport(cli plugin.CliConnection) {
	orgChan := retrieveAndPublishOrgs(cli)
	serviceChan := retrieveAndPublishServices(cli, orgChan)
	servicePlanChan := retrieveAndPublishServicePlans(cli, serviceChan)
	serviceInstanceChan := retrieveAndPublishServiceInstances(cli, servicePlanChan)
	retrieveServiceBindings(cli, serviceInstanceChan)
}

func retrieveAndPublishOrgs(cli plugin.CliConnection) <-chan orgResource {
	orgChan := make(chan orgResource)

	go func() {
		orgs, err := getOrgs(cli)

		if err != nil {
			fmt.Println("Failed to retreive orgs: " + err.Error())
			return
		}

		for _, orgResource := range orgs.Resources {
			//TODO - how should this state be captured? Should it be passed around?
			//var orgServiceInstances, boundApps string
			fmt.Printf("Org: %s\n", orgResource.Entity.Name)
			orgChan <- orgResource
		}

		close(orgChan)
	}()

	return orgChan
}

func retrieveAndPublishServices(cli plugin.CliConnection, orgChan <-chan orgResource) <-chan serviceResource {
	serviceChan := make(chan serviceResource)

	go func() {
		for orgResource := range orgChan {
			services, err := getServices(cli, orgResource)

			if err != nil {
				fmt.Printf("Failed to retrieve services for org %s; %v", orgResource.Entity.Name, err.Error())
				continue
			}

			for _, serviceResource := range services.Resources {
				serviceChan <- serviceResource
			}
		}

		close(serviceChan)
	}()

	return serviceChan
}

func retrieveAndPublishServicePlans(cli plugin.CliConnection, serviceChan <-chan serviceResource) <-chan servicePlanResource {
	servicePlanChan := make(chan servicePlanResource)

	go func() {
		for serviceResource := range serviceChan {
			servicePlans, err := getServicePlans(cli, serviceResource.Entity)

			if err != nil {
				fmt.Printf("Failed to retreive service plans for service %s: %s\n", serviceResource.Entity.Label, err.Error())
				continue
			}

			for _, servicePlanResource := range servicePlans.Resources {
				servicePlanChan <- servicePlanResource
			}
		}

		close(servicePlanChan)
	}()

	return servicePlanChan
}

func retrieveAndPublishServiceInstances(cli plugin.CliConnection, servicePlanChan <-chan servicePlanResource) <-chan serviceInstanceResource {
	serviceInstanceChan := make(chan serviceInstanceResource)

	go func() {
		for servicePlanResource := range servicePlanChan {
			serviceInstances, err := getServiceInstances(cli, servicePlanResource)

			if err != nil {
				fmt.Printf("Failed to retreive service instances for service plan %s: %s\n", servicePlanResource.Entity.Name, err.Error())
				continue
			}

			for _, serviceInstanceResource := range serviceInstances.Resources {
				//TODO - how to build this, and where?
				// if orgServiceInstances != "" {
				// 	orgServiceInstances = fmt.Sprintf("%s, %s", orgServiceInstances, serviceInstanceResource.Entity.Name)
				// } else {
				// 	orgServiceInstances = serviceInstanceResource.Entity.Name
				// }

				fmt.Printf("Service Instance: %s\n", serviceInstanceResource.Entity.Name)

				serviceInstanceChan <- serviceInstanceResource
			}
		}

		close(serviceInstanceChan)
	}()

	return serviceInstanceChan
}

func retrieveServiceBindings(cli plugin.CliConnection, serviceInstanceChan <-chan serviceInstanceResource) {
	for serviceInstanceResource := range serviceInstanceChan {
		serviceBindings, err := getServiceBindings(cli, serviceInstanceResource.Entity)

		if err != nil {
			fmt.Printf("Failed to retreive service bindings for %s: %s\n", serviceInstanceResource.Entity.Name, err.Error())
			continue
		}

		for _, serviceBindingResource := range serviceBindings.Resources {
			//TODO - boundApps is within the context of an org and a service - how do we tie these back?
			// if boundApps != "" {
			// 	boundApps = fmt.Sprintf("%s, %s", boundApps, serviceBindingResource.Entity.AppGUID)
			// } else {
			// 	boundApps = serviceBindingResource.Entity.AppGUID
			// }

			fmt.Printf("Bound App GUID: %s\n", serviceBindingResource.Entity.AppGUID)
		}
	}

	//TODO - need to print results in a table
	//fmt.Fprintf(w, "%v\t%v\t%v\n", orgResource.Entity.Name, orgServiceInstances, boundApps)
}

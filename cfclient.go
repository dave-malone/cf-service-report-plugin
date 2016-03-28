package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry/cli/plugin"
)

func getServicePlans(cli plugin.CliConnection, service service) (servicePlans, error) {
	var servicePlans servicePlans

	data, err := cfcurl(cli, service.ServicePlansURL)

	if nil != err {
		return servicePlans, err
	}

	err = json.Unmarshal(data, &servicePlans)

	if nil != err {
		fmt.Printf("Failed to parse json response from %s; error: %s\n", service.ServicePlansURL, err.Error())
		return servicePlans, err
	}

	return servicePlans, err
}

func getOrgs(cli plugin.CliConnection) (orgs, error) {
	var orgs orgs

	data, err := cfcurl(cli, "/v2/organizations?results-per-page=100")

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

func getServices(cli plugin.CliConnection, orgResource orgResource) (services, error) {
	var services services

	data, err := cfcurl(cli, fmt.Sprintf("/v2/organizations/%s/services?results-per-page=100", orgResource.Metadata.GUID))

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

func getServiceInstances(cli plugin.CliConnection, servicePlanResource servicePlanResource) (serviceInstances, error) {
	var serviceInstances serviceInstances

	data, err := cfcurl(cli, servicePlanResource.Entity.ServiceInstancesURL)

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

func getServiceBindings(cli plugin.CliConnection, serviceInstance serviceInstance) (serviceBindings, error) {
	var serviceBindings serviceBindings

	data, err := cfcurl(cli, serviceInstance.ServiceBindingsURL)

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

func cfcurl(cli plugin.CliConnection, cliCommandArgs ...string) (data []byte, err error) {
	return httpGet(cli, cliCommandArgs[0])
	// cliCommandArgs = append([]string{"curl"}, cliCommandArgs...)
	//
	// output, err := cli.CliCommandWithoutTerminalOutput(cliCommandArgs...)
	//
	// if nil != err {
	// 	return nil, err
	// }
	//
	// if nil == output || 0 == len(output) {
	// 	return nil, errors.New("CF API returned no output")
	// }
	//
	// response := strings.Join(output, " ")
	//
	// if 0 == len(response) || "" == response {
	// 	return nil, errors.New("Failed to join output")
	// }
	//
	// return []byte(response), err
}

func httpGet(cli plugin.CliConnection, url string) ([]byte, error) {
	var payload []byte

	accessToken, err := cli.AccessToken()
	if err != nil {
		fmt.Printf("Failed to get access token: %s\n", err.Error())
		return payload, err
	}

	apiEndpoint, err := cli.ApiEndpoint()
	if err != nil {
		fmt.Printf("Failed to get Api Endpoint: %s\n", err.Error())
		return payload, err
	}

	apiURL := fmt.Sprintf("%s%s", apiEndpoint, url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Printf("http.NewRequest failed: %s\n", err.Error())
		return payload, err
	}

	req.Header.Add("Authorization", accessToken)

	response, err := client.Do(req)
	if err != nil {
		return payload, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %s\n", err.Error())
		return payload, err
	}

	return contents, nil
}

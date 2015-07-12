// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/krujos/usagereport-plugin/apihelper"
)

type FakeCFAPIHelper struct {
	GetOrgsStub        func(plugin.CliConnection) ([]apihelper.Organization, error)
	getOrgsMutex       sync.RWMutex
	getOrgsArgsForCall []struct {
		arg1 plugin.CliConnection
	}
	getOrgsReturns struct {
		result1 []apihelper.Organization
		result2 error
	}
	GetQuotaMemoryLimitStub        func(plugin.CliConnection, string) (float64, error)
	getQuotaMemoryLimitMutex       sync.RWMutex
	getQuotaMemoryLimitArgsForCall []struct {
		arg1 plugin.CliConnection
		arg2 string
	}
	getQuotaMemoryLimitReturns struct {
		result1 float64
		result2 error
	}
	GetOrgMemoryUsageStub        func(plugin.CliConnection, apihelper.Organization) (float64, error)
	getOrgMemoryUsageMutex       sync.RWMutex
	getOrgMemoryUsageArgsForCall []struct {
		arg1 plugin.CliConnection
		arg2 apihelper.Organization
	}
	getOrgMemoryUsageReturns struct {
		result1 float64
		result2 error
	}
	GetOrgSpacesStub        func(plugin.CliConnection, string) ([]apihelper.Space, error)
	getOrgSpacesMutex       sync.RWMutex
	getOrgSpacesArgsForCall []struct {
		arg1 plugin.CliConnection
		arg2 string
	}
	getOrgSpacesReturns struct {
		result1 []apihelper.Space
		result2 error
	}
	GetSpaceAppsStub        func(plugin.CliConnection, string) ([]apihelper.App, error)
	getSpaceAppsMutex       sync.RWMutex
	getSpaceAppsArgsForCall []struct {
		arg1 plugin.CliConnection
		arg2 string
	}
	getSpaceAppsReturns struct {
		result1 []apihelper.App
		result2 error
	}
}

func (fake *FakeCFAPIHelper) GetOrgs(arg1 plugin.CliConnection) ([]apihelper.Organization, error) {
	fake.getOrgsMutex.Lock()
	fake.getOrgsArgsForCall = append(fake.getOrgsArgsForCall, struct {
		arg1 plugin.CliConnection
	}{arg1})
	fake.getOrgsMutex.Unlock()
	if fake.GetOrgsStub != nil {
		return fake.GetOrgsStub(arg1)
	} else {
		return fake.getOrgsReturns.result1, fake.getOrgsReturns.result2
	}
}

func (fake *FakeCFAPIHelper) GetOrgsCallCount() int {
	fake.getOrgsMutex.RLock()
	defer fake.getOrgsMutex.RUnlock()
	return len(fake.getOrgsArgsForCall)
}

func (fake *FakeCFAPIHelper) GetOrgsArgsForCall(i int) plugin.CliConnection {
	fake.getOrgsMutex.RLock()
	defer fake.getOrgsMutex.RUnlock()
	return fake.getOrgsArgsForCall[i].arg1
}

func (fake *FakeCFAPIHelper) GetOrgsReturns(result1 []apihelper.Organization, result2 error) {
	fake.GetOrgsStub = nil
	fake.getOrgsReturns = struct {
		result1 []apihelper.Organization
		result2 error
	}{result1, result2}
}

func (fake *FakeCFAPIHelper) GetQuotaMemoryLimit(arg1 plugin.CliConnection, arg2 string) (float64, error) {
	fake.getQuotaMemoryLimitMutex.Lock()
	fake.getQuotaMemoryLimitArgsForCall = append(fake.getQuotaMemoryLimitArgsForCall, struct {
		arg1 plugin.CliConnection
		arg2 string
	}{arg1, arg2})
	fake.getQuotaMemoryLimitMutex.Unlock()
	if fake.GetQuotaMemoryLimitStub != nil {
		return fake.GetQuotaMemoryLimitStub(arg1, arg2)
	} else {
		return fake.getQuotaMemoryLimitReturns.result1, fake.getQuotaMemoryLimitReturns.result2
	}
}

func (fake *FakeCFAPIHelper) GetQuotaMemoryLimitCallCount() int {
	fake.getQuotaMemoryLimitMutex.RLock()
	defer fake.getQuotaMemoryLimitMutex.RUnlock()
	return len(fake.getQuotaMemoryLimitArgsForCall)
}

func (fake *FakeCFAPIHelper) GetQuotaMemoryLimitArgsForCall(i int) (plugin.CliConnection, string) {
	fake.getQuotaMemoryLimitMutex.RLock()
	defer fake.getQuotaMemoryLimitMutex.RUnlock()
	return fake.getQuotaMemoryLimitArgsForCall[i].arg1, fake.getQuotaMemoryLimitArgsForCall[i].arg2
}

func (fake *FakeCFAPIHelper) GetQuotaMemoryLimitReturns(result1 float64, result2 error) {
	fake.GetQuotaMemoryLimitStub = nil
	fake.getQuotaMemoryLimitReturns = struct {
		result1 float64
		result2 error
	}{result1, result2}
}

func (fake *FakeCFAPIHelper) GetOrgMemoryUsage(arg1 plugin.CliConnection, arg2 apihelper.Organization) (float64, error) {
	fake.getOrgMemoryUsageMutex.Lock()
	fake.getOrgMemoryUsageArgsForCall = append(fake.getOrgMemoryUsageArgsForCall, struct {
		arg1 plugin.CliConnection
		arg2 apihelper.Organization
	}{arg1, arg2})
	fake.getOrgMemoryUsageMutex.Unlock()
	if fake.GetOrgMemoryUsageStub != nil {
		return fake.GetOrgMemoryUsageStub(arg1, arg2)
	} else {
		return fake.getOrgMemoryUsageReturns.result1, fake.getOrgMemoryUsageReturns.result2
	}
}

func (fake *FakeCFAPIHelper) GetOrgMemoryUsageCallCount() int {
	fake.getOrgMemoryUsageMutex.RLock()
	defer fake.getOrgMemoryUsageMutex.RUnlock()
	return len(fake.getOrgMemoryUsageArgsForCall)
}

func (fake *FakeCFAPIHelper) GetOrgMemoryUsageArgsForCall(i int) (plugin.CliConnection, apihelper.Organization) {
	fake.getOrgMemoryUsageMutex.RLock()
	defer fake.getOrgMemoryUsageMutex.RUnlock()
	return fake.getOrgMemoryUsageArgsForCall[i].arg1, fake.getOrgMemoryUsageArgsForCall[i].arg2
}

func (fake *FakeCFAPIHelper) GetOrgMemoryUsageReturns(result1 float64, result2 error) {
	fake.GetOrgMemoryUsageStub = nil
	fake.getOrgMemoryUsageReturns = struct {
		result1 float64
		result2 error
	}{result1, result2}
}

func (fake *FakeCFAPIHelper) GetOrgSpaces(arg1 plugin.CliConnection, arg2 string) ([]apihelper.Space, error) {
	fake.getOrgSpacesMutex.Lock()
	fake.getOrgSpacesArgsForCall = append(fake.getOrgSpacesArgsForCall, struct {
		arg1 plugin.CliConnection
		arg2 string
	}{arg1, arg2})
	fake.getOrgSpacesMutex.Unlock()
	if fake.GetOrgSpacesStub != nil {
		return fake.GetOrgSpacesStub(arg1, arg2)
	} else {
		return fake.getOrgSpacesReturns.result1, fake.getOrgSpacesReturns.result2
	}
}

func (fake *FakeCFAPIHelper) GetOrgSpacesCallCount() int {
	fake.getOrgSpacesMutex.RLock()
	defer fake.getOrgSpacesMutex.RUnlock()
	return len(fake.getOrgSpacesArgsForCall)
}

func (fake *FakeCFAPIHelper) GetOrgSpacesArgsForCall(i int) (plugin.CliConnection, string) {
	fake.getOrgSpacesMutex.RLock()
	defer fake.getOrgSpacesMutex.RUnlock()
	return fake.getOrgSpacesArgsForCall[i].arg1, fake.getOrgSpacesArgsForCall[i].arg2
}

func (fake *FakeCFAPIHelper) GetOrgSpacesReturns(result1 []apihelper.Space, result2 error) {
	fake.GetOrgSpacesStub = nil
	fake.getOrgSpacesReturns = struct {
		result1 []apihelper.Space
		result2 error
	}{result1, result2}
}

func (fake *FakeCFAPIHelper) GetSpaceApps(arg1 plugin.CliConnection, arg2 string) ([]apihelper.App, error) {
	fake.getSpaceAppsMutex.Lock()
	fake.getSpaceAppsArgsForCall = append(fake.getSpaceAppsArgsForCall, struct {
		arg1 plugin.CliConnection
		arg2 string
	}{arg1, arg2})
	fake.getSpaceAppsMutex.Unlock()
	if fake.GetSpaceAppsStub != nil {
		return fake.GetSpaceAppsStub(arg1, arg2)
	} else {
		return fake.getSpaceAppsReturns.result1, fake.getSpaceAppsReturns.result2
	}
}

func (fake *FakeCFAPIHelper) GetSpaceAppsCallCount() int {
	fake.getSpaceAppsMutex.RLock()
	defer fake.getSpaceAppsMutex.RUnlock()
	return len(fake.getSpaceAppsArgsForCall)
}

func (fake *FakeCFAPIHelper) GetSpaceAppsArgsForCall(i int) (plugin.CliConnection, string) {
	fake.getSpaceAppsMutex.RLock()
	defer fake.getSpaceAppsMutex.RUnlock()
	return fake.getSpaceAppsArgsForCall[i].arg1, fake.getSpaceAppsArgsForCall[i].arg2
}

func (fake *FakeCFAPIHelper) GetSpaceAppsReturns(result1 []apihelper.App, result2 error) {
	fake.GetSpaceAppsStub = nil
	fake.getSpaceAppsReturns = struct {
		result1 []apihelper.App
		result2 error
	}{result1, result2}
}

var _ apihelper.CFAPIHelper = new(FakeCFAPIHelper)
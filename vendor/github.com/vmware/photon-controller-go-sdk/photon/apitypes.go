// Copyright (c) 2016 VMware, Inc. All Rights Reserved.
//
// This product is licensed to you under the Apache License, Version 2.0 (the "License").
// You may not use this product except in compliance with the License.
//
// This product may include a number of subcomponents with separate copyright notices and
// license terms. Your use of these subcomponents is subject to the terms and conditions
// of the subcomponent's license, as noted in the LICENSE file.

package photon

import (
	"fmt"
)

type Entity struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

// Implement a generic sdk error
type SdkError struct {
	Message string
}

func (e SdkError) Error() string {
	return fmt.Sprintf("photon: %v", e.Message)
}

// Represents an error from the Photon API.
type ApiError struct {
	Code           string                 `json:"code"`
	Data           map[string]interface{} `json:"data"`
	Message        string                 `json:"message"`
	HttpStatusCode int                    `json:"-"` // Not part of API contract
}

// Implement Go error interface for ApiError
func (e ApiError) Error() string {
	return fmt.Sprintf(
		"photon: { HTTP status: '%v', code: '%v', message: '%v', data: '%v' }",
		e.HttpStatusCode,
		e.Code,
		e.Message,
		e.Data)
}

// Used to represent a generic HTTP error, i.e. an unexpected HTTP 500.
type HttpError struct {
	StatusCode int
	Message    string
}

// Implementation of error interface for HttpError
func (e HttpError) Error() string {
	return fmt.Sprintf("photon: HTTP %d: %v", e.StatusCode, e.Message)
}

// Represents an Photon task that has entered into an error state.
// Photon task errors can be caught and type-checked against with
// the usual Go idiom.
type TaskError struct {
	ID   string `json:"id"`
	Step Step   `json:"step,omitempty"`
}

// Implement Go error interface for TaskError.
func (e TaskError) Error() string {
	return fmt.Sprintf("photon: Task '%s' is in error state: {@step==%s}", e.ID, GetStep(e.Step))
}

// An error representing a timeout while waiting for a task to complete.
type TaskTimeoutError struct {
	ID string
}

// Implement Go error interface for TaskTimeoutError.
func (e TaskTimeoutError) Error() string {
	return fmt.Sprintf("photon: Timed out waiting for task '%s'. "+
		"Task may not be in error state, examine task for full details.", e.ID)
}

// Represents an operation (Step) within a Task.
type Step struct {
	ID                 string                 `json:"id"`
	Operation          string                 `json:"operation,omitempty"`
	State              string                 `json:"state"`
	StartedTime        int64                  `json:"startedTime"`
	EndTime            int64                  `json:"endTime,omitempty"`
	QueuedTime         int64                  `json:"queuedTime"`
	Sequence           int                    `json:"sequence,omitempty"`
	ResourceProperties map[string]interface{} `json:"resourceProperties,omitempty"`
	Errors             []ApiError             `json:"errors,omitempty"`
	Warnings           []ApiError             `json:"warnings,omitempty"`
	Options            map[string]interface{} `json:"options,omitempty"`
	SelfLink           string                 `json:"selfLink,omitempty"`
}

// Implement Go error interface for Step.
func GetStep(s Step) string {
	return fmt.Sprintf("{\"sequence\"=>\"%d\",\"state\"=>\"%s\",\"errors\"=>%s,\"warnings\"=>%s,\"operation\"=>\"%s\","+
		"\"startedTime\"=>\"%d\",\"queuedTime\"=>\"%d\",\"endTime\"=>\"%d\",\"options\"=>%s}",
		s.Sequence, s.State, s.Errors, s.Warnings, s.Operation, s.StartedTime, s.QueuedTime,
		s.EndTime, s.Options)

}

// Represents an asynchronous task.
type Task struct {
	ID                 string      `json:"id"`
	Operation          string      `json:"operation,omitempty"`
	State              string      `json:"state"`
	StartedTime        int64       `json:"startedTime"`
	EndTime            int64       `json:"endTime,omitempty"`
	QueuedTime         int64       `json:"queuedTime"`
	Entity             Entity      `json:"entity,omitempty"`
	SelfLink           string      `json:"selfLink,omitempty"`
	Steps              []Step      `json:"steps,omitempty"`
	ResourceProperties interface{} `json:"resourceProperties,omitempty"`
}

// Represents multiple tasks returned by the API.
type TaskList struct {
	Items []Task `json:"items"`
}

// Options for GetTasks API.
type TaskGetOptions struct {
	State      string `urlParam:"state"`
	Kind       string `urlParam:"kind"`
	EntityID   string `urlParam:"entityId"`
	EntityKind string `urlParam:"entityKind"`
}

type BaseCompact struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type QuotaLineItem struct {
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
	Key   string  `json:"key"`
}

// The QuotaLineItem with limit and usage in one place.
type QuotaStatusLineItem struct {
	Unit  string  `json:"unit"`
	Limit float64 `json:"limit"`
	Usage float64 `json:"usage"`
}

// Creation spec for locality.
type LocalitySpec struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

// Creation spec for disks.
type DiskCreateSpec struct {
	Flavor     string         `json:"flavor"`
	Kind       string         `json:"kind"`
	CapacityGB int            `json:"capacityGb"`
	Affinities []LocalitySpec `json:"affinities,omitempty"`
	Name       string         `json:"name"`
	Tags       []string       `json:"tags,omitempty"`
}

// Represents a persistent disk.
type PersistentDisk struct {
	Flavor     string          `json:"flavor"`
	Cost       []QuotaLineItem `json:"cost"`
	Kind       string          `json:"kind"`
	Datastore  string          `json:"datastore,omitempty"`
	CapacityGB int             `json:"capacityGb,omitempty"`
	Name       string          `json:"name"`
	State      string          `json:"state"`
	ID         string          `json:"id"`
	VMs        []string        `json:"vms"`
	Tags       []string        `json:"tags,omitempty"`
	SelfLink   string          `json:"selfLink,omitempty"`
}

// Represents multiple persistent disks returned by the API.
type DiskList struct {
	Items []PersistentDisk `json:"items"`
}

// Creation spec for projects.
type ProjectCreateSpec struct {
	Name                       string   `json:"name"`
	SecurityGroups             []string `json:"securityGroups,omitempty"`
	DefaultRouterPrivateIpCidr string   `json:"defaultRouterPrivateIpCidr,omitempty"`
	ResourceQuota              Quota    `json:"quota,omitempty"`
}

// Represents multiple projects returned by the API.
type ProjectList struct {
	Items []ProjectCompact `json:"items"`
}

// Compact representation of projects.
type ProjectCompact struct {
	Kind           string          `json:"kind"`
	Name           string          `json:"name"`
	ID             string          `json:"id"`
	Tags           []string        `json:"tags"`
	SelfLink       string          `json:"selfLink"`
	SecurityGroups []SecurityGroup `json:"securityGroups"`
	ResourceQuota  Quota           `json:"quota,omitempty"`
}

// Represents an image.
type Image struct {
	Size                int64          `json:"size"`
	Kind                string         `json:"kind"`
	Name                string         `json:"name"`
	State               string         `json:"state"`
	ID                  string         `json:"id"`
	Tags                []string       `json:"tags"`
	Scope               ImageScope     `json:"scope"`
	SelfLink            string         `json:"selfLink"`
	Settings            []ImageSetting `json:"settings"`
	ReplicationType     string         `json:"replicationType"`
	ReplicationProgress string         `json:"replicationProgress"`
	SeedingProgress     string         `json:"seedingProgress"`
}

// Represents an image scope.
type ImageScope struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

// Represents an image setting.
type ImageSetting struct {
	Name         string `json:"name"`
	DefaultValue string `json:"defaultValue"`
}

// Creation spec for images.
type ImageCreateOptions struct {
	ReplicationType string
}

// Represents multiple images returned by the API.
type Images struct {
	Items []Image `json:"items"`
}

// Represents a component with status.
type Component struct {
	Component string
	Message   string
	Status    string
}

// Represents status of the photon system.
type Status struct {
	Status     string
	Components []Component
}

// Represents a single tenant.
type Tenant struct {
	Projects       []BaseCompact   `json:"projects"`
	Kind           string          `json:"kind"`
	Name           string          `json:"name"`
	ID             string          `json:"id"`
	SelfLink       string          `json:"selfLink"`
	Tags           []string        `json:"tags"`
	SecurityGroups []SecurityGroup `json:"securityGroups"`
	ResourceQuota  Quota           `json:"quota,omitempty"`
}

// Represents multiple tenants returned by the API.
type Tenants struct {
	Items []Tenant `json:"items"`
}

// Creation spec for tenants.
type TenantCreateSpec struct {
	Name           string   `json:"name"`
	SecurityGroups []string `json:"securityGroups,omitempty"`
	ResourceQuota  Quota    `json:"quota,omitempty"`
}

// Represents the quota
type Quota struct {
	QuotaLineItems map[string]QuotaStatusLineItem `json:"quotaItems"`
}

// QuotaSpec is used when set/update/excluding QuotaLineItems from existing Quota
type QuotaSpec map[string]QuotaStatusLineItem

// Creation spec for VMs.
type VmCreateSpec struct {
	Flavor        string            `json:"flavor"`
	SourceImageID string            `json:"sourceImageId"`
	AttachedDisks []AttachedDisk    `json:"attachedDisks"`
	Affinities    []LocalitySpec    `json:"affinities,omitempty"`
	Name          string            `json:"name"`
	Tags          []string          `json:"tags,omitempty"`
	Subnets       []string          `json:"subnets,omitempty"`
	Environment   map[string]string `json:"environment,omitempty"`
}

// Represents possible operations for VMs. Valid values include:
// START_VM, STOP_VM, RESTART_VM, SUSPEND_VM, RESUME_VM
type VmOperation struct {
	Operation string                 `json:"operation"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// Represents metadata that can be set on a VM.
type VmMetadata struct {
	Metadata map[string]string `json:"metadata"`
}

// Represents tags that can be set on a VM.
type VmTag struct {
	Tag string `json:"value"`
}

// Represents a single attached disk.
type AttachedDisk struct {
	Flavor     string `json:"flavor"`
	Kind       string `json:"kind"`
	CapacityGB int    `json:"capacityGb,omitempty"`
	Name       string `json:"name"`
	State      string `json:"state"`
	ID         string `json:"id,omitempty"`
	BootDisk   bool   `json:"bootDisk"`
}

// Represents a single VM.
type VM struct {
	SourceImageID string            `json:"sourceImageId,omitempty"`
	Cost          []QuotaLineItem   `json:"cost"`
	Kind          string            `json:"kind"`
	AttachedDisks []AttachedDisk    `json:"attachedDisks"`
	Datastore     string            `json:"datastore,omitempty"`
	AttachedISOs  []ISO             `json:"attachedIsos,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	SelfLink      string            `json:"selfLink,omitempty"`
	Flavor        string            `json:"flavor"`
	Host          string            `json:"host,omitempty"`
	Name          string            `json:"name"`
	State         string            `json:"state"`
	ID            string            `json:"id"`
	FloatingIp    string            `json:"floatingIp"`
}

// Represents multiple VMs returned by the API.
type VMs struct {
	Items []VM `json:"items"`
}

// Represents an ISO.
type ISO struct {
	Size int64  `json:"size,omitempty"`
	Kind string `json:"kind,omitempty"`
	Name string `json:"name"`
	ID   string `json:"id"`
}

// Represents operations for disks.
type VmDiskOperation struct {
	DiskID    string                 `json:"diskId"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// Represents a floating IP operation related to a VM.
type VmFloatingIpSpec struct {
	NetworkId string `json:"networkId"`
}

// Creation spec for flavors.
type FlavorCreateSpec struct {
	Cost []QuotaLineItem `json:"cost"`
	Kind string          `json:"kind"`
	Name string          `json:"name"`
}

// Represents a single flavor.
type Flavor struct {
	Cost     []QuotaLineItem `json:"cost"`
	Kind     string          `json:"kind"`
	Name     string          `json:"name"`
	ID       string          `json:"id"`
	Tags     []string        `json:"tags"`
	SelfLink string          `json:"selfLink"`
	State    string          `json:"state"`
}

// Represents multiple flavors returned by the API.
type FlavorList struct {
	Items []Flavor `json:"items"`
}

// Creation spec for hosts.
type HostCreateSpec struct {
	Username string            `json:"username"`
	Password string            `json:"password"`
	Zone     string            `json:"zone,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Address  string            `json:"address"`
	Tags     []string          `json:"usageTags"`
}

// Represents a host
type Host struct {
	Username   string            `json:"username"`
	Password   string            `json:"password"`
	Address    string            `json:"address"`
	Kind       string            `json:"kind"`
	ID         string            `json:"id"`
	Zone       string            `json:"zone,omitempty"`
	Tags       []string          `json:"usageTags"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	SelfLink   string            `json:"selfLink"`
	State      string            `json:"state"`
	EsxVersion string            `json:"esxVersion"`
}

// Represents multiple hosts returned by the API.
type Hosts struct {
	Items []Host `json:"items"`
}

type Datastore struct {
	Kind     string   `json:"kind"`
	Type     string   `json:"type"`
	Tags     []string `json:"tags,omitempty"`
	ID       string   `json:"id"`
	SelfLink string   `json:"selfLink"`
}

type Datastores struct {
	Items []Datastore `json:"items"`
}

type SystemUsage struct {
	NumberHosts      int `json:"numberHosts"`
	NumberVMs        int `json:"numberVMs"`
	NumberTenants    int `json:"numberTenants"`
	NumberProjects   int `json:"numberProjects"`
	NumberDatastores int `json:"numberDatastores"`
	NumberServices   int `json:"numberServices"`
}

// Represents stats information
type StatsInfo struct {
	Enabled       bool   `json:"enabled,omitempty"`
	StoreEndpoint string `json:"storeEndpoint,omitempty"`
	StorePort     int    `json:"storePort,omitempty"`
}

// Represents authentication information
type AuthInfo struct {
	Password       string   `json:"password,omitempty"`
	Endpoint       string   `json:"endpoint,omitempty"`
	Domain         string   `json:"domain,omitempty"`
	Port           int      `json:"port,omitempty"`
	SecurityGroups []string `json:"securityGroups,omitempty"`
	Username       string   `json:"username,omitempty"`
}

// Represents ip range
type IpRange struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

// Represents creation spec for network configuration.
type NetworkConfigurationCreateSpec struct {
	Enabled         bool     `json:"sdnEnabled,omitempty"`
	Address         string   `json:"networkManagerAddress,omitempty"`
	Username        string   `json:"networkManagerUsername,omitempty"`
	Password        string   `json:"networkManagerPassword,omitempty"`
	NetworkZoneId   string   `json:"networkZoneId,omitempty"`
	TopRouterId     string   `json:"networkTopRouterId,omitempty"`
	EdgeIpPoolId    string   `json:"networkEdgeIpPoolId,omitempty"`
	HostUplinkPnic  string   `json:"networkHostUplinkPnic,omitempty"`
	IpRange         string   `json:"ipRange,omitempty"`
	ExternalIpRange *IpRange `json:"externalIpRange,omitempty"`
}

// Represents network configuration.
type NetworkConfiguration struct {
	Enabled         bool     `json:"sdnEnabled,omitempty"`
	Address         string   `json:"networkManagerAddress,omitempty"`
	Username        string   `json:"networkManagerUsername,omitempty"`
	Password        string   `json:"networkManagerPassword,omitempty"`
	NetworkZoneId   string   `json:"networkZoneId,omitempty"`
	TopRouterId     string   `json:"networkTopRouterId,omitempty"`
	EdgeIpPoolId    string   `json:"networkEdgeIpPoolId,omitempty"`
	HostUplinkPnic  string   `json:"networkHostUplinkPnic,omitempty"`
	IpRange         string   `json:"ipRange,omitempty"`
	FloatingIpRange *IpRange `json:"floatingIpRange,omitempty"`
	SnatIp          string   `json:"snatIp,omitempty"`
}

// Represents a router
type Router struct {
	ID            string `json:"id"`
	Kind          string `json:"kind"`
	Name          string `json:"name"`
	PrivateIpCidr string `json:"privateIpCidr"`
	IsDefault     bool   `json:"isDefault"`
}

// Represents multiple routers returned by the API.
type Routers struct {
	Items []Router `json:"items"`
}

type RouterCreateSpec struct {
	Name          string `json:"name"`
	PrivateIpCidr string `json:"privateIpCidr"`
}

// Represents name that can be set for router
type RouterUpdateSpec struct {
	RouterName string `json:"name"`
}

// Creation spec for Service Configuration.
type ServiceConfigurationSpec struct {
	Type    string `json:"type"`
	ImageID string `json:"imageId"`
}

// Represnts a Service configuration.
type ServiceConfiguration struct {
	Kind    string `json:"kind"`
	Type    string `json:"type"`
	ImageID string `json:"imageId"`
}

// Creation spec for services.
type ServiceCreateSpec struct {
	Name               string            `json:"name"`
	Type               string            `json:"type"`
	VMFlavor           string            `json:"vmFlavor,omitempty"`
	MasterVmFlavor     string            `json:"masterVmFlavor,omitempty"`
	WorkerVmFlavor     string            `json:"workerVmFlavor,omitempty"`
	DiskFlavor         string            `json:"diskFlavor,omitempty"`
	SubnetId           string            `json:"subnetId,omitempty"`
	ImageID            string            `json:"imageId,omitempty"`
	WorkerCount        int               `json:"workerCount"`
	BatchSizeWorker    int               `json:"workerBatchExpansionSize,omitempty"`
	ExtendedProperties map[string]string `json:"extendedProperties"`
}

// Represents a service.
type Service struct {
	Kind               string                `json:"kind"`
	Name               string                `json:"name"`
	State              string                `json:"state"`
	ID                 string                `json:"id"`
	Type               string                `json:"type"`
	ImageID            string                `json:"imageId"`
	UpgradeStatus      *ServiceUpgradeStatus `json:"upgradeStatus,omitempty"`
	ProjectID          string                `json:"projectID,omitempty"`
	ClientID           string                `json:"clientId,omitempty"`
	WorkerCount        int                   `json:"workerCount"`
	SelfLink           string                `json:"selfLink,omitempty"`
	ErrorReason        string                `json:"errorReason,omitempty"`
	ExtendedProperties map[string]string     `json:"extendedProperties"`
}

// Represents the status of a service during upgrade.
type ServiceUpgradeStatus struct {
	NewImageID           string `json:"newImageId"`
	UpgradeMessage       string `json:"upgradeMessage,omitempty"`
	TotalNodes           int    `json:"totalNodes,omitempty"`
	NumNodesUpgraded     int    `json:"numNodesUpgraded"`
	UpgradeResultMessage string `json:"upgradeResultMessage,omitempty"`
}

// Represents multiple services returned by the API
type Services struct {
	Items []Service `json:"items"`
}

// Represents service size that can be resized for service
type ServiceResizeOperation struct {
	NewWorkerCount int `json:"newWorkerCount"`
}

// Represents service imageId that can be updated during change version
type ServiceChangeVersionOperation struct {
	NewImageID string `json:"newImageId"`
}

// Represents a security group
type SecurityGroup struct {
	Name      string `json:"name"`
	Inherited bool   `json:"inherited"`
}

// Represents set_security_groups spec
type SecurityGroupsSpec struct {
	Items []string `json:"items"`
}

// Represents availability zone that can be set for host
type HostSetAvailabilityZoneOperation struct {
	AvailabilityZoneId string `json:"availabilityZoneId"`
}

// Represents single zone.
type Zone struct {
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	State    string `json:"state"`
	ID       string `json:"id"`
	SelfLink string `json:"selfLink"`
}

// Represents multiple zones returned by the API.
type Zones struct {
	Items []Zone `json:"items"`
}

// Creation spec for zones.
type ZoneCreateSpec struct {
	Name string `json:"name"`
}

// Represents the list of image datastores.
type ImageDatastores struct {
	Items []string `json:"items"`
}

// Image creation spec.
type ImageCreateSpec struct {
	Name            string `json:"name"`
	ReplicationType string `json:"replicationType"`
}

// Represents deployment info
type Info struct {
	BaseVersion   string `json:"baseVersion"`
	FullVersion   string `json:"fullVersion"`
	GitCommitHash string `json:"gitCommitHash"`
	NetworkType   string `json:"networkType"`
}

// NSX configuration spec
type NsxConfigurationSpec struct {
	NsxAddress             string            `json:"nsxAddress"`
	NsxUsername            string            `json:"nsxUsername"`
	NsxPassword            string            `json:"nsxPassword"`
	DhcpServerAddresses    map[string]string `json:"dhcpServerAddresses"`
	FloatingIpRootRange    IpRange           `json:"floatingIpRootRange"`
	T0RouterId             string            `json:"t0RouterId"`
	EdgeClusterId          string            `json:"edgeClusterId"`
	OverlayTransportZoneId string            `json:"overlayTransportZoneId"`
	TunnelIpPoolId         string            `json:"tunnelIpPoolId"`
	HostUplinkPnic         string            `json:"hostUplinkPnic"`
	HostUplinkVlanId       int               `json:"hostUplinkVlanId"`
	DnsServerAddresses     []string          `json:"dnsServerAddresses"`
}

// Represents port groups.
type PortGroups struct {
	Names []string `json:"names"`
}

// Represents a subnet
type Subnet struct {
	ID                 string            `json:"id"`
	Kind               string            `json:"kind"`
	Name               string            `json:"name"`
	Description        string            `json:"description,omitempty"`
	PrivateIpCidr      string            `json:"privateIpCidr"`
	ReservedIps        map[string]string `json:"reservedIps"`
	State              string            `json:"state"`
	IsDefault          bool              `json:"isDefault"`
	PortGroups         PortGroups        `json:"portGroups"`
	DnsServerAddresses []string          `json:"dnsServerAddresses"`
}

// Represents multiple subnets returned by the API.
type Subnets struct {
	Items []Subnet `json:"items"`
}

// Creation spec for subnets.
type SubnetCreateSpec struct {
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	PrivateIpCidr      string     `json:"privateIpCidr"`
	Type               string     `json:"type"`
	PortGroups         PortGroups `json:"portGroups"`
	DnsServerAddresses []string   `json:"dnsServerAddresses"`
}

// Represents name that can be set for subnet
type SubnetUpdateSpec struct {
	SubnetName string `json:"name"`
}

// Identity and Access Management (IAM)
// IAM Policy entry
type PolicyEntry struct {
	Principal string   `json:"principal"`
	Roles     []string `json:"roles"`
}

type PolicyDelta struct {
	Principal string `json:"principal"`
	Action    string `json:"action"`
	Role      string `json:"role"`
}

type SystemInfo struct {
	NTPEndpoint             string                 `json:"ntpEndpoint,omitempty"`
	UseImageDatastoreForVms bool                   `json:"useImageDatastoreForVms,omitempty"`
	Auth                    *AuthInfo              `json:"auth"`
	NetworkConfiguration    *NetworkConfiguration  `json:"networkConfiguration"`
	Kind                    string                 `json:"kind"`
	SyslogEndpoint          string                 `json:"syslogEndpoint,omitempty"`
	Stats                   *StatsInfo             `json:"stats,omitempty"`
	State                   string                 `json:"state"`
	ImageDatastores         []string               `json:"imageDatastores"`
	ServiceConfigurations   []ServiceConfiguration `json:"serviceConfigurations,omitempty"`
	LoadBalancerEnabled     bool                   `json:"loadBalancerEnabled"`
	LoadBalancerAddress     string                 `json:"loadBalancerAddress"`
	BaseVersion             string                 `json:"baseVersion"`
	FullVersion             string                 `json:"fullVersion"`
	GitCommitHash           string                 `json:"gitCommitHash"`
	NetworkType             string                 `json:"networkType"`
}

package mvdata

// Region is a location the platform provisions into. Status is what makes one usable — a region
// the platform has announced but cannot yet provision into is served here and refused at create.
type Region struct {
	Code        string `json:"code"`
	DisplayName string `json:"displayName"`
	Status      string `json:"status"`
	IsDefault   bool   `json:"isDefault"`
}

// VPC represents a virtual private cloud.
type VPC struct {
	Name      string  `json:"name"`
	Region    *Region `json:"region,omitempty"`
	CreatedAt string  `json:"createdAt,omitempty"`
}

// createVPCRequest is the request body for POST /vpcs. The API names a region by code on the way
// in and returns the whole region on the way out, so one field cannot carry both directions.
type createVPCRequest struct {
	Name   string `json:"name"`
	Region string `json:"region,omitempty"`
}

// Subnet represents a subnet within a VPC.
type Subnet struct {
	Name      string  `json:"name"`
	VPCName   string  `json:"vpcName"`
	CIDRBlock string  `json:"cidrBlock"`
	Region    *Region `json:"region,omitempty"`
	CreatedAt string  `json:"createdAt,omitempty"`
}

// createSubnetRequest is the request body for POST /subnets. A subnet takes the region of the VPC
// it is created in, so the caller does not choose one.
type createSubnetRequest struct {
	Name      string `json:"name"`
	VPCName   string `json:"vpcName"`
	CIDRBlock string `json:"cidrBlock"`
}

// Instance represents a compute instance.
type Instance struct {
	Name              string  `json:"name"`
	VPCName           string  `json:"vpcName,omitempty"`
	InstanceType      string  `json:"instanceType"`
	AuthorizedKeyName string  `json:"authorizedKeyName"`
	PrivateIP         string  `json:"privateIp,omitempty"`
	Status            string  `json:"status,omitempty"`
	HourlyPrice       float64 `json:"hourlyPrice,omitempty"`
	Region            *Region `json:"region,omitempty"`
	CreatedAt         string  `json:"createdAt,omitempty"`
}

// createInstanceRequest is the request body for POST /instances.
type createInstanceRequest struct {
	Name              string `json:"name"`
	VPCName           string `json:"vpcName"`
	InstanceType      string `json:"instanceType"`
	AuthorizedKeyName string `json:"authorizedKeyName"`
}

// Key represents an SSH public key registered with the console.
type Key struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// KubernetesCluster represents a managed Kubernetes cluster.
type KubernetesCluster struct {
	Name             string  `json:"name"`
	Version          string  `json:"version"`
	NodeInstanceType string  `json:"nodeInstanceType"`
	NodeCount        int     `json:"nodeCount"`
	Endpoint         string  `json:"endpoint,omitempty"`
	Status           string  `json:"status,omitempty"`
	Region           *Region `json:"region,omitempty"`
	CreatedAt        string  `json:"createdAt,omitempty"`
}

// createKubernetesClusterRequest is the request body for POST /kubernetes. Region is a code here
// for the same reason it is on a VPC create, and the fields the platform fills in — endpoint,
// status, createdAt — are not sent.
type createKubernetesClusterRequest struct {
	Name             string `json:"name"`
	Version          string `json:"version"`
	NodeInstanceType string `json:"nodeInstanceType"`
	NodeCount        int    `json:"nodeCount"`
	Region           string `json:"region,omitempty"`
}

// KubernetesClusterUpdate holds mutable fields for PATCH.
type KubernetesClusterUpdate struct {
	NodeCount int `json:"nodeCount"`
}

// InstanceType represents an available instance type with pricing.
type InstanceType struct {
	InstanceType string  `json:"instanceType"`
	HourlyPrice  float64 `json:"hourlyPrice"`
}

// APIKey represents a console API key.
type APIKey struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Key        string `json:"key,omitempty"`
	Prefix     string `json:"prefix"`
	ExpiresAt  string `json:"expiresAt,omitempty"`
	LastUsedAt string `json:"lastUsedAt,omitempty"`
	CreatedAt  string `json:"createdAt,omitempty"`
}

// APIKeyCreate is the request body for creating an API key.
type APIKeyCreate struct {
	Name      string `json:"name"`
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// DeviceAuth is the response from POST /auth/device/authorize.
type DeviceAuth struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// DeviceTokenResponse is the response from POST /auth/device/token.
type DeviceTokenResponse struct {
	Status        string `json:"status"`
	APIToken      string `json:"api_token,omitempty"`
	AccountName   string `json:"account_name,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	Email         string `json:"email,omitempty"`
}

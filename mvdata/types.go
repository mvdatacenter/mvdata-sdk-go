package mvdata

// VPC represents a virtual private cloud.
type VPC struct {
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// Subnet represents a subnet within a VPC.
type Subnet struct {
	Name      string `json:"name"`
	VPCName   string `json:"vpcName"`
	CIDRBlock string `json:"cidrBlock"`
	CreatedAt string `json:"createdAt,omitempty"`
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
	Name             string `json:"name"`
	Version          string `json:"version"`
	NodeInstanceType string `json:"nodeInstanceType"`
	NodeCount        int    `json:"nodeCount"`
	Endpoint         string `json:"endpoint,omitempty"`
	Status           string `json:"status,omitempty"`
	CreatedAt        string `json:"createdAt,omitempty"`
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

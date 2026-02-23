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
	InstanceType      string  `json:"instanceType"`
	AuthorizedKeyName string  `json:"authorizedKeyName"`
	PrivateIP         string  `json:"privateIp,omitempty"`
	Status            string  `json:"status,omitempty"`
	HourlyPrice       float64 `json:"hourlyPrice,omitempty"`
	CreatedAt         string  `json:"createdAt,omitempty"`
}

// createInstanceRequest is the request body for POST /vpcs/:vpcName/instances.
type createInstanceRequest struct {
	Name              string `json:"name"`
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

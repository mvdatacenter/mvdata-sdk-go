package mvdata

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(handler http.Handler) *Client {
	server := httptest.NewServer(handler)
	return &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}
}

// --- VPC ---

func TestCreateVPC(t *testing.T) {
	var gotMethod, gotPath, gotAuth string
	var gotBody VPC

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(VPC{Name: "production", CreatedAt: "2025-01-01T00:00:00Z"})
	}))

	result, err := client.CreateVPC(context.Background(), &VPC{Name: "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/vpcs" {
		t.Errorf("expected /vpcs, got %s", gotPath)
	}
	if gotAuth != "Bearer test-token" {
		t.Errorf("expected Bearer test-token, got %s", gotAuth)
	}
	if gotBody.Name != "production" {
		t.Errorf("expected name production, got %s", gotBody.Name)
	}
	if result.Name != "production" {
		t.Errorf("expected result name production, got %s", result.Name)
	}
}

func TestGetVPC(t *testing.T) {
	var gotMethod, gotPath string

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewEncoder(w).Encode(VPC{Name: "production"})
	}))

	result, err := client.GetVPC(context.Background(), "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/vpcs/production" {
		t.Errorf("expected /vpcs/production, got %s", gotPath)
	}
	if result.Name != "production" {
		t.Errorf("expected name production, got %s", result.Name)
	}
}

func TestGetVPC_NotFound(t *testing.T) {
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Not Found"}`))
	}))

	_, err := client.GetVPC(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
	var nfe *NotFoundError
	if !errors.As(err, &nfe) {
		t.Errorf("expected NotFoundError, got %T: %v", err, err)
	}
}

func TestGetVPC_ServerError(t *testing.T) {
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal Server Error"}`))
	}))

	_, err := client.GetVPC(context.Background(), "broken")
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
	var nfe *NotFoundError
	if errors.As(err, &nfe) {
		t.Errorf("expected non-NotFoundError for 500, but got NotFoundError: %v", err)
	}
}

func TestDeleteVPC(t *testing.T) {
	var gotMethod, gotPath string

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))

	if err := client.DeleteVPC(context.Background(), "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", gotMethod)
	}
	if gotPath != "/vpcs/production" {
		t.Errorf("expected /vpcs/production, got %s", gotPath)
	}
}

// --- Subnet ---

func TestCreateSubnet(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody Subnet

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Subnet{Name: "public", VPCName: "production", CIDRBlock: "10.0.1.0/24"})
	}))

	result, err := client.CreateSubnet(context.Background(), &Subnet{Name: "public", VPCName: "production", CIDRBlock: "10.0.1.0/24"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/subnets" {
		t.Errorf("expected /subnets, got %s", gotPath)
	}
	if gotBody.VPCName != "production" {
		t.Errorf("expected vpcName production, got %s", gotBody.VPCName)
	}
	if gotBody.CIDRBlock != "10.0.1.0/24" {
		t.Errorf("expected cidrBlock 10.0.1.0/24, got %s", gotBody.CIDRBlock)
	}
	if result.Name != "public" {
		t.Errorf("expected result name public, got %s", result.Name)
	}
}

func TestGetSubnet(t *testing.T) {
	var gotMethod, gotPath string

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewEncoder(w).Encode(Subnet{Name: "public", VPCName: "production", CIDRBlock: "10.0.1.0/24"})
	}))

	result, err := client.GetSubnet(context.Background(), "public")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/subnets/public" {
		t.Errorf("expected /subnets/public, got %s", gotPath)
	}
	if result.VPCName != "production" {
		t.Errorf("expected vpcName production, got %s", result.VPCName)
	}
}

func TestDeleteSubnet(t *testing.T) {
	var gotMethod, gotPath string

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))

	if err := client.DeleteSubnet(context.Background(), "public"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", gotMethod)
	}
	if gotPath != "/subnets/public" {
		t.Errorf("expected /subnets/public, got %s", gotPath)
	}
}

// --- Instance ---

func TestCreateInstance(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody Instance

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Instance{
			Name: "web-01", VPCName: "production", SubnetName: "public",
			InstanceType: "c1", AuthorizedKeyName: "deploy-key",
			PrivateIP: "10.0.0.2", Status: "running", HourlyPrice: 0.50,
		})
	}))

	result, err := client.CreateInstance(context.Background(), &Instance{
		Name: "web-01", VPCName: "production", SubnetName: "public",
		InstanceType: "c1", AuthorizedKeyName: "deploy-key",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/instances" {
		t.Errorf("expected /instances, got %s", gotPath)
	}
	if gotBody.VPCName != "production" {
		t.Errorf("expected vpcName production, got %s", gotBody.VPCName)
	}
	if gotBody.SubnetName != "public" {
		t.Errorf("expected subnetName public, got %s", gotBody.SubnetName)
	}
	if result.PrivateIP != "10.0.0.2" {
		t.Errorf("expected privateIp 10.0.0.2, got %s", result.PrivateIP)
	}
}

func TestGetInstance(t *testing.T) {
	var gotMethod, gotPath string

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewEncoder(w).Encode(Instance{
			Name: "web-01", VPCName: "production", SubnetName: "public",
			InstanceType: "c1", AuthorizedKeyName: "deploy-key",
			PrivateIP: "10.0.0.2", Status: "running",
		})
	}))

	result, err := client.GetInstance(context.Background(), "web-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/instances/web-01" {
		t.Errorf("expected /instances/web-01, got %s", gotPath)
	}
	if result.VPCName != "production" {
		t.Errorf("expected vpcName production, got %s", result.VPCName)
	}
}

func TestDeleteInstance(t *testing.T) {
	var gotMethod, gotPath string

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))

	if err := client.DeleteInstance(context.Background(), "web-01"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", gotMethod)
	}
	if gotPath != "/instances/web-01" {
		t.Errorf("expected /instances/web-01, got %s", gotPath)
	}
}

// --- Key ---

func TestCreateKey(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody Key

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Key{Name: "deploy-key", Key: "ssh-ed25519 AAAA..."})
	}))

	result, err := client.CreateKey(context.Background(), &Key{Name: "deploy-key", Key: "ssh-ed25519 AAAA..."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/keys" {
		t.Errorf("expected /keys, got %s", gotPath)
	}
	if gotBody.Name != "deploy-key" {
		t.Errorf("expected name deploy-key, got %s", gotBody.Name)
	}
	if result.Name != "deploy-key" {
		t.Errorf("expected result name deploy-key, got %s", result.Name)
	}
}

func TestGetKey(t *testing.T) {
	var gotMethod, gotPath, gotQuery string

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		json.NewEncoder(w).Encode(Key{Name: "deploy-key", Key: "ssh-ed25519 AAAA..."})
	}))

	result, err := client.GetKey(context.Background(), "deploy-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/keys" {
		t.Errorf("expected /keys, got %s", gotPath)
	}
	if gotQuery != "name=deploy-key" {
		t.Errorf("expected query name=deploy-key, got %s", gotQuery)
	}
	if result.Name != "deploy-key" {
		t.Errorf("expected name deploy-key, got %s", result.Name)
	}
}

func TestGetKey_NotFound(t *testing.T) {
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Not Found"}`))
	}))

	_, err := client.GetKey(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
	var nfe *NotFoundError
	if !errors.As(err, &nfe) {
		t.Errorf("expected NotFoundError, got %T: %v", err, err)
	}
}

func TestDeleteKey(t *testing.T) {
	var gotMethod, gotPath string

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))

	if err := client.DeleteKey(context.Background(), "deploy-key"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", gotMethod)
	}
	if gotPath != "/keys/deploy-key" {
		t.Errorf("expected /keys/deploy-key, got %s", gotPath)
	}
}

// --- Kubernetes ---

func TestCreateKubernetesCluster(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody KubernetesCluster

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(KubernetesCluster{
			Name: "k8s-main", Version: "1.29", NodeInstanceType: "c1", NodeCount: 3,
			Endpoint: "https://k8s-main.k8s.mvdatacenter.com", Status: "provisioning",
		})
	}))

	result, err := client.CreateKubernetesCluster(context.Background(), &KubernetesCluster{
		Name: "k8s-main", Version: "1.29", NodeInstanceType: "c1", NodeCount: 3,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/kubernetes" {
		t.Errorf("expected /kubernetes, got %s", gotPath)
	}
	if gotBody.NodeCount != 3 {
		t.Errorf("expected nodeCount 3, got %d", gotBody.NodeCount)
	}
	if result.Endpoint != "https://k8s-main.k8s.mvdatacenter.com" {
		t.Errorf("expected endpoint https://k8s-main.k8s.mvdatacenter.com, got %s", result.Endpoint)
	}
}

func TestGetKubernetesCluster(t *testing.T) {
	var gotMethod, gotPath string

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewEncoder(w).Encode(KubernetesCluster{
			Name: "k8s-main", Version: "1.29", NodeInstanceType: "c1", NodeCount: 3,
			Endpoint: "https://k8s-main.k8s.mvdatacenter.com", Status: "running",
		})
	}))

	result, err := client.GetKubernetesCluster(context.Background(), "k8s-main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/kubernetes/k8s-main" {
		t.Errorf("expected /kubernetes/k8s-main, got %s", gotPath)
	}
	if result.NodeCount != 3 {
		t.Errorf("expected nodeCount 3, got %d", result.NodeCount)
	}
}

func TestUpdateKubernetesCluster(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody KubernetesClusterUpdate

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		json.NewEncoder(w).Encode(KubernetesCluster{
			Name: "k8s-main", Version: "1.29", NodeInstanceType: "c1", NodeCount: 5,
			Endpoint: "https://k8s-main.k8s.mvdatacenter.com", Status: "running",
		})
	}))

	result, err := client.UpdateKubernetesCluster(context.Background(), "k8s-main", &KubernetesClusterUpdate{NodeCount: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch {
		t.Errorf("expected PATCH, got %s", gotMethod)
	}
	if gotPath != "/kubernetes/k8s-main" {
		t.Errorf("expected /kubernetes/k8s-main, got %s", gotPath)
	}
	if gotBody.NodeCount != 5 {
		t.Errorf("expected nodeCount 5, got %d", gotBody.NodeCount)
	}
	if result.NodeCount != 5 {
		t.Errorf("expected result nodeCount 5, got %d", result.NodeCount)
	}
}

func TestDeleteKubernetesCluster(t *testing.T) {
	var gotMethod, gotPath string

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))

	if err := client.DeleteKubernetesCluster(context.Background(), "k8s-main"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", gotMethod)
	}
	if gotPath != "/kubernetes/k8s-main" {
		t.Errorf("expected /kubernetes/k8s-main, got %s", gotPath)
	}
}

// --- Instance Types ---

func TestListInstanceTypes(t *testing.T) {
	var gotMethod, gotPath string

	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewEncoder(w).Encode([]InstanceType{
			{InstanceType: "c1", HourlyPrice: 0.50},
			{InstanceType: "c1_g1", HourlyPrice: 1.25},
		})
	}))

	result, err := client.ListInstanceTypes(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/prices" {
		t.Errorf("expected /prices, got %s", gotPath)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 instance types, got %d", len(result))
	}
	if result[0].InstanceType != "c1" {
		t.Errorf("expected first type c1, got %s", result[0].InstanceType)
	}
}

// --- Error cases ---

func TestCreateVPC_ServerError(t *testing.T) {
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal Server Error"}`))
	}))

	_, err := client.CreateVPC(context.Background(), &VPC{Name: "test"})
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}

func TestCreateInstance_ServerError(t *testing.T) {
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal Server Error"}`))
	}))

	_, err := client.CreateInstance(context.Background(), &Instance{Name: "test", InstanceType: "c1", AuthorizedKeyName: "key"})
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}

// --- New constructor ---

func TestNew(t *testing.T) {
	client := New("https://api.mvdatacenter.com", "my-token")

	if client.BaseURL != "https://api.mvdatacenter.com" {
		t.Errorf("expected BaseURL https://api.mvdatacenter.com, got %s", client.BaseURL)
	}
	if client.Token != "my-token" {
		t.Errorf("expected Token my-token, got %s", client.Token)
	}
	if client.HTTPClient == nil {
		t.Error("expected non-nil HTTPClient")
	}
}

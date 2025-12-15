package storage

import (
	"testing"
	"time"
)

// TestService_Validation tests Service struct validation
func TestService_Validation(t *testing.T) {
	tests := []struct {
		name    string
		service Service
		wantErr bool
	}{
		{
			name: "Valid service",
			service: Service{
				Name:      "GitHub",
				Secret:    "JBSWY3DPEHPK3PXP",
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Valid service with LastUsed",
			service: Service{
				Name:      "AWS",
				Secret:    "JBSWY3DPEHPK3PXP",
				CreatedAt: time.Now(),
				LastUsed:  timePtr(time.Now()),
			},
			wantErr: false,
		},
		{
			name: "Valid service with identifier",
			service: Service{
				Name:       "GitHub",
				Identifier: "user@example.com",
				Secret:     "JBSWY3DPEHPK3PXP",
				CreatedAt:  time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			service: Service{
				Name:      "",
				Secret:    "JBSWY3DPEHPK3PXP",
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Name too long",
			service: Service{
				Name:      "This is a very long service name that exceeds the maximum allowed length of fifty characters",
				Secret:    "JBSWY3DPEHPK3PXP",
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Empty secret",
			service: Service{
				Name:      "GitHub",
				Secret:    "",
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Invalid secret format",
			service: Service{
				Name:      "GitHub",
				Secret:    "INVALID!@#$",
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Name with control characters",
			service: Service{
				Name:      "Git\nHub",
				Secret:    "JBSWY3DPEHPK3PXP",
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Name with path separator",
			service: Service{
				Name:      "Git/Hub",
				Secret:    "JBSWY3DPEHPK3PXP",
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateServiceName tests service name validation
func TestValidateServiceName(t *testing.T) {
	tests := []struct {
		name    string
		svcName string
		wantErr bool
	}{
		{
			name:    "Valid name",
			svcName: "GitHub",
			wantErr: false,
		},
		{
			name:    "Valid name with spaces",
			svcName: "My Service",
			wantErr: false,
		},
		{
			name:    "Valid name with hyphens",
			svcName: "my-service",
			wantErr: false,
		},
		{
			name:    "Valid name with underscores",
			svcName: "my_service",
			wantErr: false,
		},
		{
			name:    "Empty name",
			svcName: "",
			wantErr: true,
		},
		{
			name:    "Only whitespace",
			svcName: "   ",
			wantErr: true,
		},
		{
			name:    "Too long",
			svcName: "This is a very long service name that exceeds fifty chars",
			wantErr: true,
		},
		{
			name:    "Contains newline",
			svcName: "Git\nHub",
			wantErr: true,
		},
		{
			name:    "Contains tab",
			svcName: "Git\tHub",
			wantErr: true,
		},
		{
			name:    "Contains path separator",
			svcName: "path/to/service",
			wantErr: true,
		},
		{
			name:    "Contains backslash",
			svcName: "path\\to\\service",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateServiceName(tt.svcName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateServiceName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestStorage_AddService tests adding services to storage
func TestStorage_AddService(t *testing.T) {
	storage := &Storage{
		Version:  1,
		Services: []Service{},
	}

	// Add first service
	service1 := Service{
		Name:      "GitHub",
		Secret:    "JBSWY3DPEHPK3PXP",
		CreatedAt: time.Now(),
	}

	err := storage.AddService(service1)
	if err != nil {
		t.Fatalf("AddService() error = %v", err)
	}

	if len(storage.Services) != 1 {
		t.Errorf("Storage.Services length = %d, want 1", len(storage.Services))
	}

	// Try to add duplicate (should fail)
	err = storage.AddService(service1)
	if err == nil {
		t.Error("AddService() expected error for duplicate service, got nil")
	}

	// Add second service (should succeed)
	service2 := Service{
		Name:      "AWS",
		Secret:    "JBSWY3DPEHPK3PXP",
		CreatedAt: time.Now(),
	}

	err = storage.AddService(service2)
	if err != nil {
		t.Fatalf("AddService() error = %v", err)
	}

	if len(storage.Services) != 2 {
		t.Errorf("Storage.Services length = %d, want 2", len(storage.Services))
	}
}

// TestStorage_GetService tests retrieving services
func TestStorage_GetService(t *testing.T) {
	storage := &Storage{
		Version: 1,
		Services: []Service{
			{
				Name:      "GitHub",
				Secret:    "JBSWY3DPEHPK3PXP",
				CreatedAt: time.Now(),
			},
			{
				Name:      "AWS",
				Secret:    "JBSWY3DPEHPK3PXP",
				CreatedAt: time.Now(),
			},
		},
	}

	// Test existing service
	service, err := storage.GetService("GitHub")
	if err != nil {
		t.Fatalf("GetService() error = %v", err)
	}
	if service.Name != "GitHub" {
		t.Errorf("GetService() name = %s, want GitHub", service.Name)
	}

	// Test case-insensitive lookup
	service, err = storage.GetService("github")
	if err != nil {
		t.Fatalf("GetService() case-insensitive error = %v", err)
	}
	if service.Name != "GitHub" {
		t.Errorf("GetService() case-insensitive name = %s, want GitHub", service.Name)
	}

	// Test non-existent service
	_, err = storage.GetService("NonExistent")
	if err == nil {
		t.Error("GetService() expected error for non-existent service, got nil")
	}
}

// TestStorage_UpdateLastUsed tests updating last used timestamp
func TestStorage_UpdateLastUsed(t *testing.T) {
	now := time.Now()
	storage := &Storage{
		Version: 1,
		Services: []Service{
			{
				Name:      "GitHub",
				Secret:    "JBSWY3DPEHPK3PXP",
				CreatedAt: now,
			},
		},
	}

	// Verify LastUsed is nil initially
	if storage.Services[0].LastUsed != nil {
		t.Error("Service.LastUsed should be nil initially")
	}

	// Update last used
	err := storage.UpdateLastUsed("GitHub")
	if err != nil {
		t.Fatalf("UpdateLastUsed() error = %v", err)
	}

	// Verify LastUsed is set
	if storage.Services[0].LastUsed == nil {
		t.Fatal("Service.LastUsed should be set after UpdateLastUsed")
	}

	// Verify timestamp is recent (within 1 second)
	if time.Since(*storage.Services[0].LastUsed) > time.Second {
		t.Error("Service.LastUsed timestamp is not recent")
	}
}

// Helper function to create time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}

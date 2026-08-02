package demo

import (
	"context"
	"log/slog"
	"testing"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
)

func TestSeed_AllExampleEmails(t *testing.T) {
	fixtures := Seed()
	for _, u := range fixtures.Users {
		if u.Email == "" || !endsWithExample(u.Email) {
			t.Errorf("user email %s should end with .example", u.Email)
		}
	}
	for _, c := range fixtures.Clients {
		if c.ContactEmail == "" || !endsWithExample(c.ContactEmail) {
			t.Errorf("client email %s should end with .example", c.ContactEmail)
		}
	}
}

func TestSeed_AllDemoPrefixed(t *testing.T) {
	fixtures := Seed()
	for _, u := range fixtures.Users {
		if !IsDemoRecord(string(u.ID)) {
			t.Errorf("user ID %s should be demo-prefixed", u.ID)
		}
	}
	for _, c := range fixtures.Clients {
		if !IsDemoRecord(string(c.ID)) {
			t.Errorf("client ID %s should be demo-prefixed", c.ID)
		}
	}
	for _, q := range fixtures.Quotations {
		if !IsDemoRecord(string(q.ID)) {
			t.Errorf("quotation ID %s should be demo-prefixed", q.ID)
		}
	}
}

func TestIsDemoRecord(t *testing.T) {
	if !IsDemoRecord("demo-001") {
		t.Error("demo-001 should be a demo record")
	}
	if IsDemoRecord("real-001") {
		t.Error("real-001 should not be a demo record")
	}
	if IsDemoRecord("") {
		t.Error("empty string should not be a demo record")
	}
}

// Acceptance: "Demonstration reset removes only synthetic demonstration records."
func TestReset_RemovesOnlyDemoRecords(t *testing.T) {
	svc := NewResetService(slog.Default())
	ctx := context.Background()

	// Simulate a store with demo and non-demo records.
	users := map[string]bool{
		"demo-admin-001": true,
		"demo-user-002":  true,
		"real-user-001":  true,
		"real-user-002":  true,
	}
	clients := map[string]bool{
		"demo-client-acme":  true,
		"demo-client-stark": true,
		"real-client-001":   true,
	}

	removeUsers := func(ctx context.Context, pred func(id string) bool) int {
		count := 0
		for id := range users {
			if pred(id) {
				delete(users, id)
				count++
			}
		}
		return count
	}
	removeClients := func(ctx context.Context, pred func(id string) bool) int {
		count := 0
		for id := range clients {
			if pred(id) {
				delete(clients, id)
				count++
			}
		}
		return count
	}
	removeQuotations := func(ctx context.Context, pred func(id string) bool) int { return 0 }
	removeDocuments := func(ctx context.Context, pred func(id string) bool) int { return 0 }

	result, err := svc.Reset(ctx, domain.UserID("admin"), removeUsers, removeClients, removeQuotations, removeDocuments)
	if err != nil {
		t.Fatalf("Reset: %v", err)
	}

	if result.UsersRemoved != 2 {
		t.Errorf("users removed = %d, want 2", result.UsersRemoved)
	}
	if result.ClientsRemoved != 2 {
		t.Errorf("clients removed = %d, want 2", result.ClientsRemoved)
	}

	// Non-demo records must survive.
	if !users["real-user-001"] {
		t.Error("real-user-001 was removed. Reset must not touch non-demo records")
	}
	if !users["real-user-002"] {
		t.Error("real-user-002 was removed. Reset must not touch non-demo records")
	}
	if !clients["real-client-001"] {
		t.Error("real-client-001 was removed. Reset must not touch non-demo records")
	}

	// Demo records must be gone.
	if users["demo-admin-001"] {
		t.Error("demo-admin-001 should have been removed")
	}
	if clients["demo-client-acme"] {
		t.Error("demo-client-acme should have been removed")
	}
}

func TestEnvironmentMarker(t *testing.T) {
	if EnvironmentMarker("demo") == "" {
		t.Error("demo env should have a marker")
	}
	if EnvironmentMarker("production") != "" {
		t.Error("production env should not have a demo marker")
	}
}

func endsWithExample(email string) bool {
	return len(email) > 8 && email[len(email)-8:] == ".example"
}
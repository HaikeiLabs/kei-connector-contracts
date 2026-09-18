package connectors

import (
	"reflect"
	"testing"
)

func TestLookupCapability(t *testing.T) {
	cap, ok := LookupCapability(ProviderCRM, "lead.read")
	if !ok {
		t.Fatal("expected lead.read to be defined for crm")
	}
	if cap.Name != "lead.read" || cap.Action != ActionRead {
		t.Fatalf("unexpected capability: %+v", cap)
	}

	if _, ok := LookupCapability(ProviderCRM, "admin.raw_sql"); ok {
		t.Fatal("undefined capability must not be found")
	}
	if _, ok := LookupCapability("unknown", "lead.read"); ok {
		t.Fatal("unknown provider must not yield capabilities")
	}
}

func TestCapabilityForErrors(t *testing.T) {
	if _, err := CapabilityFor(ProviderGitHub, "repository.read"); err != nil {
		t.Fatalf("defined capability errored: %v", err)
	}
	if _, err := CapabilityFor(ProviderGitHub, "repository.purge"); err == nil {
		t.Fatal("undefined capability must error")
	}
}

func TestCapabilitiesForProviderCopiesAndRejectsUnknown(t *testing.T) {
	got, err := CapabilitiesForProvider(ProviderGoogle)
	if err != nil {
		t.Fatalf("CapabilitiesForProvider(google_drive): %v", err)
	}
	if len(got) == 0 {
		t.Fatal("google_drive must declare capabilities")
	}
	if _, err := CapabilitiesForProvider("atlassian"); err == nil {
		t.Fatal("unknown provider must error, not return an empty set")
	}
}

func TestProvidersSortedAndComplete(t *testing.T) {
	got := Providers()
	want := []Provider{ProviderCRM, ProviderGitHub, ProviderGoogle, ProviderHTTPAPI, ProviderLinear, ProviderNotion, ProviderS3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Providers() = %v, want %v", got, want)
	}
}

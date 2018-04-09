package main

import "testing"

func TestHandleBootstrapErrorCommand(t *testing.T) {
	dns_out := []byte("nodename nor servname provided")
	auth_out := []byte("Authentication failed")
	timeout_out := []byte("ConnectionTimeout")
	auth_host := Host{
		Hostname: "testhost_auth",
		Domain:   "testdomain",
		ChefEnv:  "testenv",
		RunList:  "testrecipe",
	}
	dns_host := Host{
		Hostname: "testhost_dns",
		Domain:   "testdomain",
		ChefEnv:  "testenv",
		RunList:  "testrecipe",
	}
	timeout_host := Host{
		Hostname: "testhost_timeout",
		Domain:   "testdomain",
		ChefEnv:  "testenv",
		RunList:  "testrecipe",
	}

	dns_result := handle_bootstrap_error(dns_out, dns_host, 1)
	if dns_result != true {
		t.Error("Expected dns error true, got ", dns_result)
	}
	auth_result := handle_bootstrap_error(auth_out, auth_host, 1)
	if auth_result != true {
		t.Error("Expected auth error true, got ", auth_result)
	}
	timeout_result := handle_bootstrap_error(timeout_out, timeout_host, 1)
	if timeout_result != true {
		t.Error("Expected timeout error true, got ", timeout_result)
	}

}

package main

import "testing"

func TestGenerateCommand(t *testing.T) {
	var test_host Host
	test_host = Host{
		Hostname: "testhost",
		Domain:   "testdomain",
		ChefEnv:  "testenv",
		RunList:  "testrecipe",
	}
	result := generate_command(test_host)
	if result != "knife bootstrap testhost.testdomain -N testhost -E testenv --sudo --ssh-user testuser1 --ssh-password testuser1 -r testrecipe" {
		t.Error("Expected knife bootstrap testhost.testdomain -N testhost -E testenv --sudo --ssh-user testuser1 --ssh-password testuser1 -r testrecipe, got ", result)
	}
}

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
	// TODO: not currently testing the env variables for ssh user/pw, needs to be mocked.
	if result != "knife bootstrap testhost.testdomain -N testhost -E testenv --sudo --ssh-user  --ssh-password  -r testrecipe" {
		t.Error("Expected knife bootstrap testhost.testdomain -N testhost -E testenv --sudo --ssh-user  --ssh-password  -r testrecipe, got ", result)
	}
}

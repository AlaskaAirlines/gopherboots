# gopherboots

#Introduction 
TODO: Give a short introduction of your project. Let this section explain the objectives or the motivation behind this project. 

#Getting Started
TODO: Guide users through getting your code up and running on their own system. In this section you can talk about:
1.	Installation process
2.	Software dependencies
3.	Latest releases
4.	API references

#Usage
It's expected that you already have a working command line knife installation, and are able to run commands like `knife bootstrap` successfully.
You'll need to set a few environment variables specific to your organization:
`SUPERUSER_NAME` should be set to a superuser account valid on the target host.
`SUPERUSER_PW` should be set to a superuser password valid on the target host.

A tsv file containing all the hosts you'd like to bootstrap should be formatted like this:
```
hostname	domain	target chef environment	runlist
```

For example:
```
test-host	example.org	linux	chef-client,base
```

To bootstrap all hosts simply run the following (where "./hosts.tsv" is the location of your tsv file):
```
./gopherboots -file=./hosts.tsv
```

If you encounter errors on any hosts during the bootstrapping process you can view them in the `./logs` directory.

#Build and Test
TODO: Describe and show how to build your code and run the tests. 

#Contribute
TODO: Explain how other users and developers can contribute to make your code better.
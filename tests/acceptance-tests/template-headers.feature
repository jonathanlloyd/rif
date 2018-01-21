@needs_echo_server
Feature: Template Headers

	Scenario: The user makes a request with a URL template from a .rif file
		Given a .rif file is on disk that has a header template
		When the user runs RIF on that file passing in the appropriate variables
		Then RIF should return an echo of the request it made

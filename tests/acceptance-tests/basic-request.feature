@needs_echo_server
Feature: Basic Requests

	Scenario: The user makes a basic GET request from a .rif file
		Given a .rif file is on disk that describes a GET request
		When the user runs RIF on that file
		Then RIF should return an echo of the request it made

	Scenario: The user makes a basic request with headers from a .rif file
		Given a .rif file is on disk that describes a request with headers
		When the user runs RIF on that file
		Then RIF should return an echo of the request it made

	Scenario: The user makes a basic request with a body from a .rif file
		Given a .rif file is on disk that describes a request with a body
		When the user runs RIF on that file
		Then RIF should return an echo of the request it made

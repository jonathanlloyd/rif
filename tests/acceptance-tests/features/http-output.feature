@needs_echo_server
Feature: HTTP output format

	Scenario: The user makes a basic GET request with a HTTP output format
		Given a .rif file is on disk that describes a GET request
		When the user runs RIF on that file with a HTTP output format
		Then RIF should return the HTTP/1.x representation of the request/response

	Scenario: The user makes a request with an unknown output format
		Given a .rif file is on disk that describes a GET request
		When the user runs RIF on that file with an unknown output format
		Then RIF should error

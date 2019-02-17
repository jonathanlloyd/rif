@needs_echo_server
Feature: Output formats

	Scenario: The user makes a basic GET request with a HTTP output format
		Given a .rif file is on disk that describes a GET request
		When the user runs RIF on that file with a HTTP output format
		Then RIF should return the HTTP/1.x representation of the request/response

	Scenario: The user makes a request with a body with a HTTP output format
		Given a .rif file is on disk that describes a request with a body
		When the user runs RIF on that file with a HTTP output format
		Then RIF should return the HTTP/1.x representation of the request/response

	Scenario: The user makes a request with a URL template and a HTTP output format
		Given a .rif file is on disk that has a URL template
		When the user runs RIF on that file passing in the HTTP output flag and variables
		Then RIF should return the HTTP/1.x representation of the request/response

	Scenario: The user makes a request with the cURL output format
		Given a .rif file is on disk that describes a request with a body
		When the user runs RIF on that file with the cURL output format
		Then RIF should return a cURL command equivalent to the request

	Scenario: The user makes a request with a URL template and a cURL output format
		Given a .rif file is on disk that has a URL template
		When the user runs RIF on that file passing in the cURL output flag and variables
		Then RIF should return a cURL command equivalent to the request

	Scenario: The user makes a request with an unknown output format
		Given a .rif file is on disk that describes a GET request
		When the user runs RIF on that file with an unknown output format
		Then RIF should error

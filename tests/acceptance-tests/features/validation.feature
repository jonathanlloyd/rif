Feature: .rif file validation

	Scenario: The user makes a request without a rif file version
		Given a .rif file is on disk that has no rif file version
		When the user runs RIF on that file
		Then RIF should error

	Scenario: The user makes a request with a higher rif file version
		Given a .rif file is on disk that has a higher rif file version
		When the user runs RIF on that file
		Then RIF should error

	Scenario: The user makes a request without passing in the required variables
		Given a .rif file is on disk that has some required variables
		When the user runs RIF on that file without passing in those variables
		Then RIF should error

Feature: .rif file validation

	Scenario: The user makes a request with a higher rif file version
		Given a .rif file is on disk that has a higher rif file version
		When the user runs RIF on that file
		Then RIF should error

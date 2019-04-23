import pytest

from common import *


def test_no_rif_file_version():
    """RIF files must specify a version"""
    output, return_code = run_rif(
        [make_test_file_path('no-version.rif')],
    )

    expected_output = """
Invalid .rif file: 
 - rif_version is required
 """[1:-1]

    assert return_code == 1
    assert expected_output in output

def test_version_too_high():
    """The RIF file version must not be higher than the current major"""
    output, return_code = run_rif(
        [make_test_file_path('bad-version.rif')],
    )

    expected_output = """
Invalid .rif file: 
 - rif_version must not be greater than the maximum supported version (0)
 """[1:-1]

    assert return_code == 1
    assert expected_output in output

def test_url_missing():
    """RIF files must specify a URL"""
    output, return_code = run_rif(
        [make_test_file_path('missing-url.rif')],
    )

    expected_output = """
Invalid .rif file: 
 - Field "URL" is required
 """[1:-1]

    assert return_code == 1
    assert expected_output in output

def test_method_missing():
    """RIF files must specify a method"""
    output, return_code = run_rif(
        [make_test_file_path('missing-method.rif')],
    )

    expected_output = """
Invalid .rif file: 
 - Field "method" is required
 """[1:-1]

    assert return_code == 1
    assert expected_output in output

def test_invalid_method():
    """RIF files must contain a valid HTTP method"""
    output, return_code = run_rif(
        [make_test_file_path('bad-method.rif')],
    )

    expected_output = """
Invalid .rif file: 
 - Method "NOTAVALIDMETHOD" is invalid
 """[1:-1]

    assert return_code == 1
    assert expected_output in output

def test_invalid_variable_type():
    """RIF files must contain valid variable types"""
    output, return_code = run_rif(
        [make_test_file_path('bad-variable-type.rif')],
    )

    expected_output = """
Invalid .rif file: 
 - Variable "foo" is invalid: variable has invalid type "notavalidtype" (valid types are boolean, number and string)
 """[1:-1]

    assert return_code == 1
    assert expected_output in output

def test_missing_required_variables():
    """RIF files must provide values for all required variables"""
    output, return_code = run_rif(
        [make_test_file_path('required-variables.rif')],
    )

    expected_output = """
Invalid parameters: 
Missing required variable(s): VAR_A

The variables for this RIF file are as follows:
Required:
 - VAR_A ( number )
Optional:
 - VAR_B ( string, default=value )
"""[1:-1]

    assert return_code == 1
    assert expected_output in output


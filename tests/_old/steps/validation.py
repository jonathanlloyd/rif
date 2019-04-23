"""Step definitions for validation features"""
from hamcrest import *

from behave import *
from common import *


@given(u'a .rif file is on disk that has no rif file version')
def step_impl(context):
    context.filename = test_file_path('no-version.rif')
    context.expected_error_msg = """
Invalid .rif file: 
 - rif_version is required
 """[1:-1]

@given(u'a .rif file is on disk that has a higher rif file version')
def step_impl(context):
    context.filename = test_file_path('bad-version.rif')
    context.expected_error_msg = """
Invalid .rif file: 
 - rif_version must not be greater than the maximum supported version (0)
 """[1:-1]

@given(u'a .rif file is on disk without a URL')
def step_impl(context):
    context.filename = test_file_path('missing-url.rif')
    context.expected_error_msg = """
Invalid .rif file: 
 - Field "URL" is required
 """[1:-1]

@given(u'a .rif file is on disk without a method')
def step_impl(context):
    context.filename = test_file_path('missing-method.rif')
    context.expected_error_msg = """
Invalid .rif file: 
 - Field "method" is required
 """[1:-1]

@given(u'a .rif file is on disk with an invalid method')
def step_impl(context):
    context.filename = test_file_path('bad-method.rif')
    context.expected_error_msg = """
Invalid .rif file: 
 - Method "NOTAVALIDMETHOD" is invalid
 """[1:-1]

@given(u'a .rif file is on disk that has an invalid variable type')
def step_impl(context):
    context.filename = test_file_path('bad-variable-type.rif')
    context.expected_error_msg = """
Invalid .rif file: 
 - Variable "foo" is invalid: variable has invalid type "notavalidtype" (valid types are boolean, number and string)
 """[1:-1]

@given(u'a .rif file is on disk that has some required variables')
def step_impl(context):
    context.filename = test_file_path('required-variables.rif')
    context.expected_error_msg = """
Invalid parameters: 
Missing required variable(s): VAR_A

The variables for this RIF file are as follows:
Required:
 - VAR_A ( number )
Optional:
 - VAR_B ( string, default=value )
"""[1:-1]

@when(u'the user runs RIF on that file without passing in those variables')
def step_impl(context):
    context.stdout, context.returncode = run_rif([
        context.filename,
    ])

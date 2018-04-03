"""Step definitions for validation features"""
from hamcrest import *

from behave import *
from common import *


@given(u'a .rif file is on disk that has no rif file version')
def step_impl(context):
    context.filename = '/vol/tests/test-data/no-version.rif'
    context.expected_error_msg = """
Invalid RIF file: 
 - rif_version is required
 """[1:-1]

@given(u'a .rif file is on disk that has a higher rif file version')
def step_impl(context):
    context.filename = '/vol/tests/test-data/bad-version.rif'
    context.expected_error_msg = "Error parsing .rif file: rif file version " \
            "greater than maxium supported version - 0"

@given(u'a .rif file is on disk that has some required variables')
def step_impl(context):
    context.filename = '/vol/tests/test-data/required-variables.rif'
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

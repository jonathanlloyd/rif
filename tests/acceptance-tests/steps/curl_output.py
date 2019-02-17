"""Step definitions for HTTP output"""
from hamcrest import *

from behave import *
from common import *

@when(u'the user runs RIF on that file with the cURL output format')
def step_impl(context):
    context.stdout, context.returncode = run_rif([
        context.filename,
        '--output=curl',
    ])

@when(u'the user runs RIF on that file passing in the cURL output flag and variables')
def step_impl(context):
    variable_args = [
        '{}={}'.format(name, value)
        for name, value in context.variables.items()
    ]
    context.stdout, context.returncode = run_rif([
        context.filename,
        *variable_args,
        '--output=curl',
    ])

@then(u'RIF should return a cURL command equivalent to the request')
def step_impl(context):
    assert_that(
        context.returncode,
        equal_to(0),
    )

    # These are invisible anyway and mess up the test string
    stdout_without_carriage_returns = context.stdout.replace("\r", "")

    assert_that(
        stdout_without_carriage_returns,
        contains_string(context.expected_curl_output),
    )

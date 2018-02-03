"""Step definitions for HTTP output"""
from hamcrest import *

from behave import *
from common import *


@when(u'the user runs RIF on that file with a HTTP output format')
def step_impl(context):
    context.stdout, context.returncode = run_rif([
        context.filename,
        '--output=http',
    ])


@when(u'the user runs RIF on that file with an unknown output format')
def step_impl(context):
    context.expected_error_msg = "Unknown output format: NOT_AN_OUTPUT_FORMAT"
    context.stdout, context.returncode = run_rif([
        context.filename,
        '--output=NOT_AN_OUTPUT_FORMAT',
    ])


@then(u'RIF should return the HTTP/1.x representation of the request/response')
def step_impl(context):
    assert_that(
        context.returncode,
        equal_to(0),
    )

    # These are invisible anyway and mess up the test string
    stdout_without_carriage_returns = context.stdout.replace("\r", "")

    assert_that(
        stdout_without_carriage_returns,
        contains_string(context.expected_http_output),
    )

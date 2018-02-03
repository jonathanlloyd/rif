"""Step definitions for the plain output format"""
from hamcrest import *

from behave import *
from common import *


@when(u'the user runs RIF on that file')
def step_impl(context):
    context.stdout, context.returncode = run_rif([context.filename])

@then(u'RIF should return an echo of the request it made')
def step_impl(context):
    assert_that(
        context.returncode,
        equal_to(0),
    )
    assert_that(
        context.stdout,
        contains_string(context.expected_plain_output),
    )


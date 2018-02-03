"""Step definitions for validation features"""
from hamcrest import *

from behave import *
from common import *


@given(u'a .rif file is on disk that has a higher rif file version')
def step_impl(context):
    context.filename = '/vol/tests/test-data/bad-version.rif'
    context.expected_error_msg = "Error parsing .rif file: rif file version " \
            "greater than maxium supported version - 0"


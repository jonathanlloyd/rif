"""Step definitions for acceptance tests"""
import subprocess

from behave import *
from hamcrest import *


@then(u'RIF should error')
def step_impl(context):
    assert_that(
        context.returncode,
        is_not(equal_to(0)),
    )
    assert_that(
        context.stdout,
        contains_string(context.expected_error_msg),
    )

def run_rif(args):
    result = subprocess.run(
        ['./build/rif'] + args,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
    )
    return_code = result.returncode
    output = result.stdout.decode('utf8')

    if return_code != 0:
        print(output)

    return (output, return_code)


def test_file_path(filename):
    return './tests/acceptance-tests/test-data/' + filename

"""Step definitions for acceptance tests"""
import subprocess

from hamcrest import *


@given(u'a .rif file is on disk that describes a GET request')
def step_impl(context):
    context.filename = '/vol/tests/test-data/basic-get.rif'
    context.expected_result = """
GET /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.2.0"""[1:]

@given(u'a .rif file is on disk that describes a request with headers')
def step_impl(context):
    context.filename = '/vol/tests/test-data/basic-headers.rif'
    context.expected_result = """
GET /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.2.0
x-test-header: some_value"""[1:]

@given(u'a .rif file is on disk that describes a request with a body')
def step_impl(context):
    context.filename = '/vol/tests/test-data/basic-body.rif'
    context.expected_result = """
POST /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
content-length: 4
user-agent: RIF/0.2.0

test"""[1:]

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
        contains_string(context.expected_result),
    )

@given(u'a .rif file is on disk that has a URL template')
def step_impl(context):
    context.filename = '/vol/tests/test-data/url-params.rif'
    context.variables = {
        'NUM_THINGS': 20,
    }
    context.expected_result = """
GET /url-params?count=20 HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.2.0"""[1:]

@given(u'a .rif file is on disk that has a URL template with a default')
def step_impl(context):
    context.filename = '/vol/tests/test-data/url-params.rif'
    context.variables = {}
    context.expected_result = """
GET /url-params?count=10 HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.2.0"""[1:]

@given(u'a .rif file is on disk that has a header template')
def step_impl(context):
    context.filename = '/vol/tests/test-data/header-params.rif'
    context.variables = {
        'HEADER_VALUE': 'header-value',
    }
    context.expected_result = """
GET /header-params HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.2.0
x-custom-header: header-value"""[1:]

@given(u'a .rif file is on disk that has a body template')
def step_impl(context):
    context.filename = '/vol/tests/test-data/body-params.rif'
    context.variables = {
        'BODY_VALUE': 'body-value',
    }
    context.expected_result = """
POST /body-params HTTP/1.1
host: localhost:8080
accept-encoding: gzip
content-length: 17
user-agent: RIF/0.2.0

Value: body-value"""[1:]

@when(u'the user runs RIF on that file passing in the appropriate variables')
def step_impl(context):
    variable_args = [
        '{}={}'.format(name, value)
        for name, value in context.variables.items()
    ]
    context.stdout, context.returncode = run_rif(
        [context.filename] + variable_args
    )


@given(u'a .rif file is on disk that has a higher rif file version')
def step_impl(context):
    context.filename = '/vol/tests/test-data/bad-version.rif'
    context.expected_error_msg = "Error parsing .rif file: rif file version " \
            "greater than maxium supported version - 0"

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
        ['/vol/build/rif'] + args,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
    )
    return_code = result.returncode
    output = result.stdout.decode('utf8')

    if return_code != 0:
        print(output)

    return (output, return_code)

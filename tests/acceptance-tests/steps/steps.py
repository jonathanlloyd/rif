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
user-agent: RIF/0.1.0"""[1:]

@given(u'a .rif file is on disk that describes a request with headers')
def step_impl(context):
    context.filename = '/vol/tests/test-data/basic-headers.rif'
    context.expected_result = """
GET /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.1.0
x-test-header: some_value"""[1:]

@given(u'a .rif file is on disk that describes a request with a body')
def step_impl(context):
    context.filename = '/vol/tests/test-data/basic-body.rif'
    context.expected_result = """
POST /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
content-length: 4
user-agent: RIF/0.1.0

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


def run_rif(args):
    result = subprocess.run(
        ['/vol/build/rif'] + args,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
    )
    return (result.stdout.decode('utf8'), result.returncode,)

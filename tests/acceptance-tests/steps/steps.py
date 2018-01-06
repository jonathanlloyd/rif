"""Step definitions for acceptance tests"""
import subprocess

from hamcrest import *


@given(u'a .rif file is on disk that describes a GET request from that server')
def step_impl(context):
    context.filename = '/vol/tests/test-data/basic-get.rif'

@when(u'the user runs RIF on that file')
def step_impl(context):
    context.stdout, context.returncode = run_rif([context.filename])

@then(u'RIF should return an echo of the request it made')
def step_impl(context):
    expected_result = """
GET /basic-get HTTP/1.1
host: localhost:8080
user-agent: RIF/0.1.0
accept-encoding: gzip"""[1:]

    print(context.stdout)
    assert_that(
        context.returncode,
        equal_to(0),
    )
    assert_that(
        context.stdout,
        contains_string(expected_result),
    )


def run_rif(args):
    result = subprocess.run(
        ['/vol/build/rif'] + args,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
    )
    return (result.stdout.decode('utf8'), result.returncode,)

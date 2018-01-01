"""Step definitions for acceptance tests"""
import subprocess

from hamcrest import *


@given(u'a .rif file is on disk that describes a GET request from that server')
def step_impl(context):
    context.filename = '/vol/tests/test_data/basic-get.rif'

@when(u'the user runs RIF on that file')
def step_impl(context):
    context.result = run_rif([context.filename])

@then(u'RIF should return an echo of the request it made')
def step_impl(context):
    expected_result = """
HTTP/1.1 GET /basic-get

Host: localhost:8080
Accept: */*"""
    assert_that(
        preprocess_echo_server(context.result),
        equal_to(preprocess_expectation(expected_result))
    )


def run_rif(args):
    return subprocess.run(
        ['/vol/build/rif'] + args,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
    ).stdout

def preprocess_echo_server(result):
    """Trims the hostname off the output from echoserver"""
    result = result.decode('utf-8')
    return '\n'.join(result.split('\n')[2:])

def preprocess_expectation(result):
    return '\n'.join(result.split('\n')[1:])

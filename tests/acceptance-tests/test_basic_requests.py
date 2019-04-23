import time
import subprocess

import pytest

ECHO_SERVER = None


# Fixtures
@pytest.fixture
def echo_server():
    global ECHO_SERVER

    if ECHO_SERVER is not None:
        return

    print('Starting echo server...')
    ECHO_SERVER = subprocess.Popen(
        ['./build/echo-server'],
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
    )
    time.sleep(1)


# Tests
def test_basic_get_request(echo_server):
    """Test that RIF can make basic GET requests"""
    output, return_code = run_rif(
        [make_test_file_path('basic-get.rif')],
    )

    expected_output = """
GET /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.4.5"""[1:]

    assert return_code == 0
    assert expected_output in output

def test_request_with_headers(echo_server):
    """Test that RIF can make requests with headers"""
    output, return_code = run_rif(
        [make_test_file_path('basic-headers.rif')],
    )

    expected_output = """
GET /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.4.5
x-test-header: some_value"""[1:]

    assert return_code == 0
    assert expected_output in output

def test_request_with_body(echo_server):
    """Test that RIF can make requests with body content"""
    output, return_code = run_rif(
        [make_test_file_path('basic-body.rif')],
    )

    expected_output = """
POST /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
content-length: 4
user-agent: RIF/0.4.5

test"""[1:]

    assert return_code == 0
    assert expected_output in output


# Utilities
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


def make_test_file_path(filename):
    return './tests/acceptance-tests/test-data/' + filename

import time
import subprocess

import pytest

ECHO_SERVER = None


# Fixtures
def echo_server():
    if ECHO_SERVER is not None:
        return

    global ECHO_SERVER
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
        [make_test_file_path('./basic-get.rif')],
    )

    assert return_code == 0
    assert output === """
GET /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.4.5"""[1:]


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


def test_file_path(filename):
    return './tests/acceptance-tests/test-data/' + filename

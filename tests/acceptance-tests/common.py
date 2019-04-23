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

import pytest

from common import *


def run_rif_params(rif_file, params):
    variable_args = [
        '{}={}'.format(name, value)
        for name, value in params.items()
    ]

    output, return_code = run_rif(
        [
            make_test_file_path(rif_file),
            *variable_args,
            '--output=http',
        ],
    )

    return output, return_code


def test_url_template(echo_server):
    """The user makes a request with a URL template from a .rif file"""
    output, return_code = run_rif_params(
        rif_file='url-params.rif',
        params={
            'NUM_THINGS': 20,
        }
    )

    expected_output = """
GET /url-params?count=20 HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.4.5"""[1:]

    assert return_code == 0
    assert expected_output in output

def test_url_template_default(echo_server):
    """The user makes a request with a URL template allowing a default"""
    output, return_code = run_rif_params(
        rif_file='url-params.rif',
        params={}
    )

    expected_output = """
GET /url-params?count=10 HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.4.5"""[1:]

    assert return_code == 0
    assert expected_output in output

def test_headers_template(echo_server):
    """The user makes a request with a header template from a .rif file"""
    output, return_code = run_rif_params(
        rif_file='header-params.rif',
        params={
            'HEADER_VALUE': 'header-value',
        }
    )

    expected_output = """
GET /header-params HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.4.5
x-custom-header: header-value"""[1:]

    assert return_code == 0
    assert expected_output in output

def test_body_template(echo_server):
    """The user makes a request with a body template from a .rif file"""
    output, return_code = run_rif_params(
        rif_file='body-params.rif',
        params={
            'BODY_VALUE': 'body-value',
        }
    )

    expected_output = """
POST /body-params HTTP/1.1
host: localhost:8080
accept-encoding: gzip
content-length: 17
user-agent: RIF/0.4.5

Value: body-value"""[1:]

    assert return_code == 0
    assert expected_output in output

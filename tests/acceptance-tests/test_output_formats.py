import pytest

from common import *


def run_rif_http(rif_file, params=None):
    if params is None:
        params = {}

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

    # These are invisible anyway and mess up the test string
    output = output.replace("\r", "")

    return output, return_code

def run_rif_curl(rif_file, params=None):
    if params is None:
        params = {}

    variable_args = [
        '{}={}'.format(name, value)
        for name, value in params.items()
    ]

    output, return_code = run_rif(
        [
            make_test_file_path(rif_file),
            *variable_args,
            '--output=curl',
        ],
    )

    # These are invisible anyway and mess up the test string
    output = output.replace("\r", "")

    return output, return_code


def test_basic_get_request_http_output(echo_server):
    """The user makes a basic GET request with a HTTP output format"""
    output, return_code = run_rif_http('basic-get.rif')

    expected_output = """
Request
-------
GET /basic-get HTTP/1.1
Host: localhost:8080
User-Agent: RIF/0.4.5
Accept-Encoding: gzip



Response
--------
HTTP/1.1 200 OK
Content-Length: 90
Content-Type: text/plain; charset=utf-8
Date: 

GET /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.4.5"""[1:]

    assert return_code == 0
    assert expected_output in output

def test_request_with_body_http_output(echo_server):
    """The user makes a request with a body with a HTTP output format"""
    output, return_code = run_rif_http('basic-body.rif')

    expected_output = """
Request
-------
POST /basic-get HTTP/1.1
Host: localhost:8080
User-Agent: RIF/0.4.5
Content-Length: 4
Accept-Encoding: gzip

test

Response
--------
HTTP/1.1 200 OK
Content-Length: 115
Content-Type: text/plain; charset=utf-8
Date: 

POST /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
content-length: 4
user-agent: RIF/0.4.5

test"""[1:]

    assert return_code == 0
    assert expected_output in output

def test_request_with_url_template_http_output(echo_server):
    """The user makes a request with a URL template and a HTTP output format"""
    output, return_code = run_rif_http(
        'url-params.rif',
        params={
            'NUM_THINGS': 20,
        }
    )

    expected_output = """
Request
-------
GET /url-params?count=20 HTTP/1.1
Host: localhost:8080
User-Agent: RIF/0.4.5
Accept-Encoding: gzip



Response
--------
HTTP/1.1 200 OK
Content-Length: 100
Content-Type: text/plain; charset=utf-8
Date: 

GET /url-params?count=20 HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.4.5"""[1:]

    assert return_code == 0
    assert expected_output in output

def test_request_with_body_curl_output(echo_server):
    """The user makes a request with a body with a cURL output format"""
    output, return_code = run_rif_curl('basic-body.rif')

    expected_output = """
cURL command
------------
curl -X 'POST' -d 'test' -H 'User-Agent: RIF/0.4.5' 'http://localhost:8080/basic-get'

Response
--------
POST /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
content-length: 4
user-agent: RIF/0.4.5

test"""[1:]

    assert return_code == 0
    assert expected_output in output

def test_request_with_url_template_curl_output(echo_server):
    """The user makes a request with a URL template and a cURL output format"""
    output, return_code = run_rif_curl(
        'url-params.rif',
        params={
            'NUM_THINGS': 20,
        }
    )

    expected_output = """
cURL command
------------
curl -X 'GET' -d '' -H 'User-Agent: RIF/0.4.5' 'http://localhost:8080/url-params?count=20'

Response
--------
GET /url-params?count=20 HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.4.5

"""[1:]

    assert return_code == 0
    assert expected_output in output

def test_unknown_output_format(echo_server):
    """The user makes a request with an unknown output format"""
    output, return_code = run_rif(
        [
            make_test_file_path('basic-get.rif'),
            '--output=NOT_AN_OUTPUT_FORMAT',
        ],
    )

    expected_output = "Unknown output format: NOT_AN_OUTPUT_FORMAT"

    assert return_code != 0
    assert expected_output in output

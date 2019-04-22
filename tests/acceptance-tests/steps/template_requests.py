"""Steps for templated requests"""
from hamcrest import *

from behave import *
from common import *


@given(u'a .rif file is on disk that has a URL template')
def step_impl(context):
    context.filename = test_file_path('url-params.rif')
    context.variables = {
        'NUM_THINGS': 20,
    }
    context.expected_plain_output = """
GET /url-params?count=20 HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.4.5"""[1:]
    context.expected_http_output = """
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
    context.expected_curl_output="""
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


@given(u'a .rif file is on disk that has a URL template with a default')
def step_impl(context):
    context.filename = test_file_path('url-params.rif')
    context.variables = {}
    context.expected_plain_output = """
GET /url-params?count=10 HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.4.5"""[1:]


@given(u'a .rif file is on disk that has a header template')
def step_impl(context):
    context.filename = test_file_path('header-params.rif')
    context.variables = {
        'HEADER_VALUE': 'header-value',
    }
    context.expected_plain_output = """
GET /header-params HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.4.5
x-custom-header: header-value"""[1:]


@given(u'a .rif file is on disk that has a body template')
def step_impl(context):
    context.filename = test_file_path('body-params.rif')
    context.variables = {
        'BODY_VALUE': 'body-value',
    }
    context.expected_plain_output = """
POST /body-params HTTP/1.1
host: localhost:8080
accept-encoding: gzip
content-length: 17
user-agent: RIF/0.4.5

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


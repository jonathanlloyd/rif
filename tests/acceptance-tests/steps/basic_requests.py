"""Step definitions for the basic request feature"""
from hamcrest import *

from behave import *
from common import *


@given(u'a .rif file is on disk that describes a GET request')
def step_impl(context):
    context.filename = '/vol/tests/test-data/basic-get.rif'
    context.expected_plain_output = """
GET /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.3.0"""[1:]
    context.expected_http_output = """
Request
-------
GET /basic-get HTTP/1.1
Host: localhost:8080
User-Agent: RIF/0.3.0
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
user-agent: RIF/0.3.0"""[1:]


@given(u'a .rif file is on disk that describes a request with headers')
def step_impl(context):
    context.filename = '/vol/tests/test-data/basic-headers.rif'
    context.expected_plain_output = """
GET /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
user-agent: RIF/0.3.0
x-test-header: some_value"""[1:]


@given(u'a .rif file is on disk that describes a request with a body')
def step_impl(context):
    context.filename = '/vol/tests/test-data/basic-body.rif'
    context.expected_plain_output = """
POST /basic-get HTTP/1.1
host: localhost:8080
accept-encoding: gzip
content-length: 4
user-agent: RIF/0.3.0

test"""[1:]

+++
title = "Passing Parameters"
draft = false
date = "2018-02-06T09:06:17Z"

+++
RIF allows you to pass parameters into your requests from the command line:
```
$ rif my-request.rif api_key="E225A111B0A5BDC6106DEDC6476C1A83"
```

This is done using RIF Templates and request variables:
```
rif_version: 0
url: "http://httpbin.org/get"
method: "GET"
headers:
  Authorization: "Bearer $(API_KEY)"
variables:
  API_KEY:
    type: string
```

# Templates
RIF supports template strings in the path, headers and body of your requests.

**Example:**
```
rif_version: 0
url: "http://httpbin.org/get?message=$(MESSAGE)"
method: "GET"
```
In this example, the RIF File contains a template in the path of the
request. When RIF is run on this file, the value of the `MESSAGE` variable
will be interpolated into the path.

## Syntax
You add a template variable to your RIF File by prepending the variable name
with a dollar sign and surrounding it with parentheses:
```
The following is the value of a variable: $(VARIABLE)
```
In this example, the value of `VARIABLE` will be interpolated into the string.

## Supported Template Locations
Template strings are supported in the following parts of a RIF file:

  - `url`
    - URL templates can be used to pass in query parameters or change hostnames
  - `headers` (both header names and their values)
    - Header templates can be used to send custom header values such as api keys
  - `body`
    - Body templates can be used to add custom fields in the request body
      such as the parameters in a POST request

# Variables
Once you have added some templates to your RIF File you will need to define
the type (and optionally the default value) of the variables you wish to
substitute.

```
rif_version: 0
url: "http://httpbin.org/get?message=$(MESSAGE)"
method: "GET"
variables:
  MESSAGE:
    type: string
    default: "Hello"
```

## Variable Types
When defining a variable in your RIF File you must specify the type. This will
be used to parse and validate the variables that you pass in. The following
types are currently supported by RIF:

  - `boolean`
    - Represents true/false values.
    - Must be one of `true`|`false`|`True`|`False`
    - Example: `true`
  - `number`
    - Represents integer and decimal values
    - Will be rounded to the nearest 64bit floating point value
    - Example: `12`
  - `string`
    - Represents textual values
    - Full unicode support where allowed in the request.
    - Example: `"hello"`

## Variable Defaults
When defining a variable you can optionally specify a default value. This value
will be used if the variable is not passed in over the command line. If a
variable is not passed in and does not have a default value RIF will return an
error. The default value of a variable must conform to the type of that variable.

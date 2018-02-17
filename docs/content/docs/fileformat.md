+++
title = "RIF File Format"
draft = false
date = "2018-02-13T20:01:33Z"

+++
# Schema
## RIF File
- rif\_version ( integer | **required** )
 - The RIF file format version of the file. Currently, the only supported
   version is 0.
- url ( string | **required** )
 - The URL of the documented request
- method ( string | **required** )
 - The HTTP method of the documented request. Must be a standard method
   ("GET", "POST", "PUT", "PATCH", "DELETE", etc.)
- headers ( Map\<string,string\> | optional )
 - A map from header name to value. E.g. `Content-Type: "application/json"`.
- body ( string | optional )
 - The body of the documented request.
- variables ( Map\<string,VariableDefinition\> | optional)
 - A map from variable names to their definitions. These variables are used
   to populate the template strings in your request.

## VariableDefinition
- type (string | required )
 - The type of the variable you wish to use in your template (used to
   parse/validate the variable value). Must be one of: string, boolean or number.
- default ( string or boolean or number | required )
 - The default value of the variable in your template. Must be of the appropriate
   type based on the type field of the variable definition.

# Example
```
rif_version: 0
url: http://httpbin.org/post
method: POST
headers:
  Content-Type: application/json
body: |
  {
    "command": "$(commandType)"
  }
variables:
  commandType:
    type: string
    default: OPEN
```

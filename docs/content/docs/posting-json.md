+++
title = "POSTing JSON"
draft = false
date = "2018-02-05T21:32:27Z"

+++
In this guide you will learn how to use RIF to send a `POST` request with a
JSON body. Then, you will learn how to use templates to add a variable to the
JSON payload.

# Create the RIF file
First we will create a `.rif` file that describes a `POST` request to
`http://httpbin.org/post`. This endpoint will echo back the body of the request
we are going to send. Save the following file as `post-json-body.rif`:
```
rif_version: 0
url: "http://httpbin.org/post"
method: "POST"
```

Then run RIF on the file:
```
$ rif post-json-body.rif
```

You should see something like the following output:
```
{
  "args": {}, 
  "data": "", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Connection": "close", 
    "Content-Length": "0", 
    "Host": "httpbin.org", 
    "User-Agent": "RIF/0.3.1"
  }, 
  "json": null, 
  "origin": "<YOUR IP ADDRESS>", 
  "url": "http://httpbin.org/post"
}
```

# Add a JSON body
The next step is to add a JSON body to the request. You can do this using the
[multiline yaml syntax](http://www.yaml.org/spec/1.2/spec.html#id2795688):
```
rif_version: 0
url: "http://httpbin.org/post"
method: "POST"
headers:
  content-type: "application/json"
body: |
  {
    "message": "Hello World!"
  }
```

This adds a JSON body along with the `application/json` content type header
to the request. If you run rif on this file you will see that the JSON body
has been added to the echoed response:
```
{
  "args": {}, 
  "data": "{\n  \"message\": \"Hello World!\"\n}\n", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Connection": "close", 
    "Content-Length": "32", 
    "Content-Type": "application/json", 
    "Host": "httpbin.org", 
    "User-Agent": "RIF/0.3.1"
  }, 
  "json": {
    "message": "Hello World!"
  }, 
  "origin": "<YOUR IP ADDRESS>", 
  "url": "http://httpbin.org/post"
}
```

# Add a variable
The final step is to add a template variable to your JSON body so that
you can substitute this value when calling RIF:
```
rif_version: 0
url: "http://httpbin.org/post"
method: "POST"
headers:
  content-type: "application/json"
body: |
  {
    "message": "$(body)"
  }
variables:
  body:
    type: "string"
```

You can now pass in a message as a command line argument when calling RIF:
```
$ rif post-json-body.rif body="RIF Rocks!"
```

And you should see the result returned in the request response:
```
{
  "args": {}, 
  "data": "{\n  \"message\": \"RIF Rocks!\"\n}\n", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Connection": "close", 
    "Content-Length": "30", 
    "Content-Type": "application/json", 
    "Host": "httpbin.org", 
    "User-Agent": "RIF/0.3.1"
  }, 
  "json": {
    "message": "RIF Rocks!"
  }, 
  "origin": "<YOUR IP ADDRESS>", 
  "url": "http://httpbin.org/post"
}
```

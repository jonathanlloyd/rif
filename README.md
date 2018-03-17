<p align="center">
  <img src="https://github.com/turingincomplete/rif/blob/master/assets/logo-text.svg" alt="RIF logo"/>
</p>
<p align="center">
  <a href="https://travis-ci.org/turingincomplete/rif">
    <img src="https://travis-ci.org/turingincomplete/rif.svg?branch=master" alt="Build Status"/>
  </a>
  <a href="https://goreportcard.com/report/github.com/turingincomplete/rif">
    <img src="https://goreportcard.com/badge/github.com/turingincomplete/rif" alt="Go Report Card"/>
  </a>
</p>
<p align="center">
  <b>
    <a href="https://turingincomplete.github.io/rif">
      Website
    </a>
  </b>
  &nbsp;|&nbsp;
  <b>
    <a href="https://turingincomplete.github.io/rif/docs/quickstart/">
      Documentation
    </a>
  </b>
  &nbsp;|&nbsp;
  <b>
    <a href="https://github.com/turingincomplete/rif/releases">
      Download
    </a>
  </b>
</p>

---
RIF makes working with and testing HTTP APIs a breeze. Create `.rif` files
for all those repetitive requests and never wrestle with cURL again.

![Terminal Example](docs/static/img/terminal.svg)

## Quickstart
To get started, we will be making a simple GET request to
[httpbin.org/get](http://httpbin.org/get). This endpoint returns the details
of GET requests back to the client as JSON.

We will be passing in a URL parameter called `message` that is parameterised
using RIFs variable templating feature.

Open your editor of choice and save the following file to your computer
as `gethttpbin.rif`:
```
rif_version: 0
url: "http://httpbin.org/get?message=hello%20$(PLACE)"
method: "GET"
variables:
  PLACE:
    type: "string"
    default: "world"
```

Next, open your terminal in the same location and run RIF,
passing in the file you just created:
```
$ rif ./gethttpbin.rif
```

If all goes well you should see something like the following:
```
{
  "args": {
    "message": "hello world"
  }, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Connection": "close", 
    "Host": "httpbin.org", 
    "User-Agent": "RIF/0.2.0"
  }, 
  "origin": "<YOUR IP ADDRESS>", 
  "url": "http://httpbin.org/get?message=hello world"
}
```

### Variable Templating

Now let's use RIF's variable templating feature to override our welcome message.
Paste the following command into your terminal:
```
$ rif ./gethttpbin.rif PLACE=universe
```

You should now see that the response has changed:
```
{
  "args": {
    "message": "hello universe"
  }, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Connection": "close", 
    "Host": "httpbin.org", 
    "User-Agent": "RIF/0.2.0"
  }, 
  "origin": "<YOUR IP ADDRESS>", 
  "url": "http://httpbin.org/get?message=hello universe"
}
```

Congratulations! You have just made and executed your first `.rif` file!

<p align="center">
  <a href="https://jonathanlloyd.github.io/rif">
    <img src="https://github.com/jonathanlloyd/rif/blob/master/assets/logo-text.svg" alt="RIF logo"/>
  </a>
</p>
<p align="center">
  <a href="https://travis-ci.org/jonathanlloyd/rif">
    <img src="https://travis-ci.org/jonathanlloyd/rif.svg?branch=master" alt="Build Status"/>
  </a>
  <a href="https://goreportcard.com/report/github.com/jonathanlloyd/rif">
    <img src="https://goreportcard.com/badge/github.com/jonathanlloyd/rif" alt="Go Report Card"/>
  </a>
  <a href="https://github.com/jonathanlloyd/rif/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/jonathanlloyd/rif.svg" alt="GPL3 Licensed"/>
  </a>
  <a href="https://github.com/jonathanlloyd/rif/releases">
    <img src="https://img.shields.io/github/downloads/jonathanlloyd/rif/total.svg" alt="Download count"/>
  </a>
</p>
<p align="center">
  <b>
    <a href="https://jonathanlloyd.github.io/rif">
      Website
    </a>
  </b>
  &nbsp;|&nbsp;
  <b>
    <a href="https://jonathanlloyd.github.io/rif/docs/quickstart/">
      Documentation
    </a>
  </b>
  &nbsp;|&nbsp;
  <b>
    <a href="https://github.com/jonathanlloyd/rif/releases">
      Download
    </a>
  </b>
</p>

---
Keep your Requests In Files.

RIF is a command-line tool that allows you to store repetitive and complex
HTTP requests as files to make testing and interacting with HTTP APIs a breeze.

# Installation
## Quick Install/Upgrade
### Linux (64bit)
To install/upgrade RIF on a standard 64bit linux install, run the following in
your shell:
```
curl -Lo rif.tar $(curl https://api.github.com/repos/jonathanlloyd/rif/releases/latest 2>/dev/null | grep -o http[^[:space:]]*linux_amd64\.tar\.gz) 2>/dev/null && tar -xf rif.tar rif && chmod +x rif && sudo mv rif /usr/local/bin/rif && rm ./rif.tar && rif --help && rif --version
```

## Verifying Your Installation
To verify that you have installed RIF correctly, run the following command
in your terminal:
```
$ rif --version
```

If RIF is correctly installed you should see the version and build
number printed to the screen:
```
Version: <expected version number>
Build: <build number>
```

## Manual Install (other platforms)
RIF is distributed as a single binary executable. To install it all you need
to do is download the correct binary from the
[downloads page](https://github.com/jonathanlloyd/rif/releases "Downloads Page")
and put it somewhere in your path.

## Supported Platforms
### macOS
#### Download
 - If you are running a 64bit system:
   - `rif_<version>_darwin_amd64.tar.gz`
 - If you are running a 32bit system:
   - `rif_<version>_darwin_386.tar.gz`

#### Path
Extract the appropriate binary and add it to a directory in your
[path](http://osxdaily.com/2014/08/14/add-new-path-to-path-command-line/ "How to Add a New Path to PATH at Command Line the Right Way").
We recommend that you put it in `/user/local`.

### Linux
#### Download
 - If you are running a 64bit system:
   - `rif_<version>_linux_amd64.tar.gz`
 - If you are running a 32bit system:
   - `rif_<version>_linux_386.tar.gz`

#### Path
Extract the appropriate binary and add it to a directory in your
[path](https://www.cyberciti.biz/faq/how-to-add-to-bash-path-permanently-on-linux/ "How to add to bash $PATH permanently on Linux").
We recommend that you put it in `/usr/local/bin`.

### Windows
#### Download
 - If you are running a 64bit system:
   - `rif_<version>_windows_amd64.tar.gz`
 - If you are running a 32bit system:
   - `rif_<version>_windows_386.tar.gz`

#### Path
Extract the appropriate binary and add it to a directory in your
[path](https://stackoverflow.com/questions/1618280/where-can-i-set-path-to-make-exe-on-windows "Where can I set the path on Windows?").
We recommend that you put it in `c:\RIF` and add this directory to your path.

# Quickstart
## Making Your First Request
To get started, we will be making a simple GET request to
[httpbin.org/get](http://httpbin.org/get). This endpoint returns the details
of GET requests back to the client as JSON.

We will be passing in a URL parameter called `message` that is parameterised
using RIFs variable templating feature.

Open your editor of choice and save the following file to your computer
as `gethttpbin.rif`:
```
rif_version: 0
url: "http://httpbin.org/get?message=hello%20$(place)"
method: "GET"
variables:
  place:
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

## Variable Templating
Now let's use RIF's variable templating feature to override our welcome message.
Paste the following command into your terminal:
```
$ rif ./gethttpbin.rif place=universe
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

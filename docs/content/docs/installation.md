+++
title = "Installation"
draft = false
date = "2017-12-30T21:14:46Z"

+++
# Quick Install
## Linux (64bit)
To install/upgrade RIF on a standard 64bit linux install, run the following in
your shell:
```
curl -Lo rif.tar $(curl https://api.github.com/repos/jonathanlloyd/rif/releases/latest 2>/dev/null | grep -o http[^[:space:]]*linux_amd64\.tar\.gz) 2>/dev/null && tar -xf rif.tar rif && chmod +x rif && sudo mv rif /usr/local/bin/rif && rm ./rif.tar
```

# Manual Install (other platforms)
RIF is distributed as a single binary executable. To install it all you need
to do is download the correct binary from the
[downloads page](https://github.com/jonathanlloyd/rif/releases "Downloads Page")
and put it somewhere in your path.

# Verifying Your Installation
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

# Supported Platforms
## macOS
### Download
 - If you are running a 64bit system:
   - `rif_<version>_darwin_amd64.tar.gz`
 - If you are running a 32bit system:
   - `rif_<version>_darwin_386.tar.gz`

### Path
Extract the appropriate binary and add it to a directory in your
[path](http://osxdaily.com/2014/08/14/add-new-path-to-path-command-line/ "How to Add a New Path to PATH at Command Line the Right Way").
We recommend that you put it in `/user/local`.

## Linux
### Download
 - If you are running a 64bit system:
   - `rif_<version>_linux_amd64.tar.gz`
 - If you are running a 32bit system:
   - `rif_<version>_linux_386.tar.gz`

### Path
Extract the appropriate binary and add it to a directory in your
[path](https://www.cyberciti.biz/faq/how-to-add-to-bash-path-permanently-on-linux/ "How to add to bash $PATH permanently on Linux").
We recommend that you put it in `/usr/local/bin`.

## Windows
### Download
 - If you are running a 64bit system:
   - `rif_<version>_windows_amd64.tar.gz`
 - If you are running a 32bit system:
   - `rif_<version>_windows_386.tar.gz`

### Path
Extract the appropriate binary and add it to a directory in your
[path](https://stackoverflow.com/questions/1618280/where-can-i-set-path-to-make-exe-on-windows "Where can I set the path on Windows?").
We recommend that you put it in `c:\RIF` and add this directory to your path.

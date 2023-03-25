# WSnok

WSnok is a alternative to Wget. It allows you to download a file via a simple command line using multiple threads to increase the download speed (if posible).

## Snok
A snok is a hard pull in Dutch.

## Usage
```shell
wsnok [OPTION]... [URL]...
```
## Example
Download the ubuntu iso file with 50 threads.
```shell
wsnok -t50 https://releases.ubuntu.com/22.04.2/ubuntu-22.04.2-desktop-amd64.iso
```
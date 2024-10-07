# UniFi Captive Portal

A [UniFi](https://www.ubnt.com) external captive portal with ChatGPT.

# Demo
![Chat Demo](/doc/demo.png "Demo")

## Configuration

Env variables | Value | Default
---------- | ----- | -------
UnifiURL | Full URL of your UniFi Controller. Be sure to include the port it is running on | https://unifi:8443
UnifiUsername | Username of the user to make API calls with. It is recommended to use a dedicated user | admin
UnifiPassword | Password for user defined above | admin
UnifiSite | The name of the site the APs/Users reside in | default
Title | Title used in HTML pages as well as headings. Usually you will put your company name here | Captive Portal
RedirectUrl | URL to redirect users to if they do not provide one to the controller | https://captive.apple.com/
Minutes | Amount of time to register user for | 600


## Build Binary
```shell
$ make build
```

## Build Docker
```shell
$ make docker
```
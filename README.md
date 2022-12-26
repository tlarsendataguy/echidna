<div style="text-align: center">
<img height="200" src="https://raw.githubusercontent.com/tlarsendataguy/echidna/master/logo.png"/>
</div>

# Echidna

Echidna is a small webserver application to serve static websites. It automatically obtains and renews TLS certificates and can host multiple domains from a single running service.

## What is it not

Echidna is NOT:

* A framework
* A module you can import into your own application
* A server-side rendering solution
* A simple, plug-and-play API generator

## What do I do with Echidna?

You can use Echidna out-of-the-box to serve basic websites and blogs from lightweight VMs (with public IP addresses) hosted by major cloud providers. Echidna handles all the web serving and certificate tasks, leaving you free to develop your website using HTML, CSS, and Javascript. Echidna also lets you serve multiple domains from a single VM, so you can try out many ideas without needing to expand your infrastructure.

More importantly, you are encouraged to change, extend, and improve the code. Make it your own. Does your site need an API? Add new handlers! Do you want to use templating? Code it in! Echidna is not a module, so it's not opinionated about how you call it. Change it to make it work for you.

## Who is Echidna for?

Echidna is for developers who have existing sites or blogs today on hosted platforms, but want to control everything. They are looking for an easy way to get started with their own web server, without having to write the base code themselves. As they become more comfortable managing their webserver, they want to be able to change it to whatever they desire.

## Configuration

Echidna uses a single JSON configuration file with the following structure:

```
{
    "CertFolder": string,
    "ServeFolder": string,
    "Hosts": [
        {
            "Folder": string
            "HostWhiteList": [string]
        }
    ]
}
```

* CertFolder specified a folder where TLS certs will be persisted to disk.
* ServeFolder is the root folder from which website files will be served
* Hosts is a list of JSON objects identifying the sites that will be served. Each host gets it own subfolder under the ServeFolder.
  * Folder specifies the subfolder which contains the files for the host
  * HostWhiteList is an array of strings containing the domains that apply to the host. Commonly, the list will contain the raw domain and the www subdomain.

## Sample

You can download Echidna for Windows in the Releases section. The program, configuration files, required directories, and sample sites are all included in sample.zip. The sample is the easiest way to get started with Echidna.

Executables are not provided for Linux or Unix. You will need to build and deploy those yourself.

To deploy Echidna as a service on Windows, use the [Non-Sucking Service Manager](http://nssm.cc)
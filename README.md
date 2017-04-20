### Source
- https://github.com/chenhw2/changeip-ddns-cli
  
### Thanks to
- https://www.changeip.com/dns.php
  
### Docker
- https://hub.docker.com/r/chenhw2/changeip-ddns-cli/
  
### Usage
```
$ docker pull chenhw2/changeip-ddns-cli

$ docker run -d \
    -e "Username=1234567890" \
    -e "Password=abcdefghijklmn" \
    -e "Domain=ddns.changeip.com" \
    -e "Redo=600" \
    chenhw2/changeip-ddns-cli

```
### Help
```
$ docker run --rm chenhw2/changeip-ddns-cli -h
NAME:
   ChangeIP - changeip-ddns-cli

USAGE:
   changeip-ddns-cli [global options] command [command options] [arguments...]

VERSION:
   MISSING build version [git hash]

COMMANDS:
     help, h  Shows a list of commands or help for one command
   DDNS:
     update       Update ChangeIP's DNS DomainRecords Record
     auto-update  Auto-Update ChangeIP's DNS DomainRecords Record, Get IP using its getip
   GET-IP:
     getip         Get IP Combine 11 different Web-API
     getdns        Get IP of A domain Combine 5 different DNS-Server

GLOBAL OPTIONS:
   --username value, -u value  Your User ID of ChangeIP.Com
   --password value, -p value  Your Password of ChangeIP.Com
   --help, -h                  show help
   --version, -v               print the version

```
### CLI Example:
```
changeip -u ${Username} -p ${Password} \
    auto-update --domain ddns.changeip.com

changeip -u ${Username} -p ${Password} \
    update --domain ddns.changeip.com \
    --ipaddr $(ifconfig pppoe-wan | sed -n '2{s/[^0-9]*://;s/[^0-9.].*//p}')

```
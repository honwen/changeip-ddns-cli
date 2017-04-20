package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
)

type AccessKey struct {
	Username string
	Password string
}

func (ak *AccessKey) isFiled() bool {
	return len(ak.Username) > 0 && len(ak.Password) > 0
}

func (ak *AccessKey) doDDNSUpdate(fulldomain, ipaddr string) (err error) {
	// http://www.changeip.com/accounts/knowledgebase.php?action=displayarticle&id=34
	uri := "https://nic.ChangeIP.com/nic/update?u=%s&p=%s&hostname=%s&ip=%s"
	if getDNS(fulldomain) == ipaddr {
		return
	}
	resp := wGet(fmt.Sprintf(uri, ak.Username, ak.Password, fulldomain, ipaddr), minTimeout*20)
	if !strings.Contains(resp, "Successful") {
		err = errors.New(resp)
	}
	return
}

var (
	accessKey AccessKey
	version   = "MISSING build version [git hash]"
)

func main() {
	app := cli.NewApp()
	app.Name = "ChangeIP"
	app.Usage = "changeip-ddns-cli"
	app.Version = version
	app.Commands = []cli.Command{
		{
			Name:     "update",
			Category: "DDNS",
			Usage:    "Update ChangeIP's DNS DomainRecords Record",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "domain, d",
					Usage: "Specific `DomainName`. like ddns.changeip.com",
				},
				cli.StringFlag{
					Name:  "ipaddr, i",
					Usage: "Specific `IP`. like 1.2.3.4",
				},
			},
			Action: func(c *cli.Context) error {
				if err := appInit(c); err != nil {
					return err
				}
				// fmt.Println(c.Command.Name, "task: ", accessKey, c.String("domain"), c.String("ipaddr"))
				if err := accessKey.doDDNSUpdate(c.String("domain"), c.String("ipaddr")); err != nil {
					log.Printf("%+v", err)
				} else {
					log.Println(c.String("domain"), c.String("ipaddr"))
				}
				return nil
			},
		},
		{
			Name:     "auto-update",
			Category: "DDNS",
			Usage:    "Auto-Update ChangeIP's DNS DomainRecords Record, Get IP using its getip",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "domain, d",
					Usage: "Specific `DomainName`. like ddns.changeip.com",
				},
				cli.Int64Flag{
					Name:  "redo, r",
					Value: 0,
					Usage: "redo Auto-Update, every N `Seconds`; Disable if N less than 10",
				},
			},
			Action: func(c *cli.Context) error {
				if err := appInit(c); err != nil {
					return err
				}
				// fmt.Println(c.Command.Name, "task: ", accessKey, c.String("domain"), c.Int64("redo"))
				redoDurtion := c.Int64("redo")
				for {
					autoip := getIP()
					if err := accessKey.doDDNSUpdate(c.String("domain"), autoip); err != nil {
						log.Printf("%+v", err)
					} else {
						log.Println(c.String("domain"), autoip)
					}
					if redoDurtion < 10 {
						break // Disable if N less than 10
					}
					time.Sleep(time.Duration(redoDurtion) * time.Second)
				}
				return nil
			},
		},
		{
			Name:     "getip",
			Category: "GET-IP",
			Usage:    "      Get IP Combine 11 different Web-API",
			Action: func(c *cli.Context) error {
				// fmt.Println(c.Command.Name, "task: ", c.Command.Usage)
				fmt.Println(getIP())
				return nil
			},
		},
		{
			Name:     "getdns",
			Category: "GET-IP",
			Usage:    "      Get IP of A domain Combine 5 different DNS-Server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "domain, d",
					Usage: "Specific `DomainName`. like ddns.changeip.com",
				},
			},
			Action: func(c *cli.Context) error {
				// fmt.Println(c.Command.Name, "task: ", c.String("domain"))
				if 0 == len(c.String("domain")) {
					cli.ShowAppHelp(c)
				} else {
					fmt.Println(getDNS(c.String("domain")))
				}
				return nil
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "username, u",
			Usage: "Your User ID of ChangeIP.Com",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "Your Password of ChangeIP.Com",
		},
	}
	app.Action = func(c *cli.Context) error {
		return appInit(c)
	}
	app.Run(os.Args)
}

func appInit(c *cli.Context) error {
	accessKey.Username = c.GlobalString("username")
	accessKey.Password = c.GlobalString("password")
	if !accessKey.isFiled() {
		cli.ShowAppHelp(c)
		return errors.New("Username/Password is empty")
	}
	return nil
}

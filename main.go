package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
)

type AUTH struct {
	Username string
	Password string
}

func (ak *AUTH) isFiled() bool {
	return len(ak.Username) > 0 && len(ak.Password) > 0
}

func (ak *AUTH) doDDNSUpdate(fulldomain, ipaddr string) (err error) {
	// http://www.changeip.com/accounts/knowledgebase.php?action=displayarticle&id=34
	uri := "https://nic.ChangeIP.com/nic/update?u=%s&p=%s&hostname=%s&ip=%s"
	if getDNS(fulldomain) == ipaddr {
		return
	}
	resp := wGet(fmt.Sprintf(uri, ak.Username, ak.Password, fulldomain, ipaddr), minTimeout*20)
	if !strings.Contains(resp, "Successful") {
		err = errors.New("wGet Err: " + resp)
	}
	return
}

var (
	auth    AUTH
	version = "MISSING build version [git hash]"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	app := cli.NewApp()
	app.Name = "ChangeIP"
	app.Usage = "changeip-ddns-cli"
	app.Version = fmt.Sprintf("Git:[%s] (%s)", strings.ToUpper(version), runtime.Version())
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
				if err := auth.doDDNSUpdate(c.String("domain"), c.String("ipaddr")); err != nil {
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
				cli.StringFlag{
					Name:  "redo, r",
					Value: "",
					Usage: "redo Auto-Update, every N `Seconds`; Disable if N less than 10; End with [Rr] enable random delay: [N, 2N]",
				},
			},
			Action: func(c *cli.Context) error {
				if err := appInit(c); err != nil {
					return err
				}
				// fmt.Println(c.Command.Name, "task: ", accessKey, c.String("domain"), c.Int64("redo"))
				redoDurtionStr := c.String("redo")
				if len(redoDurtionStr) > 0 && !regexp.MustCompile(`\d+[Rr]?$`).MatchString(redoDurtionStr) {
					return errors.New(`redo format: [0-9]+[Rr]?$`)
				}
				randomDelay := regexp.MustCompile(`\d+[Rr]$`).MatchString(redoDurtionStr)
				redoDurtion := 0
				if randomDelay {
					// Print Version if exist
					if !strings.HasPrefix(version, "MISSING") {
						fmt.Fprintf(os.Stderr, "%s %s\n", strings.ToUpper(c.App.Name), c.App.Version)
					}
					redoDurtion, _ = strconv.Atoi(redoDurtionStr[:len(redoDurtionStr)-1])
				} else {
					redoDurtion, _ = strconv.Atoi(redoDurtionStr)
				}
				for {
					autoip := getIP()
					if err := auth.doDDNSUpdate(c.String("domain"), autoip); err != nil {
						log.Printf("%+v", err)
					} else {
						log.Println(c.String("domain"), autoip)
					}
					if redoDurtion < 10 {
						break // Disable if N less than 10
					}
					if randomDelay {
						time.Sleep(time.Duration(redoDurtion+rand.Intn(redoDurtion)) * time.Second)
					} else {
						time.Sleep(time.Duration(redoDurtion) * time.Second)
					}
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
	auth.Username = c.GlobalString("username")
	auth.Password = c.GlobalString("password")
	if !auth.isFiled() {
		cli.ShowAppHelp(c)
		return errors.New("Username/Password is empty")
	}
	return nil
}

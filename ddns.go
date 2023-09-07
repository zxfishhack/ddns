package main

import (
	"github.com/zxfishhack/ddns/pkg/dns"
	"github.com/zxfishhack/ddns/pkg/ip"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	subDomains := strings.Split(os.Getenv("DDNS_SUBDOMAIN"), ",")
	domainNames := strings.Split(os.Getenv("DDNS_DOMAINNAME"), ",")
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	stopC := make(chan os.Signal)
	signal.Notify(stopC, syscall.SIGTERM)
	signal.Notify(stopC, syscall.SIGINT)
	printIp := true
	err := dns.Init()
	if err != nil {
		log.Fatal("init dns service failed.", err)
	}

	for {
		select {
		case <-stopC:
			return
		case <-timer.C:
		}
		if curIp, err := ip.GetMyIP(); err != nil {
			log.Printf("get current ip failed: %v", err)
		} else {
			needPrintIp := false
			for _, domainName := range domainNames {
				for _, subDomain := range subDomains {
					domain := subDomain + "." + domainName
					lastIp, rid, err := dns.GetDomainRecord(subDomain, domainName)
					if printIp {
						log.Printf("check domain record: %v -> %v, err: %v", domain, lastIp, err)
					}
					if err == os.ErrNotExist {
						err = dns.CreateDomainRecord(subDomain, domainName, curIp)
						log.Printf("create domain record: %v -> %v, err: %v", domain, curIp, err)
						needPrintIp = true
					} else if err != nil {
						log.Printf("get domain record: %v, err: %v", domain, err)
						needPrintIp = true
					} else if lastIp != curIp {
						err = dns.UpdateDomainRecord(subDomain, domainName, curIp, rid)
						log.Printf("update domain record: %v -> %v, last: %v, err: %v", domain, curIp, lastIp, err)
						needPrintIp = true
					}
				}
			}
			printIp = needPrintIp
		}
		timer.Reset(15 * time.Second)
	}

}

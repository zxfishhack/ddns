package aliyun

import (
	"errors"
	"github.com/zxfishhack/ddns/pkg/dns/inf"
	"log"
	"os"

	sdkerrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

type aliyunDns struct {
	client *alidns.Client
}

func New() (s inf.IDNSService, err error) {
	aliyun := &aliyunDns{}
	region := os.Getenv("ALIDNS_REGION")
	key := os.Getenv("ALIDNS_KEY")
	secret := os.Getenv("ALIDNS_SECRET")
	if region == "" || key == "" || secret == "" {
		err = inf.ErrNotAvailable
		return
	}
	aliyun.client, err = alidns.NewClientWithAccessKey(region, key, secret)

	if err == nil {
		s = aliyun
	}

	return
}

func (s *aliyunDns) GetDomainRecord(subDomain, domainName string) (ip, rid string, err error) {
	request := alidns.CreateDescribeSubDomainRecordsRequest()

	request.SubDomain = subDomain + "." + domainName
	request.Type = "A"
	request.DomainName = domainName

	response, err := s.client.DescribeSubDomainRecords(request)
	if err != nil {
		return
	}
	if !response.IsSuccess() {
		err = errors.New("GetDomainRecord failed")
	} else if len(response.DomainRecords.Record) == 0 {
		err = os.ErrNotExist
	} else {
		ip = response.DomainRecords.Record[0].Value
		rid = response.DomainRecords.Record[0].RecordId
	}
	return
}

func (s *aliyunDns) CreateDomainRecord(subDomain, domainName, ip string) (err error) {
	log.Printf("CreateDomainRecord: subDomain: %s domainName: %s, ip: %s", subDomain, domainName, ip)
	request := alidns.CreateAddDomainRecordRequest()

	request.Value = ip
	request.Type = "A"
	request.RR = subDomain
	request.DomainName = domainName

	response, err := s.client.AddDomainRecord(request)
	if err != nil {
		if sdkerr, ok := err.(sdkerrors.Error); ok {
			if sdkerr.ErrorCode() == "DomainRecordDuplicate" {
				err = nil
			}
		}
		return
	}
	if !response.IsSuccess() {
		err = errors.New("")
	}
	return
}

func (s *aliyunDns) UpdateDomainRecord(subDomain, domainName, ip, rid string) (err error) {
	request := alidns.CreateUpdateDomainRecordRequest()

	request.Value = ip
	request.Type = "A"
	request.RR = subDomain
	request.RecordId = rid

	response, err := s.client.UpdateDomainRecord(request)
	if err != nil {
		return
	}
	if !response.IsSuccess() {
		err = errors.New("")
	}
	return
}

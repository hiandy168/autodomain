package net

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"
	"strings"
	"time"

	"yiyecp.com/autodomain/config"

	"fmt"

	"github.com/denverdino/aliyungo/dns"
)

const IPCHECKURL = "http://ddns.oray.com/checkip"
const IPCHECKURLA = "http://yiye.yiyecp.cn/index.php"
const IPCHECKURLB = "http://yiye.yiyecp.com/index.php"

func GetMyIp() (string, error) {
	JsonParse := config.NewJsonStruct()
	v := config.Configdata{}
	homepath, _ := Home()
	JsonParse.LoadJson(homepath+"/"+"config.json", &v)
	var result string
	for _, u := range v.CheckUrls {
		resp, err := http.Get(u)
		if err != nil {
			//return "", nil
			continue
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		src := string(body)
		re := regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])`)
		src = re.FindString(src)
		if src != "" {
			result = strings.TrimSpace(src)
			fmt.Println(result)
		}
	}

	return result, nil
}

func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

func ModifyDomainRecordsValue() (err error) {
	JsonParse := config.NewJsonStruct()
	v := config.Configdata{}
	homepath, _ := Home()
	JsonParse.LoadJson(homepath+"/"+"config.json", &v)

	client := dns.NewClient(v.Id, v.Secret)
	args := &dns.DescribeSubDomainRecordsArgs{}
	args.SubDomain = v.Domain
	dr, err := client.DescribeSubDomainRecords(args)
	if err != nil {
		return err
	}

	if len(dr.DomainRecords.Record) != 1 {
		return errors.New("错误:该域名下绑定地址多余１个或找不到该域名!")
	}

	ur := &dns.UpdateDomainRecordArgs{}
	ur.RR = dr.DomainRecords.Record[0].RR
	ur.RecordId = dr.DomainRecords.Record[0].RecordId
	ur.Type = "A"
	NowIp, err := GetMyIp()

	if err != nil {
		return errors.New("错误:获取当前IP错误!")
	}
	if NowIp == "" {
		return errors.New("错误:获取当前IP为空!")
	}
	ur.Value = NowIp
	if !strings.EqualFold(dr.DomainRecords.Record[0].Value, ur.Value) {
		result, err := client.UpdateDomainRecord(ur)
		if err != nil {
			return err
		} else {
			log.Println(result)
			return nil
		}
	} else {
		return errors.New("提示!:待修改地址与域名绑定地址一致，无需同步!")
	}
}

func RunStart() {
	err := ModifyDomainRecordsValue()
	if err != nil {
		if strings.Index(err.Error(), "the same as old") != -1 {
			err = errors.New("提示:待修改域名与绑定域名一致，无需同步!")
		} else if strings.Index(err.Error(), "name does not belong") != -1 {
			err = errors.New("错误:当前用户下找不到该域名!")
		}
		log.Println(err.Error())
	} else {
		log.Println("提示:同步成功!")
	}
}

func GoRun() {
	JsonParse := config.NewJsonStruct()
	v := config.Configdata{}
	homepath, _ := Home()
	JsonParse.LoadJson(homepath+"/"+"config.json", &v)

	t1 := time.NewTimer(time.Second * time.Duration(v.Time))
	for {
		select {
		case <-t1.C:
			RunStart()
			t1.Reset(time.Second * time.Duration(v.Time))
		}
	}
}

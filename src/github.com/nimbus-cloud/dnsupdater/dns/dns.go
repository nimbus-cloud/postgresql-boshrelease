package dns

import (
	"fmt"
	"strings"
	"time"

	"github.com/nimbus-cloud/dnsupdater/config"

	"github.com/cloudfoundry/bosh-utils/system"
	"github.com/cloudfoundry/bosh-utils/logger"
)

const dnsUpdateInterval = 60 * time.Second
const logTag = "DnsUpdater"

type DnsUpdater struct {
	cfg config.Config
	logger logger.Logger
	fs system.FileSystem
	cmdRunner system.CmdRunner
}

func NewDnsUpdater(cfg config.Config, log logger.Logger) *DnsUpdater {
	fs := system.NewOsFileSystem(log)
	cmdRunner := system.NewExecCmdRunner(log)
	return &DnsUpdater{cfg: cfg, logger: log, fs: fs, cmdRunner: cmdRunner}
}

func (u *DnsUpdater) Run() {
	tickChan := time.Tick(dnsUpdateInterval)

	u.updateAll()

	for {
		select {
		case <-tickChan:
			err := u.updateAll()
			if err != nil {
				u.logger.Error(logTag, "Error updating DNS: %s, will try again in: %d seconds", err, dnsUpdateInterval)
			}
		}
	}
}

func (u *DnsUpdater) updateAll() (err error) {

	if !u.cfg.Master {
		return
	}

	for _, dnsServer := range u.cfg.DnsServers {
		err = u.updateDNSServer(u.cfg.Hostname, u.cfg.Ip, dnsServer, u.cfg.DnsKey, u.cfg.DnsTTL)
		if err != nil {
			return fmt.Errorf("Error updating dns server: %s, error: %s", dnsServer, err)
		}
	}

	return
}

func (u *DnsUpdater) updateDNSServer(nameToRegister, ip, dnsServer, dnsKey string, ttl int) (err error) {

	var tmpFile system.File
	if tmpFile, err = u.fs.TempFile("dnsUpdater-"); err != nil {
		return
	}
	defer u.fs.RemoveAll(tmpFile.Name())

	idx := strings.Index(nameToRegister, ".")
	zone := nameToRegister[idx+1:]

	configBody := fmt.Sprintf(
		dnsConfigTemplate,
		dnsServer,
		zone,
		nameToRegister,
		nameToRegister,
		ttl,
		ip,
	)

	if err = u.fs.WriteFileString(tmpFile.Name(), configBody); err != nil {
		return
	}

	_, _, _, err = u.cmdRunner.RunCommand("nsupdate", "-t", "4", "-y", dnsKey, "-v", tmpFile.Name())

	return
}

const dnsConfigTemplate = `
server %s
zone %s
update delete %s A
update add %s %d A %s
send

`

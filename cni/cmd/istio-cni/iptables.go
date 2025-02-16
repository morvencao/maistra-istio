// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This is a sample chained plugin that supports multiple CNI versions. It
// parses prevResult according to the cniVersion
package main

import (
	"fmt"
	"os/exec"
	"strings"

	"istio.io/pkg/log"
)

var (
	nsSetupProg = "istio-iptables"
)

type iptables struct {
}

func newIPTables() InterceptRuleMgr {
	return &iptables{}
}

// Program defines a method which programs iptables based on the parameters
// provided in Redirect.
func (ipt *iptables) Program(netns string, rdrct *Redirect) error {
	netnsArg := fmt.Sprintf("--net=%s", netns)
	nsSetupExecutable := fmt.Sprintf("%s/%s", nsSetupBinDir, nsSetupProg)
	nsenterArgs := []string{
		netnsArg,
		"--", // separate nsenter args from the rest with `--`, needed for hosts using BusyBox binaries
		nsSetupExecutable,
		"-p", rdrct.targetPort,
		"-u", rdrct.noRedirectUID,
		"-g", rdrct.noRedirectGID,
		"-m", rdrct.redirectMode,
		"-i", rdrct.includeIPCidrs,
		"-b", rdrct.includePorts,
		"-d", rdrct.excludeInboundPorts,
		"-o", rdrct.excludeOutboundPorts,
		"-x", rdrct.excludeIPCidrs,
		"-k", rdrct.kubevirtInterfaces,
	}
	if rdrct.redirectDNS {
		nsenterArgs = append(nsenterArgs, "--redirect-dns", "--capture-all-dns")
	}
	log.Infof("nsenter args: %s", strings.Join(nsenterArgs, " "))
	out, err := exec.Command("nsenter", nsenterArgs...).CombinedOutput()
	if err != nil {
		log.WithLabels("err", err, "out", out).Errorf("nsenter failed ")
		log.Infof("nsenter out: %s", out)
	} else {
		log.Infof("nsenter done: %s", out)
	}
	return err
}

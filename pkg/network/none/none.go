package none

import (
	"context"
	"os"
	"strconv"

	"github.com/rootless-containers/rootlesskit/pkg/api"
	"github.com/rootless-containers/rootlesskit/pkg/common"
	"github.com/rootless-containers/rootlesskit/pkg/network"
)

func NewParentDriver() (network.ParentDriver, error) {
	return &parentDriver{}, nil
}

type parentDriver struct {
}

const DriverName = "none"

func (d *parentDriver) MTU() int {
	return 0
}

func (d *parentDriver) Info(ctx context.Context) (*api.NetworkDriverInfo, error) {
	return &api.NetworkDriverInfo{
		Driver: DriverName,
	}, nil
}

func (d *parentDriver) ConfigureNetwork(childPID int, stateDir string) (*common.NetworkMessage, func() error, error) {
	var cleanups []func() error

	cmds := [][]string{
		[]string{"nsenter", "-t", strconv.Itoa(childPID), "--no-fork", "-n", "-m", "-U", "--preserve-credentials", "ip", "address", "add", "127.0.0.1/8", "dev", "lo"},
		[]string{"nsenter", "-t", strconv.Itoa(childPID), "--no-fork", "-n", "-m", "-U", "--preserve-credentials", "ip", "link", "set", "lo", "up"},
	}
	if err := common.Execs(os.Stderr, os.Environ(), cmds); err != nil {
		return nil, nil, err
	}

	netmsg := common.NetworkMessage{}
	return &netmsg, common.Seq(cleanups), nil
}

package grpcplugin

import (
	"github.com/1340691923/eve-plugin-sdk-go/backend/grpcplugin"
	"github.com/hashicorp/go-hclog"
	"os"

	goplugin "github.com/hashicorp/go-plugin"
	"os/exec"
)

type PluginDescriptor struct {
	pluginID         string
	executablePath   string
	executableArgs   []string
	versionedPlugins map[int]goplugin.PluginSet
	isDebug          bool
	pluginAddr       string
	pid              int
}

func newProductClientConfig(
	executablePath string,
	args []string,
	env []string,
	log hclog.Logger,
	versionedPlugins map[int]goplugin.PluginSet) *goplugin.ClientConfig {

	cmd := exec.Command(executablePath, args...)
	cmd.Env = env
	os.Chmod(executablePath, 0755)

	return &goplugin.ClientConfig{
		Cmd:              cmd,
		HandshakeConfig:  handshake,
		VersionedPlugins: versionedPlugins,
		Logger:           log,
		AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC},
	}
}

var handshake = goplugin.HandshakeConfig{
	ProtocolVersion:  grpcplugin.ProtocolVersion,
	MagicCookieKey:   grpcplugin.MagicCookieKey,
	MagicCookieValue: grpcplugin.MagicCookieValue,
}

type PluginAddr struct {
	Net     string
	Address string
}

func NewPluginAddr(net string, address string) *PluginAddr {
	return &PluginAddr{Net: net, Address: address}
}

func (p PluginAddr) Network() string {
	return p.Net
}

func (p PluginAddr) String() string {
	return p.Address
}

func newTestClientConfig(address string, pid int, log hclog.Logger, versionedPlugins map[int]goplugin.PluginSet) *goplugin.ClientConfig {

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   log.Name(),
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	return &goplugin.ClientConfig{
		Reattach: &goplugin.ReattachConfig{
			Addr:            NewPluginAddr("tcp", address),
			Pid:             pid,
			Protocol:        goplugin.ProtocolGRPC,
			Test:            true,
			ProtocolVersion: 2,
		},
		Logger:           logger,
		Plugins:          versionedPlugins[grpcplugin.ProtocolVersion],
		HandshakeConfig:  handshake,
		VersionedPlugins: versionedPlugins,
		AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC},
	}
}

func getV2PluginSet() goplugin.PluginSet {
	return goplugin.PluginSet{
		"resource": &grpcplugin.ResourceGRPCPlugin{},
		"basic":    &grpcplugin.PluginInfoGRPCPlugin{},
	}
}

func NewBackendPlugin(
	log hclog.Logger,
	closeLogWriteCallback func() error,
	pluginID,
	executablePath string,
	isDebug bool,
	testPluginTcpAddr string,
	testPluginPid int,
	executableArgs []string) *GrpcPlugin {
	return newPlugin(PluginDescriptor{
		pluginID:       pluginID,
		executablePath: executablePath,
		executableArgs: executableArgs,
		pluginAddr:     testPluginTcpAddr,
		pid:            testPluginPid,
		isDebug:        isDebug,
		versionedPlugins: map[int]goplugin.PluginSet{
			grpcplugin.ProtocolVersion: getV2PluginSet(),
		},
	}, []string{}, log, closeLogWriteCallback)
}

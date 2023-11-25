package main

import (
	"github.com/gliderlabs/ssh"
	"iu9-networks/lab4/configs"
	"iu9-networks/lab4/pkg/core"
)

func main() {
	var (
		p              = core.SetupParameters("")
		forwardHandler = &ssh.ForwardedTCPHandler{}
		server         = ssh.Server{
			Handler:                       core.CreateSSHSessionHandler(p.Shell),
			PasswordHandler:               core.CreatePasswordHandler(configs.LocalPassword),
			PublicKeyHandler:              core.CreatePublicKeyHandler(configs.AuthorizedKey),
			LocalPortForwardingCallback:   core.CreateLocalPortForwardingCallback(p.NoShell),
			ReversePortForwardingCallback: core.CreateReversePortForwardingCallback(),
			SessionRequestCallback:        core.CreateSessionRequestCallback(p.NoShell),
			ChannelHandlers: map[string]ssh.ChannelHandler{
				"direct-tcpip": ssh.DirectTCPIPHandler,
				"session":      ssh.DefaultSessionHandler,
				"rs-info":      core.CreateExtraInfoHandler(),
			},
			RequestHandlers: map[string]ssh.RequestHandler{
				"tcpip-forward":        forwardHandler.HandleSSHRequest,
				"cancel-tcpip-forward": forwardHandler.HandleSSHRequest,
			},
			SubsystemHandlers: map[string]ssh.SubsystemHandler{
				"sftp": core.CreateSFTPHandler(),
			},
		}
	)

	core.Run(p, server)
}

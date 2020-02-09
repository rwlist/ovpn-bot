package app

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"os/exec"
	"strings"
)

const (
	prefix = "ovpn_"

	data = "data"
	tcp = "tcp"
	udp = "udp"
)

type Logic struct {
	cli *client.Client
}

func NewLogic(cli *client.Client) *Logic {
	return &Logic{
		cli: cli,
	}
}

type readCloser struct{}

func (r readCloser) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (r readCloser) Close() error {
	return nil
}

/*
	docker volume create --name $OVPN_DATA
	docker run -v $OVPN_DATA:/etc/openvpn --rm kylemanna/openvpn ovpn_genconfig -u $(PROTO)://$(HOST):$(PORT)
	docker run -v $OVPN_DATA:/etc/openvpn --rm -it kylemanna/openvpn ovpn_initpki
 */
func (l *Logic) CommandInit(w io.Writer) error {
	dataVolume := prefix + data

	err := l.execute(w, []string{"docker", "volume", "create", "--name", dataVolume})
	if err != nil {
		return err
	}

	//dock := command.NewDockerCli(readCloser{}, w, w)
	//run := container.NewRunCommand(dock)
	//args := []string{"--name", dataName, "-v", "/etc/openvpn", "busybox"}
	//run.SetArgs(args)
	//err := run.RunE(run, args)

	_, _ = fmt.Fprintf(w, "All done, init completed!")
	return nil
}

func (l *Logic) CommandStatus() (string, error) {
	list, err := l.cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return "", err
	}

	var ovpnContainers []types.Container

	text := fmt.Sprintf("Total %v contrainers:\n", len(list))
	for _, c := range list {
		text += "\n"

		name := strings.Join(c.Names, ":")
		if strings.HasPrefix(name, prefix) {
			ovpnContainers = append(ovpnContainers, c)
		}

		text += formatContainer(c)
	}

	text += fmt.Sprintf("\n\nTotal %v ovpn containers:\n", len(ovpnContainers))
	for _, c := range ovpnContainers {
		text += "\n"
		text += formatContainer(c)
	}

	return text, nil
}

func (l *Logic) CommandGenerate(profileName string) (string, error) {
	return "ok", nil
}

func (l *Logic) execute(w io.Writer, args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = w
	cmd.Stderr = w

	_, _ = fmt.Fprintf(w, "Executing command: `%s`", strings.Join(args, " "))

	err := cmd.Run()
	return err
}
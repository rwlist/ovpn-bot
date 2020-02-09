package app

import (
	"bytes"
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

func (l *Logic) CommandInit(w io.Writer, addr string) error {
	dataVolume := prefix + data

	_, _, port, ok := parseAddr(addr)
	if !ok {
		return fmt.Errorf(`"%s" is not valid addr`, addr)
	}

	// docker volume create --name $OVPN_DATA
	err := l.execute(w, []string{"docker", "volume", "create", "--name", dataVolume})
	if err != nil {
		return err
	}

	dataMount := dataVolume+":/etc/openvpn"

	// docker run -v $OVPN_DATA:/etc/openvpn --rm kylemanna/openvpn ovpn_genconfig -u $(PROTO)://$(HOST):$(PORT)
	err = l.execute(w, []string{"docker", "run", "-v", dataMount, "--rm", "kylemanna/openvpn", "ovpn_genconfig", "-u", addr})
	if err != nil {
		return err
	}

	// docker run -v $OVPN_DATA:/etc/openvpn --rm -it kylemanna/openvpn ovpn_initpki
	lnReader := bytes.NewReader([]byte{'\n', '\n'})
	err = l.execute2(w, []string{"docker", "run", "-v", dataMount, "--rm", "-i", "kylemanna/openvpn", "ovpn_initpki", "nopass"}, lnReader)
	if err != nil {
		return err
	}

	udpContainer := prefix + udp
	tcpContainer := prefix + tcp

	// docker run -v $OVPN_DATA:/etc/openvpn -d --restart=always --name $(NAME)_udp -p $(PORT):1194/udp --cap-add=NET_ADMIN kylemanna/openvpn ovpn_run --proto udp
	err = l.execute(w, []string{"docker", "run", "-v", dataMount, "-d", "--restart=always", "--name", udpContainer, "-p", port+":1194/udp", "--cap-add=NET_ADMIN", "kylemanna/openvpn", "ovpn_run", "--proto", "udp"})
	if err != nil {
		return err
	}

	// docker run -v $OVPN_DATA:/etc/openvpn -d --restart=always --name $(NAME)_tcp -p $(PORT):1194/tcp --cap-add=NET_ADMIN kylemanna/openvpn ovpn_run --proto tcp
	err = l.execute(w, []string{"docker", "run", "-v", dataMount, "-d", "--restart=always", "--name", tcpContainer, "-p", port+":1194/tcp", "--cap-add=NET_ADMIN", "kylemanna/openvpn", "ovpn_run", "--proto", "tcp"})
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(w, "All done, init completed!")
	return nil
}

func (l *Logic) CommandRemove(w io.Writer) {
	dataVolume := prefix + data
	udpContainer := prefix + udp
	tcpContainer := prefix + tcp

	var err error

	err = l.execute(w, []string{"docker", "rm", "-f", udpContainer})
	if err != nil {
		_, _ = fmt.Fprintf(w, "remove error: %v\n", err)
	}

	err = l.execute(w, []string{"docker", "rm", "-f", tcpContainer})
	if err != nil {
		_, _ = fmt.Fprintf(w, "remove error: %v\n", err)
	}

	err = l.execute(w, []string{"docker", "volume", "rm", dataVolume})
	if err != nil {
		_, _ = fmt.Fprintf(w, "remove error: %v\n", err)
	}

	_, _ = fmt.Fprintf(w, "All removed!")
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

func (l *Logic) CommandGenerate(w *botWriter, profileName string) ([]byte, error) {
	dataVolume := prefix + data
	dataMount := dataVolume+":/etc/openvpn"

	// docker run -v ovpn_data:/etc/openvpn --rm -i kylemanna/openvpn easyrsa build-client-full client_name nopass
	err := l.execute(w, []string{"docker", "run", "-v", dataMount, "--rm", "-i", "kylemanna/openvpn", "easyrsa", "build-client-full", profileName, "nopass"})
	if err != nil {
		return nil, err
	}

	// docker run -v ovpn_data:/etc/openvpn --rm kylemanna/openvpn ovpn_getclient client_name
	configWriter := bytes.NewBuffer(nil)
	err = l.execute(configWriter, []string{"docker", "run", "-v", dataMount, "--rm", "kylemanna/openvpn", "ovpn_getclient", profileName})
	if err != nil {
		return nil, err
	}

	configData := configWriter.Bytes()
	return configData, nil
}

func (l *Logic) execute(w io.Writer, args []string) error {
	return l.execute2(w, args, nil)
}

func (l *Logic) execute2(w io.Writer, args []string, stdin io.Reader) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = w
	cmd.Stderr = w

	_, _ = fmt.Fprintf(w, "Executing command: `%s`", strings.Join(args, " "))

	err := cmd.Run()
	return err
}
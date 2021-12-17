package src

import (
	"bytes"
	"context"
	"errors"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"os"
)

type remoteScriptType byte

const (
	cmdLine remoteScriptType = iota
	rawScript
	scriptFile
)

type ViaSSHDialer struct {
	client *ssh.Client
	_      *context.Context
}

type Client struct {
	client *ssh.Client
}

type remoteScript struct {
	client *ssh.Client
	_type  remoteScriptType
	script *bytes.Buffer
	scriptFile string
	err        error
	stdout     io.Writer
	stderr     io.Writer
}

type remoteShell struct {
	client         *ssh.Client
	requestPty     bool
	terminalConfig *TerminalConfig

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

type TerminalConfig struct {
	Term   string
	Height int
	Weight int
	Modes  ssh.TerminalModes
}

type Config struct {
	Addr, User, Passwd string
}

//Dial -- ----------------------------
//--> @Description get dial
//--> @Param
//--> @return
//-- ----------------------------
func (v *ViaSSHDialer) Dial(context context.Context, addr string) (net.Conn, error) {
	return v.client.Dial("tcp", addr)
}

//DialWithPasswd -- --------------------------------------------------------------------------
//--> @Description starts a client connection to the given SSH server with passwd auth method.
//--> @Param
//--> @return
//-- ----------------------------
func DialWithPasswd(cfg *Config) (*Client, error) {
	if cfg.User == "" || cfg.Passwd == "" || cfg.Addr == "" {
		return nil, errors.New("config info is empty of one")
	}
	config := &ssh.ClientConfig{
		User: cfg.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.Passwd),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return Dial("tcp", cfg.Addr, config)
}

//GetClient -- ----------------------------
//--> @Description get sesssion client
//--> @Param
//--> @return
//-- ----------------------------
func (c *Client) GetClient() *ssh.Client {
	return c.client
}

//DialWithKey -- -------------------------------------------------------------------------
//--> @Description starts a client connection to the given SSH server with key auth method.
//--> @Param
//--> @return
//-- ----------------------------
func DialWithKey(addr, user, keyfile string) (*Client, error) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return Dial("tcp", addr, config)
}

//DialWithKeyWithPassphrase -- ------------------------------------------------------------
//--> @Description same as DialWithKey but with a passphrase to decrypt the private key
//--> @Param
//--> @return
//-- ----------------------------
func DialWithKeyWithPassphrase(addr, user, keyfile string, passphrase string) (*Client, error) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(passphrase))
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}
	return Dial("tcp", addr, config)
}

//Dial -- ----------------------------
//--> @Description Dial starts a client connection to the given SSH server.
//--> @Description This is wrap the ssh.Dial
//--> @Param
//--> @return
//-- ----------------------------
func Dial(network, addr string, config *ssh.ClientConfig) (*Client, error) {
	client, err := ssh.Dial(network, addr, config)

	if err != nil {
		return nil, err
	}
	return &Client{
		client: client,
	}, nil
}

//Close -- ----------------------------
//--> @Description close result
//--> @Param
//--> @return
//-- ----------------------------
func (c *Client) Close() error {
	return c.client.Close()
}

//Cmd -- ----------------------------
//--> @Description create a command on client
//--> @Param
//--> @return
//-- ----------------------------
func (c *Client) Cmd(cmd string) *remoteScript {
	return &remoteScript{
		_type:  cmdLine,
		client: c.client,
		script: bytes.NewBufferString(cmd + "\n"),
	}
}

//Script -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func (c *Client) Script(script string) *remoteScript {
	return &remoteScript{
		_type:  rawScript,
		client: c.client,
		script: bytes.NewBufferString(script + "\n"),
	}
}

//ScriptFile -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func (c *Client) ScriptFile(fname string) *remoteScript {
	return &remoteScript{
		_type:      scriptFile,
		client:     c.client,
		scriptFile: fname,
	}
}

//Run -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func (rs *remoteScript) Run() error {
	if rs.err != nil {
		return rs.err
	}
	if rs._type == cmdLine {
		return rs.runCmdS()
	} else if rs._type == rawScript {
		return rs.runScript()
	} else if rs._type == scriptFile {
		return rs.runScriptFile()
	} else {
		return errors.New("not supported remoteScript type")
	}
}

//Output -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func (rs *remoteScript) Output() ([]byte, error) {
	if rs.stdout != nil {
		return nil, errors.New("stdout already set")
	}
	var out bytes.Buffer
	rs.stdout = &out
	err := rs.Run()
	return out.Bytes(), err
}

//SmartOutput -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func (rs *remoteScript) SmartOutput() ([]byte, error) {
	if rs.stdout != nil {
		return nil, errors.New("stdout already set")
	}
	if rs.stderr != nil {
		return nil, errors.New("stderr already set")
	}

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	rs.stdout = &stdout
	rs.stderr = &stderr
	err := rs.Run()
	if err != nil {
		return stderr.Bytes(), err
	}
	return stdout.Bytes(), err
}

//Cmd -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func (rs *remoteScript) Cmd(cmd string) *remoteScript {
	_, err := rs.script.WriteString(cmd + "\n")
	if err != nil {
		rs.err = err
	}
	return rs
}

//SetStdio -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func (rs *remoteScript) SetStdio(stdout, stderr io.Writer) *remoteScript {
	rs.stdout = stdout
	rs.stderr = stderr
	return rs
}

//runCmd -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func (rs *remoteScript) runCmd(cmd string) error {
	session, err := rs.client.NewSession()
	if err != nil {
		return err
	}
	defer func() {
		err = session.Close()
	}()

	session.Stdout = rs.stdout
	session.Stderr = rs.stderr

	if err := session.Run(cmd); err != nil {
		return err
	}
	return nil
}

//runCmdS -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func (rs *remoteScript) runCmdS() error {
	for {
		statment, err := rs.script.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := rs.runCmd(statment); err != nil {
			return err
		}
	}

	return nil
}

//runScript -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func (rs *remoteScript) runScript() error {
	session, err := rs.client.NewSession()
	if err != nil {
		return err
	}

	session.Stdin = rs.script
	session.Stdout = rs.stdout
	session.Stderr = rs.stderr
	if err := session.Shell(); err != nil {
		return err
	}
	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}

//runScriptFile -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func (rs *remoteScript) runScriptFile() error {
	var buffer bytes.Buffer
	file, err := os.Open(rs.scriptFile)
	if err != nil {
		return err
	}
	_, err = io.Copy(&buffer, file)
	if err != nil {
		return err
	}

	rs.script = &buffer
	return rs.runScript()
}

//Terminal -- ----------------------------
//--> @Description create a interactive shell on client.
//--> @Param
//--> @return
//-- ----------------------------
func (c *Client) Terminal(config *TerminalConfig) *remoteShell {
	return &remoteShell{
		client:         c.client,
		terminalConfig: config,
		requestPty:     true,
	}
}

//Shell -- ----------------------------
//--> @Description create a noninteractive shell on client.
//--> @Param
//--> @return
//-- ----------------------------
func (c *Client) Shell() *remoteShell {
	return &remoteShell{
		client:     c.client,
		requestPty: false,
	}
}

//SetStdio -- ----------------------------
//--> @Description SetStdio
//--> @Param
//--> @return
//-- ----------------------------
func (rs *remoteShell) SetStdio(stdin io.Reader, stdout, stderr io.Writer) *remoteShell {
	rs.stdin = stdin
	rs.stdout = stdout
	rs.stderr = stderr
	return rs
}

//Start -- ----------------------------
//--> @Description Start
//--> @Param
//--> @return
//-- ----------------------------
func (rs *remoteShell) Start() (err error) {
	session, err := rs.client.NewSession()
	if err != nil {
		return err
	}
	defer func() {
		err = session.Close()
	}()

	if rs.stdin == nil {
		session.Stdin = os.Stdin
	} else {
		session.Stdin = rs.stdin
	}
	if rs.stdout == nil {
		session.Stdout = os.Stdout
	} else {
		session.Stdout = rs.stdout
	}
	if rs.stderr == nil {
		session.Stderr = os.Stderr
	} else {
		session.Stderr = rs.stderr
	}

	if rs.requestPty {
		tc := rs.terminalConfig
		if tc == nil {
			tc = &TerminalConfig{
				Term:   "xterm",
				Height: 40,
				Weight: 80,
			}
		}
		if err := session.RequestPty(tc.Term, tc.Height, tc.Weight, tc.Modes); err != nil {
			return err
		}
	}

	if err := session.Shell(); err != nil {
		return err
	}

	if err := session.Wait(); err != nil {
		return err
	}
	return nil
}

//TestSShServer -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func TestSShServer(s *ssh.Session, client *Client, shell string) (str string, err error) {
	var out []byte
	//测试是否连接上
	if out, err = client.Cmd(shell).Output(); err != nil {
		panic(err)
	}
	return string(out), err
}

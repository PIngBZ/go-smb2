package smb2_test

import (
	"fmt"
	"github.com/PIngBZ/ntlmssp"
	"io"
	"io/ioutil"
	"net"

	"github.com/PIngBZ/go-smb2"
)

func Example() {
	conn, err := net.Dial("tcp", "localhost:445")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ntlmsspClient, err := ntlmssp.NewClient(
		ntlmssp.SetCompatibilityLevel(3),
		ntlmssp.SetUserInfo("Guest", ""),
		ntlmssp.SetDomain("MicrosoftAccount"))
	if err != nil {
		panic(err)
	}
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMSSPInitiator{
			NTLMSSPClient: ntlmsspClient,
		},
	}

	c, err := d.Dial(conn)
	if err != nil {
		panic(err)
	}
	defer c.Logoff()

	fs, err := c.Mount(`\\localhost\share`)
	if err != nil {
		panic(err)
	}
	defer fs.Umount()

	f, err := fs.Create("hello.txt")
	if err != nil {
		panic(err)
	}
	defer fs.Remove("hello.txt")
	defer f.Close()

	_, err = f.Write([]byte("Hello world!"))
	if err != nil {
		panic(err)
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bs))

	// Hello world!
}

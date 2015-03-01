package main

import (
	"fmt"
	//"io"
	//"os"
	"github.com/fsouza/go-dockerclient"
	"github.com/samalba/dockerclient"
	"net"
	"time"
	"io"
	"syscall"
	log "github.com/Sirupsen/logrus"
)

var dockerClient *docker.Client
var mydocker *dockerclient.DockerClient

var running bool

var containerID string

var back  *net.TCPAddr

func main() {
	containerID = "699dd412a994"

	//go-dockerclient
    endpoint := "unix:///var/run/docker.sock"
    dockerClient, _ = docker.NewClient(endpoint)

    // Init the client
    mydocker, _ = dockerclient.NewDockerClient(endpoint, nil)

    c, _ := dockerClient.InspectContainer(containerID)

    fmt.Println("", c.State.Running)
    fmt.Println("", c.HostConfig.Binds)

    if c.State.Running {
    	dockerClient.KillContainer(docker.KillContainerOptions{
    		ID: containerID,
    		})
    }

    running = false

    go func() {
    	for {
    		time.Sleep(10*time.Second)
    		if running {
    			dockerClient.KillContainer(docker.KillContainerOptions{
    		ID: containerID,
    		})
    			running = false
    			fmt.Printf("KILL")
    		}
    	}
    }()

	front,_ := net.ResolveTCPAddr("tcp", "127.0.0.1:6666")
	back,_ = net.ResolveTCPAddr("tcp", "172.17.0.167:9000")

	t, _ := NewTCPProxy(front, back)
	t.Run()
}

type TCPProxy struct {
	listener     *net.TCPListener
	frontendAddr *net.TCPAddr
	backendAddr  *net.TCPAddr
}

func NewTCPProxy(frontendAddr, backendAddr *net.TCPAddr) (*TCPProxy, error) {
	listener, err := net.ListenTCP("tcp", frontendAddr)
	if err != nil {
		return nil, err
	}
	// If the port in frontendAddr was 0 then ListenTCP will have a picked
	// a port to listen on, hence the call to Addr to get that actual port:
	return &TCPProxy{
		listener:     listener,
		frontendAddr: listener.Addr().(*net.TCPAddr),
		backendAddr:  backendAddr,
	}, nil
}

func (proxy *TCPProxy) clientLoop(client *net.TCPConn, quit chan bool) {
	backend, err := net.DialTCP("tcp", nil, proxy.backendAddr)
	if err != nil {
		log.Printf("Can't forward traffic to backend tcp/%v: %s\n", proxy.backendAddr, err)
		client.Close()
		return
	}

	event := make(chan int64)
	var broker = func(to, from *net.TCPConn) {
		written, err := io.Copy(to, from)
		if err != nil {
			// If the socket we are writing to is shutdown with
			// SHUT_WR, forward it to the other end of the pipe:
			if err, ok := err.(*net.OpError); ok && err.Err == syscall.EPIPE {
				from.CloseWrite()
			}
		}
		to.CloseRead()
		event <- written
	}

	go broker(client, backend)
	go broker(backend, client)

	var transferred int64 = 0
	for i := 0; i < 2; i++ {
		select {
		case written := <-event:
			transferred += written
		case <-quit:
			// Interrupt the two brokers and "join" them.
			client.Close()
			backend.Close()
			for ; i < 2; i++ {
				transferred += <-event
			}
			return
		}
	}
	client.Close()
	backend.Close()
}

func (proxy *TCPProxy) Run() {
	quit := make(chan bool)
	defer close(quit)
	for {
		client, err := proxy.listener.Accept()
		if err != nil {
			log.Printf("Stopping proxy on tcp/%v for tcp/%v (%s)", proxy.frontendAddr, proxy.backendAddr, err)
			return
		}

		if !running {
			mydocker.StartContainer(containerID, nil)
			//time.Sleep(time.Millisecond*100)

			c, _ := dockerClient.InspectContainer(containerID)

    		fmt.Println("", c.NetworkSettings.IPAddress)

    		proxy.backendAddr, _ = net.ResolveTCPAddr("tcp", c.NetworkSettings.IPAddress+":9000")

    		fmt.Println("NEW START")

    		//wait until app starts serving requests
		for {
			_, err := net.DialTCP("tcp", nil, proxy.backendAddr)

			if err == nil {
				break
			}

			time.Sleep(time.Millisecond * 5)
		}

			running=true
		}else{
			fmt.Println("EXISTING")
		}

		


		go proxy.clientLoop(client.(*net.TCPConn), quit)
	}
}

func (proxy *TCPProxy) Close()                 { proxy.listener.Close() }
func (proxy *TCPProxy) FrontendAddr() net.Addr { return proxy.frontendAddr }
func (proxy *TCPProxy) BackendAddr() net.Addr  { return proxy.backendAddr }
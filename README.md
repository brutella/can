# can

*can* provides an interface to a [CAN bus](https://www.kernel.org/doc/Documentation/networking/can.txt) to read and write frames. The library is based on the [SocketCAN](https://github.com/torvalds/linux/blob/097f70b3c4d84ffccca15195bdfde3a37c0a7c0f/include/uapi/linux/can.h) network stack on Linux.

# Hardware

I'm using a [Raspberry Pi 2 Model B](https://www.raspberrypi.org/products/raspberry-pi-2-model-b/) and [PiCAN2 CAN-Bus board for Raspberry Pi 2](http://skpang.co.uk/catalog/pican2-canbus-board-for-raspberry-pi-2-p-1475.html) to connect to a CAN bus.

# Software

The Raspberry Pi runs [Raspbian](https://www.raspberrypi.org/downloads/raspbian/). 

# Configuration

Update `/boot/config.txt` with

    dtparam=spi=on 
    dtoverlay=mcp2515-can0-overlay,oscillator=16000000,interrupt=25 
    dtoverlay=spi-bcm2835-overlay

and reboot.

After phyiscally connecting to the CAN bus, you have to set up the can network interface for a specific bitrate, i.e. 50 kB

    sudo ip link set can0 up type can bitrate 50000

Running `ifconfig` should now include the `can0` interface. 

#### Test your configuration

You should test if you actually receive data from the CAN bus. You can either use the `candump` tool from the [can-utils](https://github.com/linux-can/can-utils) or a simple reimplementation under `cmd/candump.go`. 

Either way you will see something like this

    > go run $GOSRC/github.com/brutella/can/cmd/candump.go -if can0
    
    can0 100  [6] 20 83 0C 00 67 29        ' ...g)'
    can0 701  [1] 05                       '.'

## Usage

#### Setup the CAN bus

```go
bus, _ := can.NewBusForInterfaceWithName("can0")
bus.ConnectAndPublish()
```

#### Send a CAN frame

```go
frm := can.Frame{
	ID:     0x701,
	Length: 1,
	Flags:  0,
	Res0:   0,
	Res1:   0,
	Data:   [8]uint8{0x05},
}

bus.Publish(frm)
```
    
#### Receive a CAN frame

```go
bus.SubscribeFunc(handleCANFrame)

func handleCANFrame(frm can.Frame) {    
    ...
}
```

There is more to learn from the [documentation](http://godoc.org/github.com/brutella/can).

# Contact

Matthias Hochgatterer

Github: [https://github.com/brutella](https://github.com/brutella/)

Twitter: [https://twitter.com/brutella](https://twitter.com/brutella)


# License

*can* is available under the MIT license. See the LICENSE file for more info.
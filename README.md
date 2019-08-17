# can

*can* provides an interface to a [CAN bus](https://www.kernel.org/doc/Documentation/networking/can.txt) to read and write frames. The library is based on the [SocketCAN](https://github.com/torvalds/linux/blob/097f70b3c4d84ffccca15195bdfde3a37c0a7c0f/include/uapi/linux/can.h) network stack on Linux.

# Usage

## Setup the CAN bus

```go
bus, _ := can.NewBusForInterfaceWithName("can0")
bus.ConnectAndPublish()
```

## Send a CAN frame

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
    
## Receive a CAN frame

```go
bus.SubscribeFunc(handleCANFrame)

func handleCANFrame(frm can.Frame) {    
    ...
}
```

There is more to learn from the [documentation](http://godoc.org/github.com/sitec-systems/can).

# License

The orignal work is done by [brutella](https://github.com/brutella). This project is a fork of https://github.com/brutella/can

*can* is available under the MIT license. See the LICENSE file for more info.

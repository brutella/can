// Package can provides a CAN bus to send and receive CAN frames.
package can

const (
	MaskID    = 0x000007FF // 11-bit identifier
	MaskIDEff = 0x1FFFFFFF // 29-bit identifier
)

const (
	FlagErr = 0x20000000 // 0 = data frame, 1 = error message
	FlagRtr = 0x40000000 // 1 = rtr frame
	FlagEff = 0x80000000 // 0 = standard 11 bit, 1 = extended 29 bit
)

const (
	MaxFrameDataLength    = 8  // ISO 11898-1
	MaxExtFrameDataLength = 64 // ISO 11898-7
)

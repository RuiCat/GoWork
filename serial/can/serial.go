package can

import "serial"

// NewSerialBus returns a new CAN bus.
func NewSerialBus(config *serial.Config) (*Bus, error) {
	p, err := serial.Open(config)
	if err != nil {
		return nil, err
	}
	return NewBus(NewReadWriteCloser(p)), err
}

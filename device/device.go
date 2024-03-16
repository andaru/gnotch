package device

// Device is the Gnotch base device interface
type Device interface {
	// // Name returns the device name
	// Name() string
	// Connected returns true if the device has a network connection underway
	Connected() bool
	// Close closes the device. Shadows io.Closer
	Close() error
}

// Commander is the interface covering the most basic device functionality, command execution.
//
// A concrete device type implementing Commander will typically also implement Device.
type Commander interface {
	// Command executes a command on the device, returning the result
	Command(command string) ([]byte, error)
}

// Provider is the interface to retrieving device information.
//
// A device provider provides a mapping between devices and their vendor.
type Provider interface {
	Vendor(deviceName string) (vendorName string)
}

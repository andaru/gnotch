package routerdb

import (
	"bufio"
	"io"
	"strings"
)

// RouterDB is a device provider for RANCID router.db files.
//
// Canonically, a router.db file is a semi-colon field separated database file.
// Three fields are defined by RANCID, those being 1) hostname, 2) vendor_string and
// 3) "up" or "down". Users may specify additional fields for their own purposes.
type RouterDB struct {
	// db is the map from device name to vendor string
	db map[string]string

	// allow "down" devices to be addressed
	allowDown bool
}

// Option describes a router.db file option
type Option func(*RouterDB)

// WithAllowDown permits "down" devices in the router.db file to be used by gnotch
func WithAllowDown() Option { return func(db *RouterDB) { db.allowDown = true } }

// New returns a new router.db file device.Provider, loading data from r.
//
// opts can be used to additionally configure the device provider.
func New(r io.Reader, opts ...Option) (*RouterDB, error) {
	db := &RouterDB{db: map[string]string{}}
	for _, opt := range opts {
		opt(db)
	}
	err := db.load(r)
	return db, err
}

// Vendor returns the vendor identifier for the given device name
func (db *RouterDB) Vendor(deviceName string) string { return db.db[deviceName] }

// load loads the router.db file data from r, returning any read errors
func (db *RouterDB) load(r io.Reader) error {
	var deviceName, vendor, updown string
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		switch f := strings.SplitN(scan.Text(), ":", 4); len(f) {
		case 4:
			deviceName, vendor, updown, _ = f[0], f[1], strings.ToLower(f[2]), f[3]
		case 3:
			deviceName, vendor, updown = f[0], f[1], strings.ToLower(f[2])
		default:
			continue
		}
		if updown == "up" || (updown == "down" && db.allowDown) {
			db.db[deviceName] = vendor
		}
	}
	return scan.Err()
}

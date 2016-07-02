package mysql

import (
	"errors"
	"fmt"
)

var (
	errInvalidConn = errors.New("Invalid Connection")
	errMalformPkt  = errors.New("Malformed Packet")
	errNoTLS       = errors.New("TLS encryption requested but server does not support TLS")
	errOldPassword = errors.New("This server only supports the insecure old password authentication. If you still want to use it, please add 'allowOldPasswords=1' to your DSN. See also https://github.com/go-sql-driver/mysql/wiki/old_passwords")
	errOldProtocol = errors.New("MySQL-Server does not support required Protocol 41+")
	errPktSync     = errors.New("Commands out of sync. You can't run this command now")
	errPktSyncMul  = errors.New("Commands out of sync. Did you run multiple statements at once?")
	errPktTooLarge = errors.New("Packet for query is too large. You can change this value on the server by adjusting the 'max_allowed_packet' variable.")
)

// error type which represents a single MySQL error
type MySQLError struct {
	Number  uint16
	Message string
}

func (me *MySQLError) Error() string {
	return fmt.Sprintf("Error %d: %s", me.Number, me.Message)
}

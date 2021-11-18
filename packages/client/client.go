package client

import (
	"crypto/tls"
	"errors"
	"time"

	"github.com/rs/zerolog"
)

type Status int8

type Client struct {
	host   string
	logger *zerolog.Logger
}

var (
	ErrConnection  error = errors.New("connection error")
	ErrTimeout     error = errors.New("connection timeout")
	ErrCertInvalid error = errors.New("invalid certificate")
	ErrCertExpired error = errors.New("certificate has been expired")
)

func New(hostname string, zerolog *zerolog.Logger) *Client {
	return &Client{
		host:   hostname,
		logger: zerolog,
	}
}

func (c *Client) GetExpiry() (*time.Time, error) {
	protocol := "tcp"
	hostPort := c.host + ":443"
	c.logInfo("Establishing TLS connection to " + protocol + "://" + hostPort)

	tlsConn, err := tls.Dial(protocol, hostPort, nil)
	if err != nil {
		c.logWarn(err)
		return nil, ErrConnection
	}
	defer tlsConn.Close()

	if err = tlsConn.VerifyHostname(c.host); err != nil {
		c.logWarn(err)
		return nil, ErrCertInvalid
	}

	expiry := tlsConn.ConnectionState().PeerCertificates[0].NotAfter
	if time.Until(expiry) <= 0 {
		c.logWarn(errors.New("certificate has been expired"))
		return &expiry, ErrCertExpired
	}

	c.logInfo("Operation succeeded")
	return &expiry, nil
}

func (c *Client) logWarn(err error) {
	if c.logger != nil {
		c.logger.Warn().Msgf("%s: %s", c.host, err.Error())
	}
}

func (c *Client) logInfo(info string) {
	if c.logger != nil {
		c.logger.Info().Msgf("%s: %s", c.host, info)
	}
}

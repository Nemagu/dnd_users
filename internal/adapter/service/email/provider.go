package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/smtp"
	"os"
	"sync"
	"time"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
)

type SMTPEmailProvider struct {
	mu       sync.RWMutex
	logger   *slog.Logger
	username string
	password string
	from     string
	host     string
	port     int
	timeout  time.Duration
}

func NewEmailProvider(
	logger *slog.Logger,
	username, password, from, host string,
	port int,
	timeout time.Duration,
) (*SMTPEmailProvider, error) {
	return &SMTPEmailProvider{
		logger:   logger,
		username: username,
		password: password,
		from:     from,
		host:     host,
		port:     port,
		timeout:  timeout,
	}, nil
}

func MustNewEmailProvider(
	logger *slog.Logger,
	username, password, from, host string,
	port int,
	timeout time.Duration,
) *SMTPEmailProvider {
	return &SMTPEmailProvider{
		logger:   logger,
		username: username,
		password: password,
		from:     from,
		host:     host,
		port:     port,
		timeout:  timeout,
	}
}

func (p *SMTPEmailProvider) SendConfirmEmail(
	message appdto.Email,
) {
	ctx, cancel := p.createContext()
	defer cancel()
	p.sendEmail(ctx, message.To, buildConfirmMessage(message))
}

func (p *SMTPEmailProvider) SendChangeEmail(
	message appdto.Email,
) {
	ctx, cancel := p.createContext()
	defer cancel()
	p.sendEmail(ctx, message.To, buildChangeEmailMessage(message))
}

func (p *SMTPEmailProvider) SendResetPasswordEmail(
	message appdto.Email,
) {
	ctx, cancel := p.createContext()
	defer cancel()
	p.sendEmail(ctx, message.To, buildResetPasswordMessage(message))
}

func (p *SMTPEmailProvider) createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), p.timeout)
}

func (p *SMTPEmailProvider) sendEmail(ctx context.Context, to string, message []byte) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	p.logger.InfoContext(ctx, "try to create connection")
	conn, err := p.createConnection()
	if err != nil {
		p.logger.ErrorContext(ctx, "connection failed", "error", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			p.logger.ErrorContext(ctx, "connection closing failed", "error", err)
		}
	}()
	p.logger.InfoContext(ctx, "try to create client")
	client, err := smtp.NewClient(conn, p.host)
	if err != nil {
		p.logger.ErrorContext(ctx, "creation client failed", "error", err)
	}
	defer func() {
		if err := client.Quit(); err != nil {
			p.logger.ErrorContext(ctx, "quit client failed", "error", err)
		}
	}()
	auth := smtp.PlainAuth("", p.username, p.password, p.host)
	p.logger.InfoContext(ctx, "try to authenticate")
	if err := client.Auth(auth); err != nil {
		p.logger.ErrorContext(ctx, "authentication failed", "error", err)
	}
	p.logger.InfoContext(ctx, "try to set sender")
	if err = client.Mail(p.from); err != nil {
		p.logger.ErrorContext(ctx, "sender setting failed", "error", err)
	}
	p.logger.InfoContext(ctx, "try to set recipients")
	if err = client.Rcpt(to); err != nil {
		p.logger.ErrorContext(ctx, "recipients setting failed", "error", err)
	}
	w, err := client.Data()
	if err != nil {
		p.logger.ErrorContext(ctx, "writer creation failed", "error", err)
	}
	defer func() {
		if err := w.Close(); err != nil {
			p.logger.ErrorContext(ctx, "writer closing failed", "error", err)
		}
	}()
	p.logger.InfoContext(ctx, "try to write email")
	if _, err = w.Write(message); err != nil {
		p.logger.ErrorContext(ctx, "message writing failed", "error", err)
	}
}

func (p *SMTPEmailProvider) createConnection() (*tls.Conn, error) {
	tlsConfig := &tls.Config{
		ServerName: p.host,
	}

	addr := fmt.Sprintf("%s:%d", p.host, p.port)

	return tls.Dial("tcp", addr, tlsConfig)
}

type FileEmailProvider struct {
	logger     *slog.Logger
	folderPath string
	timeout    time.Duration
}

func NewFileEmailProvider(
	logger *slog.Logger, folderPath string, timeout time.Duration,
) (*FileEmailProvider, error) {
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	return &FileEmailProvider{
		logger:     logger,
		folderPath: folderPath,
		timeout:    timeout,
	}, nil
}

func MustNewFileEmailProvider(
	logger *slog.Logger, folderPath string, timeout time.Duration,
) *FileEmailProvider {
	if provider, err := NewFileEmailProvider(logger, folderPath, timeout); err != nil {
		panic(err)
	} else {
		return provider
	}
}

func (p *FileEmailProvider) SendConfirmEmail(
	message appdto.Email,
) {
	ctx, cancel := p.createContext()
	defer cancel()
	p.saveFile(ctx, message.To, buildConfirmMessage(message))
}

func (p *FileEmailProvider) SendChangeEmail(
	message appdto.Email,
) {
	ctx, cancel := p.createContext()
	defer cancel()
	p.saveFile(ctx, message.To, buildChangeEmailMessage(message))
}

func (p *FileEmailProvider) SendResetPasswordEmail(
	message appdto.Email,
) {
	ctx, cancel := p.createContext()
	defer cancel()
	p.saveFile(ctx, message.To, buildResetPasswordMessage(message))
}

func (p *FileEmailProvider) createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), p.timeout)
}

func (p *FileEmailProvider) saveFile(ctx context.Context, to string, message []byte) {
	filename := fmt.Sprintf("%s/%s_%s", p.folderPath, to, time.Now().UTC().Format(time.RFC3339))
	p.logger.InfoContext(ctx, "try to create file", "filename", filename)
	file, err := os.Create(filename)
	if err != nil {
		p.logger.ErrorContext(ctx, "file creation failed", "error", err)
	}
	p.logger.InfoContext(ctx, "try to write into file")
	_, err = file.Write(message)
	if err != nil {
		p.logger.ErrorContext(ctx, "file writing failed", "error", err)
	}
}

func buildConfirmMessage(
	message appdto.Email,
) []byte {
	return fmt.Appendf(
		make([]byte, 0),
		"To: %s\r\n"+
			"Subject: Confirm your email\r\n"+
			"\r\n"+
			"To confirm your email, please click on the following link:\r\n"+
			"%s\r\n",
		message.To, message.Token,
	)
}

func buildChangeEmailMessage(
	message appdto.Email,
) []byte {
	return fmt.Appendf(
		make([]byte, 0),
		"To: %s\r\n"+
			"Subject: Confirm your email\r\n"+
			"\r\n"+
			"To confirm your email, please click on the following link:\r\n"+
			"%s\r\n",
		message.To, message.Token,
	)
}

func buildResetPasswordMessage(message appdto.Email) []byte {
	return fmt.Appendf(
		make([]byte, 0),
		"To: %s\r\n"+
			"Subject: Reset your password\r\n"+
			"\r\n"+
			"To reset your password, please click on the following link:\r\n"+
			"%s\r\n",
		message.To, message.Token,
	)
}

package smtp

import (
	"github.com/chrj/smtpd"
	"github.com/nonedotone/golog"
)

type Mailer interface {
	SendMessage([]byte) error
}

type Handler struct {
	addr string
	ms   Mailer
}

func NewHandler(addr string, ms Mailer) *Handler {
	return &Handler{
		addr: addr,
		ms:   ms,
	}
}

func (h *Handler) MailServer() error {
	srv := &smtpd.Server{
		Handler: h.MailHandler,
	}
	golog.Infof("mail server addr %s\n", h.addr)
	return srv.ListenAndServe(h.addr)
}

func (h *Handler) MailHandler(peer smtpd.Peer, env smtpd.Envelope) error {
	if err := h.ms.SendMessage(env.Data); err != nil {
		golog.Errorf("send mail by %s, data %s, error %v\n", env.Sender, string(env.Data), err)
		return err
	}
	golog.Infof("%s send mail success\n", env.Sender)
	return nil
}

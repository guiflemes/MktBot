package handlers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	secret = "jdjadjadjakdjka"
)

type PageAcesssAuth struct {
	signPrefix     string
	headerSignName string
}

func NewPageAcesssAuth() *PageAcesssAuth {
	return &PageAcesssAuth{
		signPrefix:     "sha1=",
		headerSignName: "X-Hub-Signature",
	}
}

func (p *PageAcesssAuth) Auth(c *fiber.Ctx) error {
	sig := string(c.Request().Header.Peek(p.headerSignName))

	if !strings.HasPrefix(sig, p.signPrefix) {
		return errors.New("not x-sign header provided")
	}

	validSig, err := p.isValidSignature(sig, c.Body())

	if err != nil {
		return err
	}

	if !validSig {
		return errors.New("invalid x-sign header")
	}

	return nil
}

func (p *PageAcesssAuth) isValidSignature(sig string, body []byte) (bool, error) {
	sign, err := hex.DecodeString(sig[len(p.signPrefix):])

	if err != nil {
		return false, err
	}

	hash := hmac.New(sha1.New, []byte(secret))
	hash.Reset()
	hash.Write(body)

	return hmac.Equal(hash.Sum(nil), sign), nil
}

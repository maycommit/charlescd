package error

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Error interface {
	Marshal() ([]byte, error)
	AddMeta(key, value string) *customError
	AddOperation(operation string) *customError
	LogFields() logrus.Fields
	Error() string
}

type customError struct {
	ID         string            `json:"id"`
	Title      string            `json:"title"`
	Detail     string            `json:"detail"`
	Meta       map[string]string `json:"meta"`
	Operations []string          `json:"operations"`
}

func (c *customError) AddMeta(key, value string) *customError {
	c.Meta[key] = value
	return c
}

func (c *customError) AddOperation(operation string) *customError {
	c.Operations = append(c.Operations, operation)
	return c
}

func (c *customError) LogFields() logrus.Fields {
	return logrus.Fields{
		"id": c.ID,
		"title":  c.Title,
		"detail": c.Detail,
		"operations": c.Operations,
	}
}

func (c *customError) Error() string {
	return fmt.Sprintf("%s", c.Detail)
}

func (c customError) Marshal() ([]byte, error) {
	return json.Marshal(&c)
}

func New(title string, detail string) Error {
	return &customError{
		ID:     uuid.New().String(),
		Title:  title,
		Detail: detail,
		Meta: map[string]string{
			"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
		},
		Operations: []string{},
	}
}

package sonyflake

import (
	"fmt"
	"strconv"

	sf "github.com/sony/sonyflake"
)

type Sonyflake interface {
	GenerateID() (string, error)
}

type sonyflakeImpl struct {
	flake *sf.Sonyflake
}

func New() (Sonyflake, error) {
	flake, err := sf.New(sf.Settings{})
	if err != nil {
		return nil, fmt.Errorf("init sonyflake: %w", err)
	}
	return &sonyflakeImpl{flake: flake}, nil
}

func (s *sonyflakeImpl) GenerateID() (string, error) {
	id, err := s.flake.NextID()
	if err != nil {
		return "", fmt.Errorf("generate sonyflake id: %w", err)
	}
	return strconv.FormatUint(id, 10), nil
}

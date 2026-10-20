package sonyflake

import (
	"fmt"

	sf "github.com/sony/sonyflake"
)

type Sonyflake interface {
	GenerateID() (int64, error)
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

func (s *sonyflakeImpl) GenerateID() (int64, error) {
	id, err := s.flake.NextID()
	if err != nil {
		return 0, fmt.Errorf("generate sonyflake id: %w", err)
	}
	return int64(id), nil
}

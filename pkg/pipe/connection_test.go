package pipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnection(t *testing.T) {
	c := NewConnection(1, "test")
	assert.NotNil(t, c.Data)
}

func TestNewReversedConnection(t *testing.T) {
	c := NewReversedConnection()
	assert.NotNil(t, c.Data)
}

func TestReversedConnection_GetKeys(t *testing.T) {
	c := NewReversedConnection()
	c.Data[1] = "test1"
	c.Data[2] = "test2"
	c.Data[3] = "test3"

	assert.Contains(t, c.GetKeys(), 1)
	assert.Contains(t, c.GetKeys(), 2)
	assert.Contains(t, c.GetKeys(), 3)

	c2 := NewReversedConnection()
	assert.Equal(t, 0, len(c2.GetKeys()))
}

func TestReversedConnection_GetInt(t *testing.T) {
	c := NewReversedConnection()
	c.Data[1] = "5-task"
	c.Data[2] = "10-user"
	c.Data[3] = "15-project"

	assert.Equal(t, 5, c.GetForeignID(1))
	assert.Equal(t, 10, c.GetForeignID(2))
	assert.Equal(t, 15, c.GetForeignID(3))
}
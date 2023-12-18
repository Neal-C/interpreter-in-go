package object

import (
	"fmt"
)

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
)

type ObjectType string
type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (self *Null) Type() ObjectType { return NULL_OBJ }
func (self *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (self *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (self *ReturnValue) Inspect() string  { return self.Value.Inspect() }

type Error struct {
	Message string
}

func (self *Error) Type() ObjectType { return ERROR_OBJ }
func (self *Error) Inspect() string  { return "ERROR: " + self.Message }

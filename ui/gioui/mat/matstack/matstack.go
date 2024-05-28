// This file is generated from mgl32/matstack/matstack.go; DO NOT EDIT

package matstack

import (
	"errors"

	"ui/gioui/mat"
)

// A MatStack is an OpenGL-style matrix stack,
// usually used for things like scenegraphs. This allows you
// to easily maintain matrix state per call level.
type MatStack[T mat.Float] []mat.Mat4[T]

func NewMatStack[T mat.Float] () *MatStack[T] {
	return &MatStack[T]{mat.Ident4[T]()}
}

// Push copies the top element and pushes it on the stack.
func (ms *MatStack[T]) Push() {
	(*ms) = append(*ms, (*ms)[len(*ms)-1])
}

// Pop removes the first element of the matrix from the stack, if there is only
// one element left there is an error.
func (ms *MatStack[T]) Pop() error {
	if len(*ms) == 1 {
		return errors.New("Cannot pop from mat stack, at minimum stack length of 1")
	}
	(*ms) = (*ms)[:len(*ms)-1]

	return nil
}

// RightMul multiplies the current top of the matrix by the argument.
func (ms *MatStack[T]) RightMul(m mat.Mat4[T]) {
	(*ms)[len(*ms)-1] = (*ms)[len(*ms)-1].Mul4(m)
}

// LeftMul multiplies the current top of the matrix by the argument.
func (ms *MatStack[T]) LeftMul(m mat.Mat4[T]) {
	(*ms)[len(*ms)-1] = m.Mul4((*ms)[len(*ms)-1])
}

// Peek returns the top element.
func (ms *MatStack[T]) Peek() mat.Mat4[T] {
	return (*ms)[len(*ms)-1]
}

// Load rewrites the top element of the stack with m
func (ms *MatStack[T]) Load(m mat.Mat4[T]) {
	(*ms)[len(*ms)-1] = m
}

// LoadIdent is a shortcut for Load(mgl.Ident4())
func (ms *MatStack[T]) LoadIdent() {
	(*ms)[len(*ms)-1] = mat.Ident4[T]()
}

package link

import (
	"andromeda/toolchain/assembler/object"
	"andromeda/toolchain/instruction"
	"errors"
	"fmt"
	"os"
)

type LinkerContext struct {
	objects    []*object.CodeObject
	origin     uint16
	outputSize uint16
	output     []instruction.Instruction
}

func NewLinkerContext(origin uint16, outputSize uint16) *LinkerContext {
	return &LinkerContext{
		origin:     origin,
		outputSize: outputSize,
	}
}

func (l *LinkerContext) Link(objects []*object.CodeObject) ([]instruction.Instruction, error) {
	l.objects = objects
	l.output = make([]instruction.Instruction, l.outputSize)

	hasErrors := l.putAllErrors()
	if hasErrors {
		return l.output, errors.New("error occurred during assembly")
	} else if l.doAnyOverlap() {
		return l.output, errors.New("overlapping segments")
	}

	for _, o := range l.objects {
		outputOrigin := o.Origin - l.origin
		for i := outputOrigin; i < outputOrigin+uint16(len(o.Code)); i++ {
			if int(i) >= len(l.output) {
				return l.output, errors.New("code overflows output buffer") /*TODO: How to say better??*/
			}

			l.output[i] = o.Code[i-outputOrigin]
		}
	}

	return l.output, nil
}

func (l *LinkerContext) doAnyOverlap() bool {
	for _, v := range l.objects {
		for _, w := range l.objects {
			if v.DoesOverlap(w) {
				return true
			}
		}
	}

	return false
}

func (l *LinkerContext) putAllErrors() bool {
	hasErrors := false
	for _, v := range l.objects {
		for _, e := range v.Errors {
			_, err := os.Stderr.WriteString(fmt.Sprintf("%v", e))
			hasErrors = true
			if err != nil {
				panic(err)
			}
		}
	}

	return hasErrors
}

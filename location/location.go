package location

import "fmt"

type Location struct {
	Line     int
	Column   int
	Position int
}

func (l Location) String() string {
	return fmt.Sprintf("position %d (line %d, column %d)", l.Position, l.Line, l.Column)
}

func (l Location) ShortString() string {
	return fmt.Sprintf("%d/%d:%d", l.Position, l.Line, l.Column)
}

func NewLocation(line, column, position int) Location {
	return Location{Line: line, Column: column, Position: position}
}

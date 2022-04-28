package gm

// === ID-er ===

type IDerInterface interface {
	ID() int
	setID(id int)
}
type IDer struct {
	id int
}

func (inst *IDer) ID() int      { return inst.id }
func (inst *IDer) setID(id int) { inst.id = id }

// === Type-er ===
// TODO:

type TyperInterface interface {
	Type() string
	setType(string)
	TypeID() int
	setTypeID(int)
}
type Typer struct {
	typ    string
	typeID int
}

func (inst *Typer) Type() string     { return inst.typ }
func (inst *Typer) setType(d string) { inst.typ = d }
func (inst *Typer) TypeID() int      { return inst.typeID }
func (inst *Typer) setTypeID(d int)  { inst.typeID = d }

package astparser

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	//"strings"
)

type aststate int

const (
	Root aststate = iota
	Code_block
	Null_expr
	Nil_expr
	//
	Call_expr
	Call_hash
	Call_vector
	Call_function
	//
	Definition_expr
	Iter_expr_definition
	//
	Multi_identifier
	Identifier
	//
	Assignment_expr
	Multi_assign
	//
	Function
	Parameter
	//Control Flow
	While_expr
	For_expr
	Iter_expr
	For_each_expr
	For_index_expr
	Condition_expr
	If_expr
	Continue_expr
	//
	Break_expr
	Return_expr
	//
	Tuple_expr
	Vector_expr
	Hash_expr
	//
	Ternary_operator
	Binary_operator
	Unary_operator
	//
	Hash_pair
	Slice_vector
	//
	Number_literal
	String_literal
	Bool_literal
)

func from(s string) aststate {
	var out aststate
	switch s {
	case "null":
		out = Null_expr
	case "nil":
		out = Nil_expr
	case "number":
		out = Number_literal
	case "string":
		out = String_literal
	case "identifier":
		out = Identifier
	case "bool":
		out = Bool_literal
	case "vector":
		out = Vector_expr
	case "hash":
		out = Hash_expr
	case "pair":
		out = Hash_pair
	case "function":
		out = Function
	case "block":
		out = Code_block
	case "parameter":
		out = Parameter
	case "ternary_operator":
		out = Ternary_operator
	case "binary_operator":
		out = Binary_operator
	case "unary_operator ":
		out = Unary_operator
	case "call_expr ":
		out = Call_expr
	case "call_hash ":
		out = Call_hash
	case "call_vector ":
		out = Call_vector
	case "call_function ":
		out = Call_function
	case "slice":
		out = Slice_vector
	case "definition":
		out = Definition_expr
	case "assignment":
		out = Assignment_expr
	case "multiple_assignment":
		out = Multi_assign
	case "multiple_identifier":
		out = Multi_identifier
	case "tuple":
		out = Tuple_expr
	case "multi_assignment":
		out = Multi_assign
	case "while":
		out = While_expr
	case "for":
		out = For_expr
	case "iterator":
		out = Iter_expr
	case "iterator_definition":
		out = Iter_expr_definition
	case "foreach":
		out = For_each_expr
	case "forindex":
		out = For_index_expr
	case "condition":
		out = Condition_expr
	case "if":
		out = If_expr
	case "continue":
		out = Continue_expr
	case "break":
		out = Break_expr
	case "return":
		out = Return_expr
	case "call_expr":
		out = Call_expr
	case "call_function":
		out = Call_function
	case "call_vector":
		out = Call_vector
	case "call_hash":
		out = Call_hash
	case "unary_operator":
		out = Unary_operator
	default:
		log.Fatalf("Error: not a comprehenedable string: %s", s)
	}
	return out
}

type identEvent int

const (
	Def identEvent = iota
	Assign
	Call
)

type NaVar struct {
	ident   string
	defline int
	astline int
	event   identEvent
}

type MemAST struct {
	rawast      *string
	definitions map[string]*NaVar
	assigns     map[string]*NaVar
	calls       map[string]*NaVar
}

type statemachine struct {
	lines []string
	idx   int
	state aststate
}

func (self *statemachine) parseline(mast *MemAST) bool {
	words, linenum, perr := mast.parseLine(self.lines[self.idx])
	if perr != nil {
		return false
	}
	current := Null_expr
	last := Null_expr
	for i := 0; i < len(words); i++ {
		if i == 0 {
			current = from(words[0])
			last = self.state
			self.state = current
		}
		switch current {
		case Definition_expr:
			return true
		case Identifier:
			v := NaVar{
				ident:   words[1],
				defline: linenum,
				astline: self.idx,
			}
			switch last {
			case Definition_expr:
				mast.definitions[words[1]] = &v
			case Assignment_expr:
				mast.assigns[words[1]] = &v
			case Call_expr:
				mast.calls[words[1]] = &v
			default:
				return false
			}
		default:
			return false
		}
	}
	return false
}

func Run(mast *MemAST) {
	var statemach statemachine = statemachine{
		lines: strings.Split(*mast.rawast, "\n"),
		idx:   0,
		state: Null_expr,
	}
	for statemach.idx = 0; statemach.idx < len(statemach.lines); statemach.idx++ {
		statemach.parseline(mast)
	}
}

func get_ast(filepath string) *MemAST {
	out, err := exec.Command("./external/nasal", "-a", filepath).Output()
	if err != nil {
		println(err.Error())
	}
	s := string(out)
	a := MemAST{
		rawast:      &s,
		definitions: make(map[string]*NaVar),
		assigns:     make(map[string]*NaVar),
		calls:       make(map[string]*NaVar),
	}
	return &a
}

type NotAParseableASTLine struct{ line string }

func (self NotAParseableASTLine) Error() string {
	return (fmt.Sprintf("This line is not a parsable AST line:\n%s\n", self.line))
}

func (self *MemAST) parseLine(line string) ([]string, int, error) {
	var words = make([]string, 0)
	halves := strings.Split(line, "->")
	if len(halves) < 2 {
		return words, 0, NotAParseableASTLine{line}
	}
	sections := strings.Split(halves[1], ":")
	linenum, linenumerr := strconv.Atoi(sections[1])
	if linenumerr != nil {
		log.Printf("Error: AST Line doesn't have parseable line number:\n%s\n     :%s", line, linenumerr.Error())
		return words, 0, NotAParseableASTLine{line}
	}

	words = strings.Split(strings.TrimLeft(halves[0], " -+|"), " ")

	return words, linenum, nil
}

func (self *MemAST) DumpDefs() {
	for k, v := range self.definitions {
		fmt.Println(k, v.astline)
	}
}
func (self *MemAST) DumpAssigns() {
	for k, v := range self.assigns {
		fmt.Println(k, v.astline)
	}
}
func (self *MemAST) DumpCalls() {
	for k, v := range self.calls {
		fmt.Println(k, v.astline)
	}
}

func GetAST(filepath string) *MemAST {
	ast := get_ast(filepath)
	return ast
}

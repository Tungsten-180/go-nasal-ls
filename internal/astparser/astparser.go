package astparser

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
	//"strings"
)

type NaVar struct {
	line       int
	ident      string
	natype     string
	parameters []string
}

type AST struct {
	line        map[int]string
	definitions map[string]*NaVar
}

func from(s string) int {
	var out int
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
	case "multiple_identifier ":
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
	default:
		log.Fatalf("Error: not a comprehenedable string: %s", s)
	}
	return out
}

const (
	Root int = iota
	Code_block
	Null_expr
	Nil_expr
	//
	Call_expr = 100
	Call_hash
	Call_vector
	Call_function
	//
	Definition_expr = 200
	Iter_expr_definition
	//
	Multi_identifier
	Identifier
	//
	Assignment_expr = 300
	Multi_assign
	//
	Function = 400
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

type statemachine struct {
	current_line int
}

func (self *statemachine) next(string) {
	switch self.current_line {
	case Definition_expr:
	}
}

func get_ast(filepath string) *AST {
	out, err := exec.Command("./external/nasal", "-a", filepath).Output()
	if err != nil {
		println(err.Error())
	}
	a := AST{
		line:        make(map[int]string),
		definitions: make(map[string]*NaVar),
	}
	a.load(string(out))
	return &a
}

func (self *AST) load(rawAst string) {
	lines := strings.Split(rawAst, "\n")
	for i, line := range lines {
		halves := strings.Split(line, "->")
		if len(halves) < 2 {
			continue
		}
		sections := strings.Split(halves[1], ":")
		linenum, linenumerr := strconv.Atoi(sections[1])
		if linenumerr != nil {
			log.Printf("Error: AST Line %d doesn't have parseable line number:\n%s\n     :%s", i, line, linenumerr.Error())
		}
		self.line[linenum] += line
		if strings.Contains(line, "definition") {

		}
	}

}

func GetAST(filepath string) *AST {
	ast := get_ast(filepath)
	return ast
}

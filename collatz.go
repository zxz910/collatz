package main

import (
	"fmt"
	"math"
)

var debug bool = true
var kCoeff int = 3;
var kConstant int = 1;

// A number of the form:
// coeff * o + constant
type CollatzExpr struct {
	coeff float64;
	constant float64;
}

// a collatz node that can describe a list of possible odds of the form: [1, 2, 3, x] for all x > start. And their exprs to go with it.
type VariableCollatzNode struct {
	odds_prefix []float64;
	odd_start float64;
	left_coeff_base float64;
	right_coeff_base float64;
	left_constant float64;
	right_constant float64;
}

type CollatzNode struct {
	odd []float64;
	expr CollatzExpr;
	level int;
}

// CollatzNodes grouped by sum of their odds lists.
// Allows us to iterate over possibilities in order by cardinality.
type CollatzList struct {
	nodes_map map[float64][]CollatzNode;
	level float64;
	index float64;
	max_level float64;
}

func appendCollatzNode(n *CollatzList, node CollatzNode) {
	if (debug) {
		fmt.Println("{");
		printCollatz(node);
		fmt.Println("}");
	}
	cardinality := 0.0;
	for _, i := range(node.odd) {
		cardinality += i;
	}
	n.max_level = math.Max(cardinality, n.max_level);
	n.nodes_map[cardinality] = append(n.nodes_map[cardinality], node);
}

func appendCollatzList(n *CollatzList, nodes ...CollatzNode) {
	for _, node := range(nodes) {
		appendCollatzNode(n, node);
	}
}

func GetNext(n *CollatzList) (bool, CollatzNode) {
	_, ok := n.nodes_map[n.level]
	if (ok) {
		if len(n.nodes_map[n.level]) > int(n.index) {
			n.index++;
			return true, n.nodes_map[n.level][int(n.index) - 1];
		} else {
			n.index = 0;
			n.level++;
			return GetNext(n);
		}

	} else {
		if (n.level < n.max_level) {
			n.level++;
			n.index = 0;
			return GetNext(n);
		}
	}
	return false, CollatzNode{};
}

/* I need a tree of values.
	 1
      1     2
     1 2      1

For each of these potential values I want to calculate the o' that could possibly satisfy these equations.

For the first step, i have:
(3o + 1) / 2 = o    => o = -1;
(3o + 1) / 4 = o    => o = 1;     This is indeed a cycle, but maybe the only 1!

... */

func printCollatz(n CollatzNode) {
	if (debug) {
		fmt.Println(n.odd);
		printCollatzExpr(calculateO(n.odd));
		//fmt.Println(n.level);
		printCollatzExpr(n.expr);
		fmt.Println("__________");
	}
}

func printCollatzExpr(n CollatzExpr) {
	fmt.Printf("%dx + %d\n", int(n.coeff), int(n.constant));
} 

func add(left CollatzExpr, right CollatzExpr) CollatzExpr {
	return CollatzExpr { coeff: left.coeff + right.coeff, constant: left.coeff + right.coeff }; 
}

func addc(left CollatzExpr, constant float64) CollatzExpr {
	return CollatzExpr { coeff: left.coeff, constant: left.constant + constant };
}

func mult(factor float64, expr CollatzExpr) CollatzExpr {
	return CollatzExpr { coeff: factor * expr.coeff, constant: factor * expr.constant };
}

func collatz(expr CollatzExpr) CollatzExpr {
	// 3x + 1
	return addc(mult(float64(kCoeff), expr), float64(kConstant));
}

func substituteO(expr CollatzExpr, power float64) CollatzExpr {
	// Plug in 2^i o' + 1 for o, in the expression: coeff * o + constant;
	// coeff * (2^i o' + 1) + constant
	// (coeff * 2^i) * o' + (coeff + constant)
	o_prime := CollatzExpr { coeff: math.Pow(2, power), constant: 1 };
	return addc(mult(expr.coeff, o_prime), expr.constant);
}

func calculateO(n []float64) CollatzExpr {
	// iterate over the list backwards, each time modifying 
	expr := CollatzExpr { coeff: 1, constant: 0 };
	for i := 0; i < len(n); i++ {
		expr = substituteO(expr, n[i]);
	}
	return expr;
}

func evaluateCollatz (n int, expr CollatzExpr) float64 {
	x := expr.coeff * float64(n) + expr.constant;
	return x;
}

func lessThan(left CollatzExpr, right CollatzExpr) bool {
	a := left.coeff - right.coeff;
	b := right.constant - left.constant;

	is_less_than := false;
	good := "XXX";
	if (left.coeff > right.coeff && left.constant > right.constant) {
		good = "XXX";
	 	is_less_than = false;
	} else if (left.coeff < right.coeff && left.constant < right.constant) {
		good = ">>>";	
		is_less_than = true;
	} else if (left.coeff < right.coeff && left.constant >= right.constant) {
		good = "<<<";
		is_less_than = true;
	} else {
		good = "<<<";
		is_less_than = true;
	}
	if (left.coeff == right.coeff) {
		fmt.Println("This should never happen");
		return false;
	}
	if (debug && false) {		
		fmt.Printf("%s      %dx + %d = %dx + %d;  ", good, int(left.coeff), int(left.constant), int(right.coeff), int(right.constant));
		fmt.Printf("constant: %d, coeff: %d, mid: %f\n", int(b), int(a), b/a);
	}

	if (b/a > 0) {
		fmt.Println("YAYAYAYAYAYAYA");
	}
	return is_less_than;
}

func copySlice(n []float64) []float64 {
	var slice []float64;
	for _, i := range n {
		slice = append(slice, i);
	}
	return slice;
}

func reduceCollatz(node CollatzNode) []CollatzNode {
	var nodes []CollatzNode;
	next_expr_b := node.expr
	var i float64;
	i = 1;

	
	for math.Remainder(next_expr_b.coeff + next_expr_b.constant, 2) == 0  {
		if (math.Remainder(next_expr_b.coeff, 2) == 0) {
			next_expr_b.coeff = next_expr_b.coeff / 2;
			next_expr_b.constant = next_expr_b.constant / 2;
		} else {
			break;
		}
	}

	if (math.Remainder(next_expr_b.coeff + next_expr_b.constant, 2) != 0) {			
		if (lessThan(calculateO(node.odd), next_expr_b)) {
			next_collatz := CollatzNode{odd : (copySlice(node.odd)), expr: next_expr_b}
			nodes = append(nodes, next_collatz);
		}
		return nodes;
	}
				
	for i < 7 {
		next_expr := (next_expr_b);
		// Reduce 
		if (math.Remainder(next_expr.coeff + next_expr.constant, 2) == 0) {
			d1 := evaluateCollatz(1, calculateO(node.odd));
			d2 := evaluateCollatz(1, next_expr);
			if (math.Remainder(d2 / d1, 2) == 0 && d1 != 1) {
				fmt.Println(d1);
				fmt.Println(d2);
				fmt.Println("FOUND");
				break;
			}
			next_collatz := CollatzNode{odd : append(copySlice(node.odd), i), expr: substituteO(next_expr, i)}
			next_nodes := reduceCollatz(next_collatz);
			nodes = append(nodes, next_nodes...);
			i++;
			if len(next_nodes) == 0 {
				break;
			} else {
				continue;
			}
		}
	}
	return nodes;
}


func generateCollatz(n int) {
	// Push 1, (coeff, constant)
	// calculate 3 * (2 o' + 1) + 1 => 6o' + 4 => 3o' + 2. So store <1>, (3, 2)
	// Need to make sure it doesn't go below o.
	var start CollatzNode;
	
	start.expr.coeff = 1;
	start.expr.constant = 0;
	start.level = 1;
	//start.odd = append(start.odd, 1);
	
	var collatz_list CollatzList;
	collatz_list.nodes_map = make(map[float64][]CollatzNode);
	
	pointer := start;
	pointer.odd = copySlice(start.odd);
	
	for len(pointer.odd) < n {
		printCollatz(pointer);
		new_collatz_nodes := reduceCollatz(CollatzNode {odd: pointer.odd, expr: collatz(pointer.expr)})
		appendCollatzList(&collatz_list, new_collatz_nodes...);
		has_next, p := GetNext(&collatz_list);
		if (has_next) {
			pointer = p;
		} else {
			break;
		}
	}
}

func runCollatz(n int) {
	start := n;
	started := false;
	for (n != 1) {
		if (n == start && started) {
			break;
		}
		started = true;
		if (math.Remainder(float64(n), 2) == 0) {
			n = n / 2;
		} else {
			fmt.Println(n);
			n = kCoeff * n + kConstant;
		}
	}	
	fmt.Println();
}

func runTests() {
	// expected:
	
	// [1 1 1 1 2 1 1]
	// 256x + 223
	// 729x + 638
	collatz_1 := CollatzNode{odd: []float64{1, 1, 1, 1, 2, 1}, expr: collatz(CollatzExpr{coeff : 243, constant: 182}), level: 1};
	for _, i := range reduceCollatz(collatz_1) {
		printCollatz(i);
	}
	
	// [1 1]
	// 4x + 3
	// 9x + 8
	// _______
	// [1 2]
	// 8x + 3
	// 9x + 4
	collatz_2 := CollatzNode{odd: []float64{1}, expr: collatz(CollatzExpr{coeff : 3, constant: 2}), level: 1};
	for _, i := range reduceCollatz(collatz_2) {
		printCollatz(i);
	}
	
	
	// __________
	// [1 2 1 1 1]
	// 64x + 59
	// 81x + 76
	// __________
	// [1 2 1 2]
	// 64x + 27
	// 162x + 71
	// __________
	// [1 2 1 3]
	// 128x + 27
	// 324x + 71
	// __________
	// [1 2 1 4]
	// 256x + 27
	// 648x + 71
	// __________
	// [1 2 1 5]
	// 512x + 27
	// 1296x + 71
	// __________
	// [1 2 1 6]
	// 1024x + 27
	// 2592x + 71
	// ___________
	collatz_3 := CollatzNode{odd: []float64{1, 2, 1}, expr: collatz(CollatzExpr{coeff : 27, constant : 20}), level: 1};
	for _, i := range reduceCollatz(collatz_3) {
		printCollatz(i);
	}
	
	// __________
	// [1 2 1]
	// 16x + 11
	// 27x + 20
	// __________
	collatz_4 := CollatzNode{odd: []float64{1, 2}, expr: collatz(CollatzExpr{coeff : 9 , constant : 4}), level: 1};
	for _, i := range reduceCollatz(collatz_4) {
		printCollatz(i);
	}
	
	// __________
	// [1 2 1 3]
	// 128x + 27
	// 486x + 107
	// __________
	collatz_5 := CollatzNode{odd: []float64{1, 2, 1, 3}, expr: collatz(CollatzExpr{coeff : 324 , constant : 71}), level: 1};
	for _, i := range reduceCollatz(collatz_5) {
		printCollatz(i);
	}
	
}


func main() {
	//generateCollatz(9);
	//evaluateCollatz(1, calculateO([]float64 {1, 2, 1, 1, 1, 1, 1, 1}));
	//runCollatz(61);
	runTests();
}


#ifndef COLLATZ_EXPR_H
#define COLLATZ_EXPR_H

class CollatzExpr {
 public:
  CollatzExpr(int coeff, int constant_);
  string ToString();

 private:
  int coeff_;
  int constant_;
};

#endif /* COLLATZ_EXPR_H */

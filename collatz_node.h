#ifndef COLLATZ_NODE_H
#define COLLATZ_NODE_H
class CollatzNode {
 public:
  CollatzNode(std::vector<int> odd_expos, CollatzExpr expr, int level);
  string ToString();

 private:
  std::vector<int> odd_expos_;
  CollatzExpr expr_;
  int level_;
}

#endif /* COLLATZ_NODE_H */

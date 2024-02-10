 # Interpreter in go
 
## Followed Thorsten ball's book : https://interpreterbook.com/
 
###  Learned and done :
- Pratt parsing
- Abstract Syntax Trees (AST)
- 2 pointer technique (slow & fast pointer)
- Tree/Graph traversal patterns
- Graph theory
- Test Driven Development 

## features and support:
- functions, functions calls, higher-order functions and closures
- Hashmaps
- arrays
- and scalar data types

### Fully functional interpreter of the Monkey-lang

To try it out:

requirements : Go >= 1.21.3 or Docker

```shell
git clone git@github.com:Neal-C/interpreter-in-go.git
cd interpreter-in-go
go run . 
# build executable command : go build -o ./bin/interpreter
# run executable : ./bin/interpreter
```

or via Docker by running my image

```shell
git clone git@github.com:Neal-C/interpreter-in-go.git
cd interpreter-in-go
docker build -t nealc:interpreter .
# builds the image
docker run -it --name nealc-interpreter nealc:interpreter 
# runs the image
```

Try a few a lines:

- try some of the built-in functions from ./evaluator/builtins.go
- try higher-order functions
- try lists and index access
- try to define a function and call it !
- Check the test cases in ./**/*_test.go files to see what behaviors and features are supported



```shell
puts("Hello!")
# Hello!
# null
puts(1234)
# 1234
# null
let people = [{"name": "Alice", "age": 24},{"name": "Neal-C", "age": 999}, {"name": "Anna", "age": 22}];
people[0]["name"];
# Alice
len(people)
# 3
first(people)
# {"name": "Alice", "age": 24}
last(people)
# {"name": "Anna", "age": 22} 
if (true) { 42 } else { "never"};
# 42
if (false) { 42 } else { return "never" }
# "never"
let a = 20
let b = 22;
a == b
# false
a + b;
# 42
let sum = fn(x,y) { return x + y };
sum(a,b)
# 42
```



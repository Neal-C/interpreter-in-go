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
requirements : Go >= 1.13

```shell
git clone git@github.com:Neal-C/interpreter-in-go.git
cd interpreter-in-go
go run . 
# build executable command : go build -o ./bin/interpreter
# run executable : ./bin/interpreter
```

Try a few a lines:

- try some of the built-in functions from ./evaluator/builtins.go
- try higher-order functions
- try lists and index access 

```shell
>> puts("Hello!")
#Hello!
#null
>> puts(1234)
#1234
#null
>> let people = [{"name": "Alice", "age": 24}, {"name": "Anna", "age": 22}];
>> people[0]["name"];
# Alice
>> len(people)
#2
>> first(people)
#{"name": "Alice", "age": 24}
```



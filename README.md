<p align="center"><img src="./assets/vorn.svg" width="128" height="128"></p>

# vorn

vorn is a simple interpreted C-like scripting language that I'm making to learn more about language design and implementation.

## Current features

* Integers, Floats, Booleans and Strings
* Arithmetic operations
* Comparison operators
* Variables
* Comments
* Arrays and Hashes
* If statements
* For and While loops
* Functions
* Built-in functions (See [evaluator.go](evaluator/evaluator.go#L23) for a complete list)
* A REPL

## Planned features

* Show line numbers in all parser and evaluator error messages
* Error handling
* Ternaries
* Modules/Namespaces/Importing
* Standard library
* Code formatter/Linter

## Example

```vorn
func fib(n) {
    if n <= 1 {
        return n;
    }

    return fib(n - 1) + fib(n - 2);
}

print(fib(10));
```

This script calculates the 10th number in the Fibonacci sequence, which is `55`.

## Building

To build vorn, you need to have [Go](https://golang.org/) installed. Then, run the following command:

```sh
go build
```

## Running

To run vorn, you can either run the executable directly or use the REPL. To run the REPL, run the following command:

```sh
./vorn
```

To run a script, run the following command:

```sh
./vorn path/to/script.vorn
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

This project was started using [Writing An Interpreter In Go](https://interpreterbook.com/), a book by Thorsten Ball. I highly recommend it if you're interested in writing your own interpreter.

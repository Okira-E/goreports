# SafeGo

safego is a Go library found at [github:Okira-E/safego](https://github.com/Okira-E/safego) that provides a Rust-like Option enum for safer handling of optional values. The Option type allows you to represent the presence or absence of a value in a type-safe manner, reducing the chances of runtime errors related to nil or null values.

## Installation

It will always be faster to make a file called `option.go` and copy the source code.

## Usage

### Creating an Option

To create an `Option` with a value, use the `Some` function:

```go
opt := safego.Some(42)
```

### Checking if an Option has a value

There are multiple ways to get the value of an Option. You can use the Expect method to retrieve the value and panic with a custom error message if the Option is empty:

```go
val := opt.Expect("Option is empty")
```

If you are confident that the Option has a value and want to retrieve it without error handling, you can use the Unwrap method. However, if the Option is empty, it will panic with a default error message:

```go
val := opt.Unwrap()
```

To safely get the value with a default fallback value, you can use the UnwrapOr method. If the Option is empty, it will return the provided default value:

```go
val := opt.UnwrapOr(0)
```

If you want to specify a function to generate the default value, you can use the UnwrapOrElse method. If the Option is empty, it will invoke the provided function and return its result:

```go
val := opt.UnwrapOrElse(func() int {
    // Generate and return the default value
})
```

## License

This project is licensed under the Unlicensed License - see the [UNLICENSE](UNLICENSE) file for details

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or a pull request on the GitHub repository.

## Acknowledgements

safego was inspired by Rust's Option enum and aims to provide a similar level of safety and expressiveness in Go code.


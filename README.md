# Spiker
[![PkgGoDev](https://pkg.go.dev/badge/github.com/shockerli/spiker)](https://pkg.go.dev/github.com/shockerli/spiker)
[![Go Report Card](https://goreportcard.com/badge/github.com/shockerli/spiker)](https://goreportcard.com/report/github.com/shockerli/spiker)
[![codecov](https://codecov.io/gh/shockerli/spiker/branch/master/graph/badge.svg)](https://codecov.io/gh/shockerli/spiker)
![GitHub top language](https://img.shields.io/github/languages/top/shockerli/spiker)

> A Go package implementation of real-time computing
>
> lexer is inspired by [tdop](https://github.com/cristiandima/tdop)

English | [中文](README_ZH.md)

## Install
```sh
go get -u github.com/shockerli/spiker
```


## Usage

- Execute


```go
spiker.Execute(`1 + 2 * 3 / 4`)
```

- ExecuteWithScope

```go
var scopes = spiker.NewScopeTable("demo", 1, nil)
scopes.Set("a", 3)
scopes.Set("b", 4)

spiker.ExecuteWithScope(`a * 2 + b`, )
```

- Format

```go
spiker.Format(`a + b * 3`)
```

## Architecture
![architecture](architecture.png)


## Examples

### Keywords
```js
true/false/if/else/in/...
```

### Instruction separation
Spiker requires instructions to be terminated with a semicolon at the end of each statement

### Data type
- Number
> Integer, Float
```js
123;
-123;
12.34;
-12.34;
```

- String
```js
"abc"
```

- Boolean
```js
true;
false;
```

- Array
```js
[];
[1, 2, "a"];
[1, [], [2,], [3, 4,], 5];
```

- Dict
```js
v = [9:99, 8:8.8, "hello":12.02];
v = ["name":"jioby", "age":18, "log":[1:2, 3:4]];
v = [1, 9:9.9, 3];
```

### Arithmetic Operators
```js
1 + 2 - 3 * 4 / 5 % 6;
-0.19 + 3.2983;
3 ** 2;
```

### Bitwise Operators
```js
1 & 2;
1 | 2;
1 ^ 2;
1 >> 2;
1 << 2;
```

### Comparison Operators
```js
3 == 2;
3 != 2;
3 > 2;
3 >= 2;
3 < 2;
3 <= 2;
```

### Logical Operators
```js
!2;
1 && 2;
1 || 2;
```

### Assignment Operators
```js
v = 2;
v += 2;
v -= 2;
v *= 2;
v /= 2;
v %= 2;

a = b = c = 100;
```

### In Operator
```js
"john" in ["joy", "john"]; // true
100 in [100, 200, 300]; // true
9 in [9:"999", 8:"888"]; // true
9 in "123456789"; // true
```

### Build-in functions
- export
> return the expression value and interrupt script
```js
export(v);
export(v[i]);
```

- exist
> determines whether a variable or index exists
```js
exist(v);
exist(v[i]);
```

- len
> return the length of a value
```js
len("123");
len(v);
len(v[i]);
```

- del
> delete one or more variable or index
```js
del(a)
del(a, b)
del(a, b[1])
del(a[i], b[i])
```

### Custom function
- single

```js
sum = (a, b) -> a + b;

export(sum(1, 2)); # 3
```

```js
pow2 = x -> x ** 2;

export(pow2(5)); # 25
```

- block

```js
max = (a, b) -> {
    if (a > b) {
        return a;
    } else {
        return b;
    }
};

export(max(1, 2)); # 2
```

### Control structures

- if/else

```js
if (a > b && c != d) {
    a = b;
} else if (c > b) {
    a = c;
} else {
    export(c + d);
}
```

- while

```js
while (true) {
  a = 0;
  while (a < 10) {
    print(a, "\n");
    a += 1;
  }
}
```


### More
```js
a = 101;
b = 102;
c = "103";
d = [1, 2, a, b, c+b, "abc"];

if (a > b) {
    export(len(d));
} else if (a < b) {
    if (c) {
        d = -1000;
    } else {
        d = 0;
    }

    d += len(c);
    export(d);
}

export(d[2]);
```


## License
This project is licensed under the terms of the [MIT](LICENSE) license.

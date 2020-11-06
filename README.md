# Spiker
> A Golang package implementation of real-time computing
>
> Inspired by [cristiandima/tdop](https://github.com/cristiandima/tdop)


## Install
```sh
go get -u github.com/shockerli/spiker
```

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

### Control structures
```js
if (a > b && c != d) {
    a = b;
} else if (c > b) {
    a = c;
} else {
    export(c + d);
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

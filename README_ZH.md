# Spiker
[![PkgGoDev](https://pkg.go.dev/badge/github.com/shockerli/spiker)](https://pkg.go.dev/github.com/shockerli/spiker) [![Go Report Card](https://goreportcard.com/badge/github.com/shockerli/spiker)](https://goreportcard.com/report/github.com/shockerli/spiker)
> 实时规则计算引擎
>
> lexer is inspired by [tdop](https://github.com/cristiandima/tdop)

[英文](README.md)

## 安装
```sh
go get -u github.com/shockerli/spiker
```

## 示例

### 关键词
```js
true/false/if/else/in/...
```

### 指令分隔符
每个语句后用英文分号（;）结束指令

### 数据类型
- 数字
> 整型/浮点型
```js
123;
-123;
12.34;
-12.34;
```

- 字符串
```js
"abc"
```

- 布尔
```js
true;
false;
```

- 数组
```js
[];
[1, 2, "a"];
[1, [], [2,], [3, 4,], 5];
```

- 字典
```js
v = [9:99, 8:8.8, "hello":12.02];
v = ["name":"jioby", "age":18, "log":[1:2, 3:4]];
v = [1, 9:9.9, 3];
```

### 算术运算符
```js
1 + 2 - 3 * 4 / 5 % 6;
-0.19 + 3.2983;
3 ** 2;
```

### 位运算符
```js
1 & 2;
1 | 2;
1 ^ 2;
1 >> 2;
1 << 2;
```

### 比较运算符
```js
3 == 2;
3 != 2;
3 > 2;
3 >= 2;
3 < 2;
3 <= 2;
```

### 逻辑运算符
```js
!2;
1 && 2;
1 || 2;
```

### 赋值运算符
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

### 内置函数
- export
> 中断脚本，并返回表达式的值
```js
export(v);
export(v[i]);
```

- exist
> 判断变量或索引是否存在
```js
exist(v);
exist(v[i]);
```

- len
> 返回值的长度或数量
```js
len("123");
len(v);
len(v[i]);
```

- del
> 删除一个或多个的变量或索引
```js
del(a)
del(a, b)
del(a, b[1])
del(a[i], b[i])
```

### 自定义函数
- 单行函数

```js
sum = (a, b) -> a + b;

export(sum(1, 2)); # 3
```

```js
pow2 = x -> x ** 2;

export(pow2(5)); # 25
```

- 块函数

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

### 流程控制

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


### 更多
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

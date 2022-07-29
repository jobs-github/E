/* comment */

func double(x) {
    x * 2;
};

/* comment */

var n = 8;
var r = double(n);
printf("%v\n", r);
var s = sprintf("origin: %v, double: %v\n", n, r);
print(s);
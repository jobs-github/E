/* comment */

func double(x) {
    x * 2;
};

/* comment */

const n = 8;
const r = double(n);
printf("%v\n", r);
const s = sprintf("origin: %v, double: %v\n", n, r);
print(s);
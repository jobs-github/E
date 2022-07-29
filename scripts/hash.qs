var h = {
    "k1": "v1",
    "k2": "v2",
    "k3": "v3"
};
println(h.keys());
var r = h.map(func(k, v) { k + v });
println(r);
var h2 = {
    "k1": 1,
    "k2": 2,
    "k3": 3
};
var acc = h2.reduce(func(x, y) { x + y }, 0);
println(acc);
var fr = h2.filter(func(k, v) { v > 1 });
println(fr);

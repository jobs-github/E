var arr = [1,2,3];
var r = arr.map(func(i, item) { item * 2 });
println(r);
var acc = arr.reduce(func(x, y) { x + y }, 0);
println(acc);
var fr = arr.filter(func(x) { x > 1 });
println(fr);
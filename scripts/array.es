const arr = [1,2,3];
const r = map(arr, func(i, item) { item * 2 });
println(r);
const acc = reduce(arr, func(x, y) { x + y }, 0);
println(acc);
const fr = filter(arr, func(i, item) { item > 1 });
println(fr);
const arr = [1,2,3];
const r = map(arr, func(i, item) { item * 2 });
println(r);
const acc = reduce(arr, 0, func(x, y) { x + y });
println(acc);
const fr = arr.filter(func(x) { x > 1 });
println(fr);
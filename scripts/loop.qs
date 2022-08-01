var n = 10;
var cond = func(i){ i < 10 };
var next = func(i){ i + 1 };

var r0 = for(0, cond, next, func(i, state){
    println(i + n);
}, 0);
println(r0);

var r1 = for(0, cond, next, func(i, state){
    i + state;
}, 0);

println(r1);
var s = "{\"k3\":{\"k31\":true,\"k32\":[1,2,3]},\"k2\":\"2\",\"k1\":1}"; 
var obj = loads(s);
println(obj);
var ss = dumps(obj);
println(ss);

var obj1 = { "k1": null, "k2": 123 };
var objstr = dumps(obj1);
println(objstr);
var obj2 = loads(objstr);
println(obj2);
println(obj1 == obj2);
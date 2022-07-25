func test_array_iter() {
    var arr = [10,20,30];
	var result = [];
    for var iter = arr.iter(); iter != null; iter = iter.next() {
        var idx = iter.key();
        var val = iter.value();
		result.push(idx);
		result.push(val);
    };
	println(result);
};

func test_hash_iter() {
    var h = {"k1": "v1", "k2": "v2", "k3": "v3"};
	var result = [];
    for var iter = h.iter(); iter != null; iter = iter.next() {
        var key = iter.key();
        var val = iter.value();
		result.push(key);
		result.push(val);
    };
	println(result);
};

func main() {
    test_array_iter();
    test_hash_iter();
};

main();
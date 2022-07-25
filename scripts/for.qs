func test_for1() {
    var i = 0;
    var result = 0;
    for ;; {
        if (i > 4) {
            break;
        }
        result = result + i;
        i = i + 1;
    }
    println("test_for1: ", result);
};

func test_for2() {
    var result = 0;
    for var i = 0; i < 5; i = i + 1 {
        result = result + i;
    }
    println("test_for2: ", result);
};

func test_for3() {
    var result = 0;
	for var i = 0; i < 5;  {
		result = result + i;
		i = i + 1;
	}
	println("test_for3: ", result);
};

func test_for4() {
    var i = 0;
	var result = 0;
	for ; i < 5;  {
		result = result + i;
		i = i + 1;
	}
	println("test_for4: ", result);
};

func test_for5() {
    var i = 0;
	var result = 0;
	for i = 0; i < 5; i = i + 1 {
		result = result + i;
	}
	println("test_for5: ", result);
};

func test_for6() {
    var i = 0;
	var result = 0;
	for {
		if (i > 4) {
			break;
		}
		result = result + i;
		i = i + 1;
	}
	println("test_for6: ", result);
};

func main() {
    var fns = [test_for1, test_for2, test_for3, test_for4, test_for5, test_for6];
    var sz = fns.len();
    for var i = 0; i < sz; i = i + 1 {
        fns[i]();
    }
};

main();
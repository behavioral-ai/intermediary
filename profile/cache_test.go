package profile

import (
	"encoding/json"
	"fmt"
)

func ExampleCache() {
	p := Cache{}

	buf, err := json.Marshal(p)

	fmt.Printf("test: Cache() -> [%v] [err:%v]\n", string(buf), err)

	//Output:
	//fail

}

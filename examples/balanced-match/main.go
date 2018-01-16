package main

import (
	"fmt"
	"regexp"

	"github.com/sniperkit/xfilter/backend/goblma"
)

func main() {
	v := goblma.Balanced("{", "}", "pre{in{nested}}post")
	fmt.Println(v) //=> { Start:3, End:14, Pre:"pre", Body:"in{nested}" Post:"post"}

	v = goblma.Balanced("{", "}", "pre{first}between{second}post")
	fmt.Println(v) //=> { Start:3, End:9, Pre:"pre", Body:"first" Post:"between{second}post"}

	//Regex example
	regStart := regexp.MustCompile(`\s+\{\s+`)
	regEnd := regexp.MustCompile(`\s+\}\s+`)
	v = goblma.Balanced(regStart, regEnd, "pre  {   in{nest}   }  post")
	fmt.Println(v) //=> { Start:3, End:17, Pre:"pre", Body:"in{nested}" Post:"post"}

}

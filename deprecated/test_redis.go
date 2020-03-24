package deprecated


import (
	"fmt"
"github.com/gomodule/redigo/redis"
"log"
)

var (
	c redis.Conn
	err error
	reply interface{}
)

func init() {
	c, err = redis.Dial("tcp", "a7e2fc632297911ea973806f8086ebf8-23313424.ca-central-1.elb.amazonaws.com:6379")
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.Do("AUTH", "P4mayz4dAK")
	if err != nil {
		panic(err)
	}
}

func main() {
	defer c.Close()
	//c.Do("FLUSHALL")
	//argsExample()
	//boolExample()
	//intExample()
	//intsExample()
	//scanExample()
	//scanSliceExample()
	//stringExample()
}

func argsExample() {
	var p1, p2 struct {
		Title  string `redis:"title"`
		Author string `redis:"author"`
		Body   string `redis:"body"`
	}

	p1.Title = "Example"
	p1.Author = "Gary"
	p1.Body = "Hello"

	if _, err := c.Do("HMSET", redis.Args{}.Add("id1").AddFlat(&p1)...); err != nil {
		fmt.Println(err)
		return
	}

	m := map[string]string{
		"title":  "Example2",
		"author": "Steve",
		"body":   "Map",
	}

	if _, err := c.Do("HMSET", redis.Args{}.Add("id2").AddFlat(m)...); err != nil {
		fmt.Println(err)
		return
	}

	for _, id := range []string{"id1", "id2"} {

		v, err := redis.Values(c.Do("HGETALL", id))
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := redis.ScanStruct(v, &p2); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%+v\n", p2)
	}
}

func boolExample() {
	c.Do("SET", "foo", 1)
	exists, _ := redis.Bool(c.Do("EXISTS", "foo"))
	fmt.Printf("%#v\n", exists)
}

func intExample() {
	c.Do("SET", "k1", 1)
	n, _ := redis.Int(c.Do("GET", "k1"))
	fmt.Printf("%#v\n", n)
	n, _ = redis.Int(c.Do("INCR", "k1"))
	fmt.Printf("%#v\n", n)
}

func intsExample() {
	c.Do("SADD", "set_with_integers", 4, 5, 6)
	ints, _ := redis.Ints(c.Do("SMEMBERS", "set_with_integers"))
	fmt.Printf("%#v\n", ints)
}

func scanExample() {
	c.Send("HMSET", "album:1", "title", "Red", "rating", 5)
	c.Send("HMSET", "album:2", "title", "Earthbound", "rating", 1)
	c.Send("HMSET", "album:3", "title", "Beat")
	c.Send("LPUSH", "albums", "1")
	c.Send("LPUSH", "albums", "2")
	c.Send("LPUSH", "albums", "3")
	values, err := redis.Values(c.Do("SORT", "albums",
		"BY", "album:*->rating",
		"GET", "album:*->title",
		"GET", "album:*->rating"))
	if err != nil {
		fmt.Println(err)
		return
	}

	for len(values) > 0 {
		var title string
		rating := -1 // initialize to illegal value to detect nil.
		values, err = redis.Scan(values, &title, &rating)
		if err != nil {
			fmt.Println(err)
			return
		}
		if rating == -1 {
			fmt.Println(title, "not-rated")
		} else {
			fmt.Println(title, rating)
		}
	}
}

func scanSliceExample() {
	c.Send("HMSET", "album:1", "title", "Red", "rating", 5)
	c.Send("HMSET", "album:2", "title", "Earthbound", "rating", 1)
	c.Send("HMSET", "album:3", "title", "Beat", "rating", 4)
	c.Send("LPUSH", "albums", "1")
	c.Send("LPUSH", "albums", "2")
	c.Send("LPUSH", "albums", "3")
	values, err := redis.Values(c.Do("SORT", "albums",
		"BY", "album:*->rating",
		"GET", "album:*->title",
		"GET", "album:*->rating"))
	if err != nil {
		fmt.Println(err)
		return
	}

	var albums []struct {
		Title  string
		Rating int
	}
	if err := redis.ScanSlice(values, &albums); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%v\n", albums)
}

func stringExample() {
	c.Do("SET", "hello", "world")
	s, err := redis.String(c.Do("GET", "hello"))
	fmt.Printf("%#v %v\n", s, err)
}
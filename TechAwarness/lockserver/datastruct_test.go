package lockserver

import "testing"
import "fmt"

/*****************************************************************************/

type clt struct {
	n int
}

func (*clt) Reply(r *MessageReply) {}

/*****************************************************************************/

func TestLockArea(t *testing.T) {

	la := NewLockArea()

	var c [3]*clt
	for i := 0; i < 3; i++ {
		c[i] = &clt{n: i}
		la.AddClient(c[i])
	}

	var ret bool
	var r Replier

	ret = la.Add(c[0], "toto")
	if ret != true {
		t.Error("Add failed c0")
	}
	ret = la.Add(c[1], "toto")
	if ret == true {
		t.Error("Add succeeded c1")
	}
	ret = la.Add(c[2], "toto")
	if ret == true {
		t.Error("Add succeeded c2")
	}
	r, ret = la.Remove(c[1], "toto")
	if ret == true || r != nil {
		t.Error("Remove oddity c1")
	}
	r, ret = la.Remove(c[0], "toto")
	if ret == false || r != c[1] {
		t.Error("Remove failed c0")
	}
	r, ret = la.Remove(c[1], "toto")
	if ret == false || r != c[2] {
		t.Error("Remove failed c1")
	}
	ret = la.Add(c[0], "toto")
	if ret == true {
		t.Error("Add succeeded c2")
	}
	if la.locks["toto"].Len() != 2 {
		t.Error("Wrong locks length")
	}
	rr := la.RemoveClient(c[2])
	if len(rr) != 1 || rr[0] != c[0] {
		t.Error("RemoveClient c2 wrong")
	}
}

/*****************************************************************************/

func ExampleLockArea() {

	la := NewLockArea()
	c0 := &clt{n: 0}
	la.AddClient(c0)

	locks := []string{"toto", "tutu", "titi"}
	for _, v := range locks {
		la.Add(c0, v)
	}

	for i, cnum := range []int{10, 1, 3, 2, 6, 7, 5, 4, 9, 8} {
		c := &clt{n: cnum}
		la.AddClient(c)
		la.Add(c, locks[i%len(locks)])
	}

	for e := la.locks["toto"].Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value.(*clt))
	}

	// Output:
	// &{0}
	// &{10}
	// &{2}
	// &{5}
	// &{8}
}

/*****************************************************************************/

func BenchmarkLockArea(b *testing.B) {

	la := NewLockArea()
	c := &clt{n: 1}
	la.AddClient(c)

	for i := 0; i < b.N; i++ {
		la.Add(c, "toto")
		la.Remove(c, "toto")
	}
}

/*****************************************************************************/

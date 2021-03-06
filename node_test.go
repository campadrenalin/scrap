package scrap

import (
	"code.google.com/p/go.net/html"
	"strings"
	"testing"
)

func setupNode(t *testing.T) Node {
	trq := NewTestRQ()
	request := trq.CreateRequest("foo")

	raw_node, err := html.Parse(strings.NewReader(sample_html))
	if err != nil {
		t.Fatal(err)
	}

	return Node{raw_node, &request}
}

func TestNode_Find_NoSuchElement(t *testing.T) {
	n := setupNode(t)
	found := n.Find("blah")
	if len(found) != 0 {
		t.Fatalf("Slice should have 0 elements, has %d!", len(found))
	}
}

func TestNode_Find_OneElement(t *testing.T) {
	n := setupNode(t)
	tagname := "body"
	found := n.Find(tagname)
	if len(found) != 1 {
		t.Fatalf("Slice should have 1 element, has %d!", len(found))
	}
	found_tagname := found[0].Node.Data
	if found_tagname != tagname {
		t.Fatalf("Expected %s, got %s", tagname, found_tagname)
	}
	if found[0].req != n.req {
		t.Fatalf(
			"Found node has different req! %r vs %r",
			found[0].req,
			n.req,
		)
	}
}

func TestNode_Find_MultipleElements(t *testing.T) {
	n := setupNode(t)
	tagname := "p"
	found := n.Find(tagname)
	if len(found) != 3 {
		t.Fatalf("Slice should have 1 element, has %d!", len(found))
	}
	for _, f := range found {
		found_tagname := f.Node.Data
		if found_tagname != tagname {
			t.Fatalf("Expected %s, got %s", tagname, found_tagname)
		}
		if f.req != n.req {
			t.Fatalf("Found node has different req! %r vs %r", f.req, n.req)
		}
	}
}

func TestNode_Find_BadSelector(t *testing.T) {
	n := setupNode(t)
	rq := n.req.RequestQueue
	selector := "*&{"

	found := n.Find(selector)
	if len(found) != 0 {
		t.Fatalf("Slice should have 0 elements, has %d!", len(found))
	}

	expected_remarks := ""
	expected_debug := "foo: End of Selector\n"
	got_remarks := rq.(*TestRQ).Remarks.String()
	got_debug := rq.(*TestRQ).Debug.String()

	compare(t, expected_remarks, got_remarks)
	compare(t, expected_debug, got_debug)
}

func TestNode_Queue_NoHref(t *testing.T) {
	n := setupNode(t)
	rq := n.req.RequestQueue

	n.Find("body").Queue()
	queued := rq.(*TestRQ).Queued
	if len(queued) != 0 {
		t.Fatal("Queue should be empty")
	}
}

func TestNode_Queue_WithHref(t *testing.T) {
	n := setupNode(t)
	rq := n.req.RequestQueue

	n.Find("link").Queue()
	queued := rq.(*TestRQ).Queued
	if len(queued) != 1 {
		t.Fatal("Queue should have one item")
	}
	queued_item := queued[0]
	compare(t, "/blah/", queued_item.Url)
	compare(t, rq, queued_item.RequestQueue)
}

func TestNode_Queue_Multiple(t *testing.T) {
	n := setupNode(t)
	rq := n.req.RequestQueue

	n.Find("a").Queue()
	queued := rq.(*TestRQ).Queued
	if len(queued) != 3 {
		t.Fatal("Queue should have three items")
	}
	expected_urls := []string{"/first", "/second", "/third"}
	for i, url := range expected_urls {
		queued_item := queued[i]
		compare(t, url, queued_item.Url)
		compare(t, rq, queued_item.RequestQueue)
	}
}

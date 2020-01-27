package goshopify

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestBlogList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/blogs.json", globalApiPathPrefix),
		httpmock.NewStringResponder(
			200,
			`{"blogs": [{"id":1},{"id":2}]}`,
		),
	)

	blogs, err := client.Blog.List(nil)
	if err != nil {
		t.Errorf("Blog.List returned error: %v", err)
	}

	expected := []Blog{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(blogs, expected) {
		t.Errorf("Blog.List returned %+v, expected %+v", blogs, expected)
	}

}

func TestBlogCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/blogs/count.json", globalApiPathPrefix),
		httpmock.NewStringResponder(
			200,
			`{"count": 5}`,
		),
	)

	cnt, err := client.Blog.Count(nil)
	if err != nil {
		t.Errorf("Blog.Count returned error: %v", err)
	}

	expected := 5
	if cnt != expected {
		t.Errorf("Blog.Count returned %d, expected %d", cnt, expected)
	}

}

func TestBlogGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/blogs/1.json", globalApiPathPrefix),
		httpmock.NewStringResponder(
			200,
			`{"blog": {"id":1}}`,
		),
	)

	blog, err := client.Blog.Get(1, nil)
	if err != nil {
		t.Errorf("Blog.Get returned error: %v", err)
	}

	expected := &Blog{ID: 1}
	if !reflect.DeepEqual(blog, expected) {
		t.Errorf("Blog.Get returned %+v, expected %+v", blog, expected)
	}

}

func TestBlogCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/blogs.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("blog.json"),
		),
	)

	blog := Blog{
		Title: "Mah Blog",
	}

	returnedBlog, err := client.Blog.Create(blog)
	if err != nil {
		t.Errorf("Blog.Create returned error: %v", err)
	}

	expectedInt := int64(241253187)
	if returnedBlog.ID != expectedInt {
		t.Errorf("Blog.ID returned %+v, expected %+v", returnedBlog.ID, expectedInt)
	}

}

func TestBlogUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"PUT",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/blogs/1.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("blog.json"),
		),
	)

	blog := Blog{
		ID:    1,
		Title: "Mah Blog",
	}

	returnedBlog, err := client.Blog.Update(blog)
	if err != nil {
		t.Errorf("Blog.Update returned error: %v", err)
	}

	expectedInt := int64(241253187)
	if returnedBlog.ID != expectedInt {
		t.Errorf("Blog.ID returned %+v, expected %+v", returnedBlog.ID, expectedInt)
	}
}

func TestBlogDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://fooshop.myshopify.com/%s/blogs/1.json", globalApiPathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.Blog.Delete(1)
	if err != nil {
		t.Errorf("Blog.Delete returned error: %v", err)
	}
}

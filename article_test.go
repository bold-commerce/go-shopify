package goshopify

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestBlogArticleList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/blogs/1/articles.json", client.pathPrefix),
		httpmock.NewStringResponder(
			200,
			`{"articles": [{"id": 1051293780,"title": "Welcome to the world of tomorrow!","created_at": "2013-11-06T19:00:00-05:00","body_html": "Good news, everybody!","blog_id": 241253187,"author": "dennis","user_id": null,"published_at": null,"updated_at": "2022-04-05T13:17:47-04:00","summary_html": null,"template_suffix": null,"handle": "welcome-to-the-world-of-tomorrow","tags": "","admin_graphql_api_id": "gid://shopify/OnlineStoreArticle/1051293780"}]}`,
		),
	)

	options := ArticleListOptions{
		ListOptions: ListOptions{Limit: 1},
	}

	articles, err := client.Article.ListBlog("1", options)
	if err != nil {
		panic(fmt.Sprintf("Cannot get blog list err: %s", err))
	}

	for _, article := range articles {
		fmt.Printf("ID: %v\n", article.ID)
		fmt.Printf("Title: %v\n", article.Title)
		fmt.Printf("CreatedAt: %v\n", article.CreatedAt)
		fmt.Printf("BodyHtml: %v\n", article.BodyHtml)
		fmt.Printf("ArticleId: %v\n", article.ArticleId)
		fmt.Printf("Author: %v\n", article.Author)
		fmt.Printf("UserId: %v\n", article.UserId)
		fmt.Printf("PublishedAt: %v\n", article.PublishedAt)
		fmt.Printf("UpdatedAt: %v\n", article.UpdatedAt)
		fmt.Printf("SummaryHtml: %v\n", article.SummaryHtml)
		fmt.Printf("TemplateSuffix: %v\n", article.TemplateSuffix)
		fmt.Printf("Handle: %v\n", article.Handle)
		fmt.Printf("Tags: %v\n", article.Tags)
		fmt.Printf("AdminGraphqlApiId: %v\n", article.AdminGraphqlApiId)
		fmt.Println("========")

	}
}

func TestBlogArticleCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/blogs/1/articles/count.json", client.pathPrefix),
		httpmock.NewStringResponder(
			200,
			`{"count": 4}`,
		),
	)

	total, err := client.Article.Count("1")
	if err != nil {
		panic(fmt.Sprintf("Cannot get blog count err: %s", err))
	}

	fmt.Println(total)
}

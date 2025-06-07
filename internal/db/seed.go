package db

import (
	"context"
	"fmt"
	"github.com/caturandi-labs/go-social/internal/store"
	"log"
	"math/rand"
)

var usernames = []string{
	"andy", "bob", "james", "alice", "charlie",
	"david", "eve", "frank", "grace", "heidi",
	"ivan", "judy", "kevin", "lisa", "mike",
	"natalie", "oscar", "peter", "quincy", "rachel",
	"steve", "tina", "ursula", "victor", "wendy",
	"xavier", "yvonne", "zach", "anna", "brian",
	"claire", "daniel", "elizabeth", "felix", "gabriel",
	"hannah", "ian", "julia", "karen", "leo",
	"monica", "nathan", "olivia", "patrick", "quinn",
	"rebecca", "samuel", "tracy", "umar", "veronica",
}

var titles = []string{
	"Mastering Go: A Beginner's Guide",
	"Concurrency in Go: Unlocking the Power of Goroutines",
	"Building RESTful APIs with Go and Gin",
	"Go vs. Python: Which Language Should You Choose?",
	"Understanding Go's Memory Management and Garbage Collection",
	"How to Write Efficient Code in Go: Best Practices",
	"Building Microservices with Go and gRPC",
	"The Ultimate Guide to Go Modules for Dependency Management",
	"Testing in Go: Writing Reliable Unit Tests",
	"Deploying Go Applications to Kubernetes",
	"Error Handling in Go: Why It’s Different and How to Do It Right",
	"Creating Web Scrapers with Go and Colly",
	"Go for DevOps: Automating Tasks with Go Scripts",
	"Exploring Go's Standard Library: Hidden Gems You Should Know",
	"Building a CLI Tool in Go: From Zero to Hero",
	"Performance Optimization in Go: Tips and Tricks",
	"Introduction to Go Templates for Dynamic HTML Rendering",
	"Working with Databases in Go: SQL and Beyond",
	"Real-Time Applications with Go and WebSockets",
	"Go in 2024: Predictions and Trends for the Future",
}

var tags = []string{
	"golang", "programming", "webdev", "backend", "frontend",
	"database", "api", "microservices", "docker", "kubernetes",
	"cloud", "devops", "linux", "security", "performance",
	"testing", "ci-cd", "version-control", "networking", "debugging",
}

var contents = []string{
	"Discover the power of Go's concurrency model and how goroutines can simplify your code...",
	"Learn the basics of Go's type system and why it's both simple and powerful...",
	"Explore best practices for writing clean, maintainable Go code in large projects...",
	"Understand the importance of interfaces in Go and how they enable flexible designs...",
	"Dive into Go's built-in testing framework and write robust tests for your applications...",
	"Master error handling in Go and why explicit errors are a feature, not a bug...",
	"Uncover the secrets of efficient memory management using Go's garbage collector...",
	"Get started with Go modules and manage dependencies like a pro...",
	"Build RESTful APIs in Go using popular frameworks like Gin or Echo...",
	"Learn how to deploy Go applications to the cloud with Docker and Kubernetes...",
	"Optimize your Go programs for performance using profiling tools like pprof...",
	"Understand the role of context in Go and how it helps manage request-scoped data...",
	"Explore the world of microservices architecture with Go as the backbone...",
	"Write idiomatic Go code by following the principles outlined in Effective Go...",
	"Discover how Go's simplicity makes it an ideal choice for DevOps tooling...",
	"Build command-line tools in Go and automate repetitive tasks efficiently...",
	"Learn about Go's HTTP package and how to create high-performance web servers...",
	"Understand the trade-offs between embedding and composition in Go's struct design...",
	"Implement secure authentication and authorization in Go-based web applications...",
	"Take a deep dive into Go's sync package and explore advanced concurrency patterns...",
}

var postComments = []string{
	"Great post! Very informative and well-written.",
	"I learned a lot from this article. Thanks for sharing!",
	"This was a fantastic read. Keep up the good work!",
	"I have a question about the second point you made. Can you clarify?",
	"Thanks for breaking this down into simple terms. It really helped me understand.",
	"Interesting perspective! I hadn't thought about it that way before.",
	"I disagree with your conclusion, but I appreciate the discussion.",
	"Could you provide more examples for the third section? It felt a bit brief.",
	"Loved the practical tips. I'll definitely try them out!",
	"Your writing style is engaging and easy to follow. Kudos!",
	"This article was exactly what I needed. Thank you!",
	"I think there’s a typo in the first paragraph. Did you mean 'example' instead of 'sample'?",
	"The visuals in this post were awesome. They really complemented the text.",
	"I found this post through a friend, and I’m glad I did. Subscribed!",
	"Can you recommend any further reading on this topic? I’m intrigued.",
	"The step-by-step guide was incredibly helpful. Thanks so much!",
	"I wish you’d expanded more on the last section. It felt rushed.",
	"This is one of the best articles I’ve read on this subject. Well done!",
	"I tried implementing your suggestions, and they worked like a charm!",
	"Looking forward to your next post. Keep them coming!",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(50)

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user", user, err)
			return
		}
	}

	posts := generatePosts(20, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post", post, err)
			return
		}
	}

	comments := generateComments(20, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment", comment, err)
			return
		}
	}

	log.Println("Seeding Complete")

}

func generateUsers(count int) []*store.User {
	users := make([]*store.User, count)
	for i := 0; i < count; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(users)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(users)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "password",
		}
	}
	return users
}

func generatePosts(count int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, count)

	for i := 0; i < count; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Version: 0,
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}

	}
	return posts
}

func generateComments(count int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, count)

	for i := 0; i < count; i++ {
		comments[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: postComments[rand.Intn(len(postComments))],
		}

	}
	return comments
}

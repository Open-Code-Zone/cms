package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"githum.com/Open-Code-Zone/cms/types"
	"githum.com/Open-Code-Zone/cms/utils"
	"githum.com/Open-Code-Zone/cms/views/components"
	"githum.com/Open-Code-Zone/cms/views/pages"
)

const (
	contentPath = "content/posts"
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}
}

func sanitizeFilename(title string) string {
	title = strings.ToLower(title)
	title = strings.ReplaceAll(title, " ", "-")

	// Remove any character that isn't alphanumeric or hyphen
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' {
			return r
		}
		return -1
	}, title)
}

func (h *Handler) HandleAddBlogPost(w http.ResponseWriter, r *http.Request) {
	blogPost := &types.BlogPost{
		Metadata: types.BlogMetadata{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Date:        r.FormValue("date"),
			Authors:     []string{r.FormValue("authors")},
			Image:       r.FormValue("image"),
			Tags:        []string{r.FormValue("tags")},
		},
		Content: r.FormValue("content"),
	}
	fmt.Println(blogPost.Content)

	gc, err := utils.NewGitHubClient()
	if err != nil {
		log.Println("Error creating GitHub client:", err)
		http.Error(w, "Failed to create GitHub client", http.StatusInternalServerError)
	}

	fileName := sanitizeFilename(blogPost.Metadata.Title) + ".md"
	filePath := filepath.Join(contentPath, fileName)

	fileContent, err := utils.GenerateMarkdown(blogPost.Metadata, blogPost.Content)
	if err != nil {
		log.Println("Error generating markdown:", err)
		http.Error(w, "Failed to generate markdown", http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Add new blog post: %s", blogPost.Metadata.Title)

	err = gc.CreateFile(filePath, fileContent, message)
	if err != nil {
		log.Println("Error creating file:", err)
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Blog post created successfully"))
}

func (h *Handler) HandleUpdateBlogPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["id"]

	blogPost := &types.BlogPost{
		Metadata: types.BlogMetadata{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Date:        r.FormValue("date"),
			Authors:     []string{r.FormValue("authors")},
			Image:       r.FormValue("image"),
			Tags:        []string{r.FormValue("tags")},
		},
		Content: r.FormValue("content"),
	}

	gc, err := utils.NewGitHubClient()
	if err != nil {
		log.Println("Error creating GitHub client:", err)
		http.Error(w, "Failed to create GitHub client", http.StatusInternalServerError)
	}

	fileContent, err := utils.GenerateMarkdown(blogPost.Metadata, blogPost.Content)
	if err != nil {
		log.Println("Error generating markdown:", err)
		http.Error(w, "Failed to generate markdown", http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("updated the blog post: %s", blogPost.Metadata.Title)
	filePath := filepath.Join(contentPath, fileName)

	err = gc.UpdateFile(filePath, fileContent, message)
	if err != nil {
		log.Println("Error updating file:", err)
		http.Error(w, "Failed to update file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	components.Toaster("Blog Post updated succesfully", "success").Render(r.Context(), w)
}

func (h *Handler) HandleDeleteBlogPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["id"]

	gc, err := utils.NewGitHubClient()
	if err != nil {
		log.Println("Error creating GitHub client:", err)
		http.Error(w, "Failed to create GitHub client", http.StatusInternalServerError)
	}

	message := fmt.Sprintf("deleted the blog post: %s", fileName)
	filePath := filepath.Join(contentPath, fileName)

	err = gc.DeleteFile(filePath, message)
	if err != nil {
		log.Println("Error deleting file:", err)
		w.WriteHeader(http.StatusInternalServerError)
		components.Toaster("Couldn't able to delete through GitHub API", "danger")
		return
	}

	components.Toaster("Blog Post deleted succesfully", "success").Render(r.Context(), w)
}

func (h *Handler) HandleBlogPostEditPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	gc, err := utils.NewGitHubClient()
	if err != nil {
		log.Printf("Error creating GitHub client: %v", err)
		http.Error(w, "Failed to initialize GitHub client", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(contentPath, id)
	content, err := gc.GetFileContent(filePath)

	blogPost, err := utils.ParseMarkdown(content)

	if err != nil {
		log.Printf("Error parsing markdown: %v", err)
		http.Error(w, "Failed to parse markdown", http.StatusInternalServerError)
		return
	}

	pages.EditBlogPosts(id, blogPost).Render(r.Context(), w)
}

func (h *Handler) HandleShowAllBlogPostsPage(w http.ResponseWriter, r *http.Request) {
	gc, err := utils.NewGitHubClient()
	if err != nil {
		log.Println("Error creating GitHub client:", err)
		http.Error(w, "Failed to create GitHub client", http.StatusInternalServerError)
	}

	contents, err := gc.GetContents(contentPath)
	if err != nil {
		log.Printf("Error fetching repository contents: %v", err)
		http.Error(w, "Failed to fetch repository contents", http.StatusInternalServerError)
		return
	}
	var markdownFiles []*types.MarkdownFile

	for _, content := range contents {
		if content.GetType() == "file" && strings.HasSuffix(content.GetName(), ".md") {
			//commits, err := gc.ListCommits(content.GetPath())
			//if err != nil || len(commits) == 0 {
			//	log.Printf("Error fetching commit history for file %s: %v", content.GetName(), err)
			//	continue
			//}
			markdownFiles = append(markdownFiles, &types.MarkdownFile{
				FileName: content.GetName(),
			})
		}
	}

	pages.ShowBlogPosts(markdownFiles).Render(r.Context(), w)
}

func (h *Handler) HandleNewBlogPostPage(w http.ResponseWriter, r *http.Request) {
	pages.EditBlogPosts("new-draft", nil).Render(r.Context(), w)
}

func (h *Handler) HandleNewBlogPost(w http.ResponseWriter, r *http.Request) {
	blogPost := &types.BlogPost{
		Metadata: types.BlogMetadata{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Date:        r.FormValue("date"),
			Authors:     []string{r.FormValue("authors")},
			Image:       r.FormValue("image"),
			Tags:        []string{r.FormValue("tags")},
		},
		Content: r.FormValue("content"),
	}

	fmt.Println(blogPost)

	gc, err := utils.NewGitHubClient()
	if err != nil {
		log.Println("Error creating GitHub client:", err)
		http.Error(w, "Failed to create GitHub client", http.StatusInternalServerError)
	}

	fileContent, err := utils.GenerateMarkdown(blogPost.Metadata, blogPost.Content)
	if err != nil {
		log.Println("Error generating markdown:", err)
		http.Error(w, "Failed to generate markdown", http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("updated the blog post: %s", blogPost.Metadata.Title)

	fileName := sanitizeFilename(blogPost.Metadata.Title) + ".md"
	filePath := filepath.Join(contentPath, fileName)

	err = gc.CreateFile(filePath, fileContent, message)
	if err != nil {
		log.Println("Error updating file:", err)
		http.Error(w, "Failed to update file", http.StatusInternalServerError)
		return
	}

	components.Toaster("Blog Post created and published succesfully", "success").Render(r.Context(), w)
	// redirect to show all blog posts page
	//http.Redirect(w, r, "/blog-post", http.StatusSeeOther)
}

package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Open-Code-Zone/cms/config"
	"github.com/Open-Code-Zone/cms/types"
	"github.com/Open-Code-Zone/cms/views/components"
	"github.com/Open-Code-Zone/cms/views/pages"
	"github.com/gorilla/mux"
)

const (
	contentPath = "content/posts"
)

func (h *Handler) HandleBlogPostEditPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	blogConfig := config.Envs.BlogPostConfig

	filePath := filepath.Join(contentPath, id)
	blogPost, err := h.githubClient.GetFileContent(filePath)

	// go routine to have asynchromenous changes to db
	if err != nil {
		log.Printf("Error parsing markdown: %v", err)
		http.Error(w, "Failed to parse markdown", http.StatusInternalServerError)
		return
	}

	pages.EditBlogPosts(id, &blogConfig, &blogPost).Render(r.Context(), w)
}

func (h *Handler) HandleShowAllBlogPostsPage(w http.ResponseWriter, r *http.Request) {
	contents, err := h.githubClient.GetContents(contentPath)
	if err != nil {
		log.Printf("Error fetching repository contents: %v", err)
		http.Error(w, "Failed to fetch repository contents", http.StatusInternalServerError)
		return
	}
	var markdownFiles []*types.MarkdownFile

	for _, content := range contents {
		if content.GetType() == "file" && strings.HasSuffix(content.GetName(), ".md") {
			markdownFiles = append(markdownFiles, &types.MarkdownFile{
				FileName: content.GetName(),
			})
		}
	}

	pages.ShowBlogPosts(markdownFiles).Render(r.Context(), w)
}

func (h *Handler) HandleNewBlogPostPage(w http.ResponseWriter, r *http.Request) {
	pages.EditBlogPosts("new-draft.md", &config.Envs.BlogPostConfig, nil).Render(r.Context(), w)
}

func (h *Handler) HandleNewBlogPost(w http.ResponseWriter, r *http.Request) {
	fileName := r.FormValue("fileName")
	fileContent := r.FormValue("content")

	fmt.Println(fileContent)

	message := fmt.Sprintf("🤖 - Added the new blog post: %s", fileName)

	filePath := filepath.Join(contentPath, fileName)

	err := h.githubClient.CreateFile(filePath, fileContent, message)
	if err != nil {
		log.Println("Error updating file:", err)
		components.Toaster("Looks like same blog post with the same file name exists", "danger").Render(r.Context(), w)
		return
	}

	components.Toaster("Blog Post created and published succesfully", "success").Render(r.Context(), w)
	// redirect to show all blog posts page
	//http.Redirect(w, r, "/blog-post", http.StatusSeeOther)
}

func (h *Handler) HandleDeleteBlogPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["id"]

	message := fmt.Sprintf("🤖 - Deleted the blog post: %s", fileName)
	filePath := filepath.Join(contentPath, fileName)

	err := h.githubClient.DeleteFile(filePath, message)
	if err != nil {
		log.Println("Error deleting file:", err)
		w.WriteHeader(http.StatusInternalServerError)
		components.Toaster("Couldn't able to delete through GitHub API", "danger")
		return
	}

	components.Toaster("Blog Post deleted succesfully", "success").Render(r.Context(), w)
}

func (h *Handler) HandleUpdateBlogPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["id"]
	fileContent := r.FormValue("content")

	message := fmt.Sprintf("🤖 - Updated the blog post: %s")
	filePath := filepath.Join(contentPath, fileName)

	err := h.githubClient.UpdateFile(filePath, fileContent, message)
	if err != nil {
		log.Println("Error updating file:", err)
		http.Error(w, "Failed to update file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	components.Toaster("Blog Post updated succesfully", "success").Render(r.Context(), w)
}

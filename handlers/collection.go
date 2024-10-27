package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Open-Code-Zone/cms/types"
	"github.com/Open-Code-Zone/cms/views/components"
	"github.com/Open-Code-Zone/cms/views/pages"
	"github.com/gorilla/mux"
)

// listing all the collection items for example blog posts, authors
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionConfig := h.store.Collections.GetCollectionConfig(vars["collection"])
	if collectionConfig == nil {
		http.Error(w, "Collection doesn't exists", http.StatusInternalServerError)
		return
	}
	// error handeling that collection is not found

	contents, err := h.githubClient.GetContents(collectionConfig.GitPath)
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

	pages.ShowCollectionItems(markdownFiles, collectionConfig).Render(r.Context(), w)
}

// rendering markdown editor with metadata form for creating new collection item
func (h *Handler) New(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionConfig := h.store.Collections.GetCollectionConfig(vars["collection"])
	if collectionConfig == nil {
		http.Error(w, "Collection doesn't exists", http.StatusInternalServerError)
		return
	}

	pages.EditCollection("new-draft.md", nil, collectionConfig).Render(r.Context(), w)
}

// creating collection item
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionConfig := h.store.Collections.GetCollectionConfig(vars["collection"])
	if collectionConfig == nil {
		http.Error(w, "Collection doesn't exists", http.StatusInternalServerError)
		return
	}

	fileName := r.FormValue("fileName")
	fileContent := r.FormValue("content")

	fmt.Println(fileContent)

	message := fmt.Sprintf("🤖 - Added the new blog post: %s", fileName)

	filePath := filepath.Join(collectionConfig.GitPath, fileName)

	err := h.githubClient.CreateFile(filePath, fileContent, message)
	if err != nil {
		log.Println("Error updating file:", err)
		components.Toaster("Looks like same blog post with the same file name exists", "danger").Render(r.Context(), w)
		return
	}

	components.Toaster("Blog Post created and published succesfully", "success").Render(r.Context(), w)
}

// rendering markdown editor with metadata form for editing collection item
func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	collectionConfig := h.store.Collections.GetCollectionConfig(vars["collection"])
	if collectionConfig == nil {
		http.Error(w, "Collection doesn't exists", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(collectionConfig.GitPath, id)
	fileContent, err := h.githubClient.GetFileContent(filePath)

	// go routine to have asynchromenous changes to db
	if err != nil {
		log.Printf("Error parsing markdown: %v", err)
		http.Error(w, "Failed to parse markdown", http.StatusInternalServerError)
		return
	}

	pages.EditCollection(id, &fileContent, collectionConfig).Render(r.Context(), w)
}

// updating collection item with new content
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionConfig := h.store.Collections.GetCollectionConfig(vars["collection"])
	if collectionConfig == nil {
		http.Error(w, "Collection doesn't exists", http.StatusInternalServerError)
		return
	}

	fileName := vars["id"]
	fileContent := r.FormValue("content")

	message := fmt.Sprintf("🤖 - Updated the blog post: %s")
	filePath := filepath.Join(collectionConfig.GitPath, fileName)

	err := h.githubClient.UpdateFile(filePath, fileContent, message)
	if err != nil {
		log.Println("Error updating file:", err)
		http.Error(w, "Failed to update file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	components.Toaster("Blog Post updated succesfully", "success").Render(r.Context(), w)
}

// deleting collection item
func (h *Handler) Destroy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionConfig := h.store.Collections.GetCollectionConfig(vars["collection"])
	if collectionConfig == nil {
		http.Error(w, "Collection doesn't exists", http.StatusInternalServerError)
		return
	}

	fileName := vars["id"]

	message := fmt.Sprintf("🤖 - Deleted the blog post: %s", fileName)
	filePath := filepath.Join(collectionConfig.GitPath, fileName)

	err := h.githubClient.DeleteFile(filePath, message)
	if err != nil {
		log.Println("Error deleting file:", err)
		w.WriteHeader(http.StatusInternalServerError)
		components.Toaster("Couldn't able to delete through GitHub API", "danger")
		return
	}

	components.Toaster("Blog Post deleted succesfully", "success").Render(r.Context(), w)
}

package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Open-Code-Zone/cms/config"
	"github.com/Open-Code-Zone/cms/internal/database"
	"github.com/Open-Code-Zone/cms/services/auth"
	"github.com/Open-Code-Zone/cms/utils"
	"github.com/Open-Code-Zone/cms/views/components"
	"github.com/Open-Code-Zone/cms/views/pages"
	"github.com/gorilla/mux"
)

// listing all the collection items for example blog posts, authors
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	collectionConfig := h.store.Collections.GetCollectionConfig(vars["collection"])
	userConfig := r.Context().Value(auth.UserContextKey).(*config.User) // error is not handled since it is already handled in auth middleware

	collectionPermissions := userConfig.GetCollectionPermission(collectionConfig.Collection)

	// we need to handle this permissioning through getting user and it's permission info from yaml file
	if collectionConfig == nil {
		http.Error(w, "Collection doesn't exist", http.StatusInternalServerError)
		return
	}
	queryBuilder, err := database.BuildCollectionQuery(r, collectionConfig)
	if err != nil {
		log.Printf("Error building query: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Execute the query using sqlc
	rows, err := h.store.DB.QueryContext(r.Context(), queryBuilder.Query, queryBuilder.Args...)
	log.Printf("Executing query: %s with args: %v", queryBuilder.Query, queryBuilder.Args)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Failed to filter collection items", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []database.ListAllCollectionItemsRow
	for rows.Next() {
		var item database.ListAllCollectionItemsRow
		if err := rows.Scan(&item.Filename, &item.Content, &item.Metadata, &item.CreatedAt); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		items = append(items, item)
	}

	// Render the filtered items using the templ component
	pages.ShowCollectionItems(items, collectionConfig, collectionPermissions).Render(r.Context(), w)
}

// rendering markdown editor with metadata form for creating new collection item
func (h *Handler) New(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionConfig := h.store.Collections.GetCollectionConfig(vars["collection"])
	if collectionConfig == nil {
		http.Error(w, "Collection doesn't exists", http.StatusInternalServerError)
		return
	}

	userConfig := r.Context().Value(auth.UserContextKey).(*config.User) // error is not handled since it is already handled in auth middleware
	log.Println("name from INDEX", userConfig.Email)
	collectionPermissions := userConfig.GetCollectionPermission(collectionConfig.Collection)

	pages.EditCollection("new-draft.md", nil, collectionConfig, collectionPermissions).Render(r.Context(), w)
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
	fileContentWithFrontMatter := r.FormValue("content")
	message := fmt.Sprintf("🤖 - Added the new blog post: %s", fileName)
	filePath := filepath.Join(collectionConfig.GitPath, fileName)

	err := h.githubClient.CreateFile(filePath, fileContentWithFrontMatter, message)
	if err != nil {
		log.Println("Error updating file:", err)
		components.Toaster("Looks like same blog post with the same file name exists", "danger").Render(r.Context(), w)
		return
	}

	// db query after succesful request to github
	fileContent, metadata := utils.ExtractFrontMatter(fileContentWithFrontMatter)

	collectionItem := database.CreateCollectionItemParams{
		Filename:       fileName,
		CollectionName: collectionConfig.Collection,
		Content:        fileContent,
		Metadata:       metadata,
	}

	fmt.Println(collectionItem)

	a, err := h.store.Queries.CreateCollectionItem(r.Context(), h.store.DB, collectionItem)
	fmt.Println(a)
	if err != nil {
		log.Println("Error updating in db", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
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

	userConfig := r.Context().Value(auth.UserContextKey).(*config.User) // error is not handled since it is already handled in auth middleware
	log.Println("name from INDEX", userConfig.Email)
	collectionPermissions := userConfig.GetCollectionPermission(collectionConfig.Collection)

	params := database.GetCollectionItemParams{
		Filename:       id,
		CollectionName: collectionConfig.Collection,
	}

	collectionItem, err := h.store.Queries.GetCollectionItem(r.Context(), h.store.DB, params)
	if err != nil {
		log.Printf("Error getting markdown from db: %v", err)
		http.Error(w, "Failed to get markdown of file", http.StatusInternalServerError)
	}

	fileContent := utils.GenerateMarkdownFile(collectionItem)
	pages.EditCollection(id, &fileContent, collectionConfig, collectionPermissions).Render(r.Context(), w)
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
	fileContentWithFrontMatter := r.FormValue("content")

	fileContent, metadata := utils.ExtractFrontMatter(fileContentWithFrontMatter)

	params := database.UpdateCollectionItemParams{
		Filename:       fileName,
		Content:        fileContent,
		Metadata:       metadata,
		CollectionName: collectionConfig.Collection,
	}

	err := h.store.Queries.UpdateCollectionItem(r.Context(), h.store.DB, params)
	if err != nil {
		log.Printf("Error updating markdown in db: %v", err)
		http.Error(w, "Failed to update markdown of file", http.StatusInternalServerError)
	}

	message := fmt.Sprintf("🤖 - Updated the blog post: %s")
	filePath := filepath.Join(collectionConfig.GitPath, fileName)

	err = h.githubClient.UpdateFile(filePath, fileContent, message)
	if err != nil {
		log.Println("Error updating file:", err)
		http.Error(w, "Failed to update file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	params := database.DeleteCollectionItemParams{
		Filename:       fileName,
		CollectionName: collectionConfig.Collection,
	}

	err = h.store.Queries.DeleteCollectionItem(r.Context(), h.store.DB, params)
	if err != nil {
		log.Println("Error deleting file:", err)
		w.WriteHeader(http.StatusInternalServerError)
		components.Toaster("Couldn't able to delete through Database", "danger")
		return
	}

	components.Toaster("Blog Post deleted succesfully", "success").Render(r.Context(), w)
}

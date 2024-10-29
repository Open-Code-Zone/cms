package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Open-Code-Zone/cms/internal/database"
	"github.com/Open-Code-Zone/cms/utils"
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
	// fetching files from database
	markdownFiles, err := h.store.Queries.ListAllCollectionItems(r.Context(), h.store.DB, collectionConfig.Collection)
	if err != nil {
		log.Printf("Error fetching collection items: %v", err)
		http.Error(w, "Failed to fetch collection items", http.StatusInternalServerError)
		return
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

// Filter handles filtering collection items and returns HTMX compatible response
func (h *Handler) Filter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionConfig := h.store.Collections.GetCollectionConfig(vars["collection"])
	if collectionConfig == nil {
		http.Error(w, "Collection doesn't exist", http.StatusInternalServerError)
		return
	}

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Build the filter query based on metadata schema
	query := `
		SELECT filename, content, metadata, created_at
		FROM collections
		WHERE collection_name = ?
	`
	args := []interface{}{collectionConfig.Collection}

	// Get filterable fields from metadata schema
	for _, field := range collectionConfig.MetadataSchema {
		if !field.Filter {
			continue
		}

		value := r.Form.Get(field.Name)
		if value == "" {
			continue
		}

		switch field.Type {
		case "string":
			query += " AND json_extract(metadata, '$." + field.Name + "') LIKE ?"
			args = append(args, "%"+value+"%")
		case "datetime":
			dateValue := strings.Split(value, "T")[0]
			date, err := time.Parse("2006-01-02", dateValue)
			if err != nil {
				continue
			}
			query += " AND DATE(json_extract(metadata, '$." + field.Name + "')) = ?"
			args = append(args, date.Format("2006-01-02"))
		case "array":
			// Handle array types (authors, tags, etc.)
			values := r.Form[field.Name]
			if len(values) > 0 {
				placeholders := make([]string, len(values))
				for i, v := range values {
					// For SQLite, we can check if the string representation contains the value
					// This assumes the array is stored as a JSON array string
					placeholders[i] = "instr(json_extract(metadata, '$." + field.Name + "'), ?)"
					args = append(args, v)
				}
				query += " AND (" + strings.Join(placeholders, " > 0 OR ") + " > 0)"
			}
		}
	}

	// Add ordering
	query += " ORDER BY created_at DESC"

	// Execute the query using sqlc
	rows, err := h.store.DB.QueryContext(r.Context(), query, args...)
	log.Printf("Executing query: %s with args: %v", query, args)
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

	if len(items) == 0 {
		components.Toaster("No items found matching the filters", "info").Render(r.Context(), w)
		return
	}

	// Render the filtered items using the templ component
	components.FilteredItems(items, collectionConfig).Render(r.Context(), w)
}

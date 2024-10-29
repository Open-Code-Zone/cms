package main

import (
	"context"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"

	"github.com/Open-Code-Zone/cms/config"
	"github.com/Open-Code-Zone/cms/internal/database"
	"github.com/Open-Code-Zone/cms/store"
	"github.com/Open-Code-Zone/cms/utils"
)

func main() {
	ctx := context.Background()
	githubClient, err := utils.NewGitHubClient()
	if err != nil {
		log.Printf("Error creating GitHub client: %v", err)
	}

	db, err := store.NewSQLiteStorage("cms.db")
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New()

	for _, collection := range *config.Envs.CollectionConfig {
		contents, err := githubClient.GetContents(collection.GitPath)
		if err != nil {
			log.Printf("Error fetching repository contents: %v", err)
			return
		}
		for _, file := range contents {
			if file.GetType() == "file" && strings.HasSuffix(file.GetName(), ".md") {
				// db seed
				content, err := githubClient.GetFileContent(collection.GitPath + "/" + file.GetName())
				if err != nil {
					log.Println("error occured while getting the content of collection item file", err)
					return
				}
				markdown, frontMatter := utils.ExtractFrontMatter(content)
				params := database.CreateCollectionItemParams{
					Filename:       *file.Name,
					Content:        markdown,
					CollectionName: collection.Collection,
					Metadata:       frontMatter,
				}
				_, err = queries.CreateCollectionItem(ctx, db, params)
				if err != nil {
					log.Println("error occured while creating collection item", err)
					return
				}
			}
		}
	}
}

package database

import (
	"net/http"
	"strings"
	"time"

	"github.com/Open-Code-Zone/cms/config"
)

type QueryBuilder struct {
	Query string
	Args  []interface{}
}

// buildCollectionQuery constructs a SQL query with filters based on form data and collection config
func BuildCollectionQuery(r *http.Request, collectionConfig *config.Collection) (*QueryBuilder, error) {
	// Parse the form data if not already parsed
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	qb := &QueryBuilder{
		Query: `
			SELECT filename, content, metadata, created_at
			FROM collections
			WHERE collection_name = ?
		`,
		Args: []interface{}{collectionConfig.Collection},
	}

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
			qb.Query += " AND json_extract(metadata, '$." + field.Name + "') LIKE ?"
			qb.Args = append(qb.Args, "%"+value+"%")
		case "datetime":
			dateValue := strings.Split(value, "T")[0]
			date, err := time.Parse("2006-01-02", dateValue)
			if err != nil {
				continue
			}
			qb.Query += " AND DATE(json_extract(metadata, '$." + field.Name + "')) = ?"
			qb.Args = append(qb.Args, date.Format("2006-01-02"))
		case "array":
			values := r.Form[field.Name]
			if len(values) > 0 {
				placeholders := make([]string, len(values))
				for i, v := range values {
					placeholders[i] = "instr(json_extract(metadata, '$." + field.Name + "'), ?)"
					qb.Args = append(qb.Args, v)
				}
				qb.Query += " AND (" + strings.Join(placeholders, " > 0 OR ") + " > 0)"
			}
		}
	}

	// Add ordering
	qb.Query += " ORDER BY created_at DESC"

	return qb, nil
}

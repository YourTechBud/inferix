package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/YourTechBud/inferix/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/go-libsql"
	"github.com/valyala/fastjson"
)

// LibSQLConfigDriver manages the configuration stored in a sqlite database
type LibSQLConfigDriver struct {
	db            *sqlx.DB
	defaultConfig map[string]any
}

// NewLibSQLConfigDriver creates a new SqliteConfigDriver
func NewLibSQLConfigDriver(opts Options) (*LibSQLConfigDriver, error) {
	// First create the containing directory if it doesn't exist
	if err := utils.CreateDirIfNotExists(opts.ConfigPath); err != nil {
		return nil, err
	}

	// Open the database
	db, err := sqlx.Open("libsql", fmt.Sprintf("file://%s", opts.ConfigPath))
	if err != nil {
		return nil, err
	}

	// Load the default configuration if provided
	var defaultConfig map[string]any
	if opts.DefaultConfigPath != "" {
		// Read the default configuration
		if err := utils.ReadYAMLFile(opts.DefaultConfigPath, &defaultConfig); err != nil {
			return nil, err
		}
	}

	// Create the table
	if _, err := db.Exec(sqliteConfigSchema); err != nil {
		return nil, err
	}

	return &LibSQLConfigDriver{
		db:            db,
		defaultConfig: defaultConfig,
	}, nil
}

func (s *LibSQLConfigDriver) Close() error {
	return s.db.Close()
}

func (s *LibSQLConfigDriver) ReadAll(ctx context.Context) (json.RawMessage, error) {
	elems := sqliteConfigElements{}
	if err := s.db.SelectContext(ctx, &elems, "SELECT * FROM config ORDER BY module, path, element_id"); err != nil {
		return nil, err
	}

	// Create a map of the elements
	cfg := elems.convertToMap()

	// Add all the modules if they don't exist
	if _, p := cfg["llm"]; !p {
		cfg["llm"] = map[string]any{}
	}

	// Merge with the default configuration
	// The default config always has the highest priority
	utils.MergeMaps(cfg, s.defaultConfig)

	// Convert the map to JSON
	return json.Marshal(cfg)
}

func (s *LibSQLConfigDriver) Get(ctx context.Context, module, path string) (json.RawMessage, error) {
	elems, err := s.getElementsByModuleAndPath(ctx, module, path)
	if err != nil {
		return nil, err
	}

	return elems.getValueAtPath(module, path)
}

func (s *LibSQLConfigDriver) SetInArray(ctx context.Context, module, path, id string, element interface{}) error {
	return s.setElement(ctx, module, path, id, "ARRAY", element)
}

func (s *LibSQLConfigDriver) SetInObject(ctx context.Context, module, path, id string, element interface{}) error {
	return s.setElement(ctx, module, path, id, "OBJECT", element)
}

func (s *LibSQLConfigDriver) DeleteFromArray(ctx context.Context, module, path, id string) error {
	return s.deleteElement(ctx, module, path, id)
}

func (s *LibSQLConfigDriver) DeleteFromObject(ctx context.Context, module, path, id string) error {
	return s.deleteElement(ctx, module, path, id)
}

func (s *LibSQLConfigDriver) getElementsByModuleAndPath(ctx context.Context, module, path string) (sqliteConfigElements, error) {
	elems := sqliteConfigElements{}
	if err := s.db.SelectContext(ctx, &elems, "SELECT * FROM config WHERE module = ? AND path = ? ORDER BY element_id", module, path); err != nil {
		return nil, err
	}

	return elems, nil
}

func (s *LibSQLConfigDriver) setElement(ctx context.Context, module, path, id, elementType string, element interface{}) error {
	query := `INSERT INTO config (module, path, element_id, element_type, config) VALUES (?, ?, ?, ?, ?)
	ON CONFLICT (module, path, element_id) DO UPDATE SET config = ?`

	// Convert the element to JSON
	data, err := json.Marshal(element)
	if err != nil {
		return err
	}

	// Execute the query
	_, err = s.db.ExecContext(ctx, query, module, path, id, elementType, string(data), string(data))
	return err
}

func (s *LibSQLConfigDriver) deleteElement(ctx context.Context, module, path, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM config WHERE module = ? AND path = ? AND element_id = ?", module, path, id)
	return err
}

type (
	sqliteConfigElement struct {
		Module      string `db:"module"`
		Path        string `db:"path"`
		ElementID   string `db:"element_id"`
		ElementType string `db:"element_type"`
		Config      string `db:"config"`
	}

	sqliteConfigElements []sqliteConfigElement
)

func (c sqliteConfigElements) convertToMap() map[string]any {
	// Create a map of the elements
	cfg := map[string]any{}
	for _, elem := range c {
		// Start by getting the module
		moduleT, p := cfg[elem.Module]
		if !p {
			moduleT = map[string]any{}
			cfg[elem.Module] = moduleT
		}
		module := moduleT.(map[string]any)

		// Get the path
		var obj map[string]any
		var next any
		pathArr := strings.Split(elem.Path, "/")
		pathLength := len(pathArr)
		for i, pathElement := range pathArr {
			if i == 0 {
				obj = module
			} else {
				obj = next.(map[string]any)
			}

			// Check if the path exists
			var p bool
			next, p = obj[pathElement]
			if !p {
				if i == pathLength-1 && elem.ElementType == "ARRAY" {
					// This is the last element and it is an array
					next = []any{}
				} else {
					// All intermediate elements have to be objects
					next = map[string]any{}
				}
				obj[pathElement] = next
			}
		}

		if elem.ElementType == "OBJECT" {
			next.(map[string]any)[elem.ElementID] = json.RawMessage(elem.Config)
		} else {
			// Check if any element in the array has the provided id
			found := false
			arr := next.([]any)
			for i, v := range arr {
				id := fastjson.GetString(v.(json.RawMessage), "id")
				if id == elem.ElementID {
					// Update the element
					arr[i] = json.RawMessage(elem.Config)
					found = true
					break
				}
			}

			if !found {
				// Add the element to the same slice so we don't have to update the map
				arr = append(arr, json.RawMessage(elem.Config))
				obj[pathArr[pathLength-1]] = arr
			}
		}
	}

	return cfg
}

func (c sqliteConfigElements) getValueAtPath(module, path string) (json.RawMessage, error) {
	// Create a map of the elements
	cfg := c.convertToMap()

	// Get the leaf element
	if _, p := cfg[module]; !p {
		return nil, nil
	}

	// Get the leaf element
	value := utils.GetValueAtPath(cfg[module].(map[string]any), path, nil)

	// Check if the value is nil
	if value == nil {
		return nil, nil
	}

	// Check if the value is a JSON message
	if jsonValue, ok := value.(json.RawMessage); ok {
		return jsonValue, nil
	}

	return json.Marshal(value)
}

var sqliteConfigSchema = `
CREATE TABLE IF NOT EXISTS config (
	module TEXT NOT NULL,					-- The module name
	path TEXT NOT NULL,						-- The path to the module
	element_id TEXT NOT NULL,				-- The id of the configuration resource
	element_type TEXT NOT NULL,				-- "OBJECT" or "ARRAY"
	config TEXT NOT NULL,					-- The configuration in JSON format
	PRIMARY KEY (module, path, element_id)
) STRICT;
`

var _ ConfigDriver = (*LibSQLConfigDriver)(nil)

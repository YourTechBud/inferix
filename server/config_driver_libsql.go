package server

import (
	"context"
	"encoding/json"
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
func NewLibSQLConfigDriver(db *sqlx.DB, defaultConfig map[string]any) (*LibSQLConfigDriver, error) {
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
	return nil
}

func (s *LibSQLConfigDriver) ReadAll(ctx context.Context) (json.RawMessage, error) {
	elems := sqliteConfigElements{}
	if err := s.db.SelectContext(ctx, &elems, "SELECT * FROM config ORDER BY module, path, element_id"); err != nil {
		return nil, err
	}

	// Create a map of the elements
	cfg := elems.convertToMap()

	// Add all the modules if they don't exist
	for _, m := range modulesList {
		if _, p := cfg[m.Name]; !p {
			cfg[m.Name] = map[string]any{}
		}
	}

	// Merge with the default configuration
	// The default config always has the highest priority
	utils.MergeMaps(cfg, s.defaultConfig)

	// Convert the map to JSON
	return json.Marshal(cfg)
}

func (s *LibSQLConfigDriver) GetAllResources(ctx context.Context, module, path string) ([]*utils.ResourceObject, error) {
	elems, err := s.getElementsByModuleAndPath(ctx, module, path)
	if err != nil {
		return nil, err
	}

	// Convert the elements to resources
	resources := make([]*utils.ResourceObject, len(elems))
	for i, elem := range elems {
		resources[i] = &utils.ResourceObject{
			Config:    json.RawMessage(elem.Config),
			Metadata:  json.RawMessage(elem.Metadata),
			CreatedAt: elem.CreatedAt,
			UpdatedAt: elem.UpdatedAt,
		}
	}

	return resources, nil
}

func (s *LibSQLConfigDriver) GetResource(ctx context.Context, module, path, id string) (*utils.ResourceObject, bool, error) {
	elem, found, err := s.getElementsByID(ctx, module, path, id)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}

	return &utils.ResourceObject{
		Config:    json.RawMessage(elem.Config),
		Metadata:  json.RawMessage(elem.Metadata),
		CreatedAt: elem.CreatedAt,
		UpdatedAt: elem.UpdatedAt,
	}, true, nil
}

func (s *LibSQLConfigDriver) CheckIfResourceExists(ctx context.Context, module, path, id string) (bool, error) {
	var count int
	if err := s.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM config WHERE module = ? AND path = ? AND element_id = ?", module, path, id); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *LibSQLConfigDriver) SetResource(ctx context.Context, module, path, id string, element, metadata any) error {
	return s.setElement(ctx, module, path, id, element, metadata)
}

func (s *LibSQLConfigDriver) SetResourceMetadata(ctx context.Context, module, path, id string, metadata any) error {
	// Convert the metadata to JSON
	metadataData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, "UPDATE config SET metadata = ? WHERE module = ? AND path = ? AND element_id = ?", string(metadataData), module, path, id)
	return err
}

func (s *LibSQLConfigDriver) DeleteResource(ctx context.Context, module, path, id string) error {
	return s.deleteElement(ctx, module, path, id)
}

func (s *LibSQLConfigDriver) getElementsByModuleAndPath(ctx context.Context, module, path string) (sqliteConfigElements, error) {
	elems := sqliteConfigElements{}
	if err := s.db.SelectContext(ctx, &elems, "SELECT * FROM config WHERE module = ? AND path = ? ORDER BY element_id", module, path); err != nil {
		return nil, err
	}

	return elems, nil
}

func (s *LibSQLConfigDriver) getElementsByID(ctx context.Context, module, path, id string) (sqliteConfigElement, bool, error) {
	elems := sqliteConfigElements{}
	if err := s.db.SelectContext(ctx, &elems, "SELECT * FROM config WHERE module = ? AND path = ? AND element_id = ? ORDER BY element_id", module, path, id); err != nil {
		return sqliteConfigElement{}, false, err
	}

	if len(elems) == 0 {
		return sqliteConfigElement{}, false, nil
	}

	return elems[0], true, nil
}

func (s *LibSQLConfigDriver) setElement(ctx context.Context, module, path, id string, element, metadata any) error {
	query := `INSERT INTO config (module, path, element_id, config, metadata, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT (module, path, element_id) DO UPDATE SET config = ?, metadata = ?, updated_at = ?`

	// Make sure metadata is not nil
	if metadata == nil {
		metadata = map[string]any{}
	}

	// Convert the element to JSON
	elementData, err := json.Marshal(element)
	if err != nil {
		return err
	}

	// Convert the metadata to JSON
	metadataData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// Execute the query
	currentTime := utils.CurrentTime()
	_, err = s.db.ExecContext(ctx, query, module, path, id, string(elementData), string(metadataData), currentTime, currentTime, string(elementData), string(metadataData), currentTime)
	return err
}

func (s *LibSQLConfigDriver) deleteElement(ctx context.Context, module, path, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM config WHERE module = ? AND path = ? AND element_id = ?", module, path, id)
	return err
}

type (
	sqliteConfigElement struct {
		Module    string `db:"module"`
		Path      string `db:"path"`
		ElementID string `db:"element_id"`
		Config    string `db:"config"`
		Metadata  string `db:"metadata"`
		CreatedAt int64  `db:"created_at"`
		UpdatedAt int64  `db:"updated_at"`
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
				if i == pathLength-1 {
					// This is the last element and it is an array
					next = []any{}
				} else {
					// All intermediate elements have to be objects
					next = map[string]any{}
				}
				obj[pathElement] = next
			}
		}

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

	return cfg
}

var sqliteConfigSchema = `
CREATE TABLE IF NOT EXISTS config (
	module TEXT NOT NULL,					-- The module name
	path TEXT NOT NULL,						-- The path to the module
	element_id TEXT NOT NULL,				-- The id of the configuration resource
	config TEXT NOT NULL,					-- The configuration in JSON format
	metadata TEXT NOT NULL,					-- The metadata in JSON format
	created_at INTEGER NOT NULL,			-- The time the configuration was created
	updated_at INTEGER NOT NULL,			-- The time the configuration was last updated
	PRIMARY KEY (module, path, element_id)
) STRICT;
`

var _ utils.ConfigDriver = (*LibSQLConfigDriver)(nil)

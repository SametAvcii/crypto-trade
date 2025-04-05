package query

var Query = map[string]string{
	"id":                 "id = ?",
	"user_id":            "user_id = ?",
	"name":               "name = ?",
	"email":              "email = ?",
	"entity_id":          "entity_id = ?",
	"entity_id_in":       "entity_id IN (?)",
	"entity":             "entity = ?",
	"entity_type":        "entity_type = ?",
	"created_at_between": "created_at BETWEEN ? AND ?",
}

const (
	WhereID               string = "id"
	WhereCreatedAtBetween string = "created_at_between"
	WhereUserID           string = "user_id"
	WhereEntityID         string = "entity_id"
	WhereEntityIDIN       string = "entity_id_in"
	WhereEntity           string = "entity"
	WhereEntityType       string = "entity_type"
	WhereEmail            string = "email"
	WhereName             string = "name"
)

func BuildQuery(keys ...string) string {
	if len(keys) == 0 {
		return ""
	}
	if len(keys) == 1 {
		return Query[keys[0]]
	}

	var query string
	for _, key := range keys {
		query += Query[key] + " AND "
	}

	return query[:len(query)-5] // remove last " AND "
}

type QueryBuilder struct {
	Keys []QueryKey `json:"keys"`
}

type QueryKey struct {
	Key    string
	Values []interface{}
	Skip   bool
}

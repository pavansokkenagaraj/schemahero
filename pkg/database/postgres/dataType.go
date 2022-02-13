package postgres

import(
	"strings"

	schemasv1alpha4 "github.com/schemahero/schemahero/pkg/apis/schemas/v1alpha4"
	"github.com/schemahero/schemahero/pkg/database/types"
)
const(
	columnENUM = "ENUM"
)


// iteration -2
// columns:
// - name: current_mood
//   type: ENUM ('sad', 'ok', 'happy')
// 	   - name: mood (optional)

// CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy');
// CREATE TABLE person (
//     name text,
//     current_mood mood
// );

func isStaticColumnType(requestedType string)bool{
	return strings.HasPrefix(requestedType, columnENUM)
}


func staticSchemaColumnToColumn(schemaColumn *schemasv1alpha4.PostgresqlTableColumn) (*types.Column, error) {
	column := &types.Column{
		Name:          schemaColumn.Name,
		DataType: schemaColumn.Name,
		IsStatic: true,
		StaticDataType: schemaColumn.Type,
	}
	return column, nil
}
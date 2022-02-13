package postgres

import (
	"reflect"
	"testing"

	schemasv1alpha4 "github.com/schemahero/schemahero/pkg/apis/schemas/v1alpha4"
	"github.com/schemahero/schemahero/pkg/database/types"
)

func Test_staticSchemaColumnToColumn(t *testing.T) {
	type args struct {
		schemaColumn *schemasv1alpha4.PostgresqlTableColumn
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Column
		wantErr bool
	}{
		{
			name: `expect ENUM () to convert`,
			args: args{
				schemaColumn: &schemasv1alpha4.PostgresqlTableColumn{
					Name: "mood",
					Type: "ENUM ('sad', 'ok', 'happy')",
				},
			},
			want: &types.Column{
				Name: "mood",
				DataType: "mood",
				IsStatic: true,
				StaticDataType: "ENUM ('sad', 'ok', 'happy')",
			},
		},{
			name: `expect ENUM () to convert`,
			args: args{
				schemaColumn: &schemasv1alpha4.PostgresqlTableColumn{
					Name: "mood",
					Type: "ENUM ()",
				},
			},
			want: &types.Column{
				Name: "mood",
				DataType: "mood",
				IsStatic: true,
				StaticDataType: "ENUM ()",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := staticSchemaColumnToColumn(tt.args.schemaColumn)
			if (err != nil) != tt.wantErr {
				t.Errorf("staticSchemaColumnToColumn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("staticSchemaColumnToColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

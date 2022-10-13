package storage

import "github.com/lfritz/toydb/types"

type SampleData struct {
	Database      *Database
	Films, People *types.Relation
}

func GetSampleData() *SampleData {
	filmsSchema := types.TableSchema{
		Columns: []types.ColumnSchema{
			types.ColumnSchema{"id", types.TypeDecimal},
			types.ColumnSchema{"name", types.TypeText},
			types.ColumnSchema{"release_date", types.TypeDate},
			types.ColumnSchema{"director", types.TypeDecimal},
		},
	}
	filmsRows := [][]types.Value{
		{types.NewDecimal("1"), types.NewText("The General"), types.NewDate(1926, 12, 31), types.NewDecimal("1")},
		{types.NewDecimal("2"), types.NewText("The Kid"), types.NewDate(1921, 1, 21), types.NewDecimal("2")},
		{types.NewDecimal("3"), types.NewText("Sherlock Jr."), types.NewDate(1924, 4, 21), types.NewDecimal("1")},
	}
	films := &types.Relation{
		Schema: filmsSchema,
		Rows:   filmsRows,
	}

	peopleSchema := types.TableSchema{
		Columns: []types.ColumnSchema{
			types.ColumnSchema{"id", types.TypeDecimal},
			types.ColumnSchema{"name", types.TypeText},
		},
	}
	peopleRows := [][]types.Value{
		{types.NewDecimal("1"), types.NewText("Buster Keaton")},
		{types.NewDecimal("2"), types.NewText("Charlie Chaplin")},
	}
	people := &types.Relation{
		Schema: peopleSchema,
		Rows:   peopleRows,
	}

	database := &Database{
		tables: map[string]*types.Relation{
			"films":  films,
			"people": people,
		},
	}

	return &SampleData{
		Database: database,
		Films:    films,
		People:   people,
	}
}

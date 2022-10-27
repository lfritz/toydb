package storage

import "github.com/lfritz/toydb/types"

type SampleData struct {
	Database      *Database
	Films, People *types.Relation
}

func GetSampleData() *SampleData {
	filmsSchema := types.TableSchema{
		Columns: []types.ColumnSchema{
			types.ColumnSchema{"id", types.TypeDecimal, false},
			types.ColumnSchema{"name", types.TypeText, false},
			types.ColumnSchema{"release_date", types.TypeDate, false},
			types.ColumnSchema{"director", types.TypeDecimal, false},
		},
	}
	filmsRows := [][]types.Value{
		{types.Dec("1"), types.Txt("The General"), types.Dat(1926, 12, 31), types.Dec("1")},
		{types.Dec("2"), types.Txt("The Kid"), types.Dat(1921, 1, 21), types.Dec("2")},
		{types.Dec("3"), types.Txt("Sherlock Jr."), types.Dat(1924, 4, 21), types.Dec("1")},
	}
	films := &types.Relation{
		Schema: filmsSchema,
		Rows:   filmsRows,
	}

	peopleSchema := types.TableSchema{
		Columns: []types.ColumnSchema{
			types.ColumnSchema{"id", types.TypeDecimal, false},
			types.ColumnSchema{"name", types.TypeText, false},
		},
	}
	peopleRows := [][]types.Value{
		{types.Dec("1"), types.Txt("Buster Keaton")},
		{types.Dec("2"), types.Txt("Charlie Chaplin")},
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

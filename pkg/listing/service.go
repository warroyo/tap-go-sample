package listing

import "github.com/warroyo/tap-go-sample/pkg/database"

func GetAllCompanies() []Company {
	db := database.ConnectToDB(false)
	defer db.Close()

	results, err := db.Query(`SELECT id, name FROM companies`)

	if err != nil {
		panic(err.Error())
	}

	defer results.Close()

	var cs []Company

	for results.Next() {
		var c Company

		err = results.Scan(
			&c.Id,
			&c.Name)
		if err != nil {
			panic(err.Error())
		}

		cs = append(cs, c)
	}

	return cs
}

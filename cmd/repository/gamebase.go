package repository

import (
	"log"

	f "github.com/fauna/faunadb-go/v4/faunadb"
)

type GameBase struct {
	ID        int
	Number    int
	RolesList []string
}

// func GetByID(filename string, id int) (gb *GameBase, err error) {
// 	// Open our csvFile
// 	csvFile, err := os.Open(filename)
// 	// if we os.Open returns an error then handle it
// 	if err != nil {
// 		log.Println(err)
// 		return nil, errors.New("False open gamebase file")
// 	}
// 	csvReader := csv.NewReader(csvFile)
// 	data, err := csvReader.ReadAll()
// 	if err != nil {
// 		log.Println(err)
// 		return nil, errors.New("False read gamebase file")
// 	}
// 	log.Println("read file OK")
// 	log.Println(data)
// 	// parse to RoleList
// 	for i, line := range data {
// 		if i == id-1 {
// 			var rec GameBase
// 			for index, column := range line {
// 				switch index {
// 				case 0:
// 					rec.ID, _ = strconv.Atoi(column)
// 				case 1:
// 					rec.Number, _ = strconv.Atoi(column)
// 				case 2:
// 					rec.RolesList = strings.Split(column, ",")
// 				}
// 			}
// 			gb = &rec
// 			break
// 		}
// 	}
// 	time.Sleep(3)
// 	if gb == nil {
// 		return nil, errors.New("No gamebase avaiable")
// 	}
// 	// defer the closing of our jsonFile so that we can parse it later on
// 	defer csvFile.Close()
// 	return gb, nil
// }

func GetByID(db_key string, id int) (gb *GameBase, err error) {
	client := f.NewFaunaClient(db_key, f.Endpoint("https://db.us.fauna.com/"))

	res, err := client.Query(f.Get(f.Ref(f.Collection("gamebase"), "357184493348454484")))
	if err != nil {
		return nil, err
	}

	if err := res.At(f.ObjKey("data")).Get(&gb); err != nil {
		return nil, err
	}

	log.Println(gb)
	return gb, nil
}

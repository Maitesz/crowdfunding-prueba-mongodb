package model

import "github.com/globalsign/mgo/bson"

// Region : proyecto de un usuario
type Region struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Name  string        `json:"name" bson:"name"`
	Comunas []Comuna `json:"comunas" bson:"comunas"`
	}
type Comuna struct {
	Name  string        `json:"name" bson:"name"`
}
// Page : Pagina de resultado


// Create : Crear proyecto por ID
func (regionModel *Region) Create(regionDoc *Region) error {
	col, session := GetCollection(CollectionNameRegion)
	defer session.Close()
	regionDoc.ID = bson.NewObjectId()
	err := col.Insert(regionDoc)

	return err
}

// Get : Obtener proyecto por ID
func (regionModel *Region) Get(id string) (*Region, error) {
	col, session := GetCollection(CollectionNameRegion)
	defer session.Close()
	var regionDoc Region
	err := col.FindId(bson.ObjectIdHex(id)).One(&regionDoc)

	return &regionDoc, err
}

// Update : Actualizar proyecto por ID
func (regionModel *Region) Update(id string, regionDoc Region) error {

	col, session := GetCollection(CollectionNameRegion)
	defer session.Close()
	err := col.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": regionDoc})
	return err
}

// Delete : Eliminar proyecto por ID
func (regionModel *Region) Delete(id string) error {

	col, session := GetCollection(CollectionNameRegion)
	defer session.Close()
	err := col.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return err
}

// Find : Obtener proyecto
func (regionModel *Region) Find(query bson.M) ([]Region, error) {

	col, session := GetCollection(CollectionNameRegion)
	defer session.Close()
	regions := []Region{}

	err := col.Find(query).All(&regions)
	return regions, err
}

// FindPaginate : Obtener proyecto
func (regionModel *Region) FindPaginate(query bson.M, limit int, offset int) (Page, error) {

	col, session := GetCollection(CollectionNameRegion)
	defer session.Close()
	pag := []bson.M{{"$skip": offset}}
	if limit > 0 {
		pag = append(pag, bson.M{"$limit": limit})
	}
	pipeline := []bson.M{
		bson.M{"$match": query},
		bson.M{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total"}},
			"data":     pag, // add region here wish you re-shape the docs
		}},
	}

	pageDoc := Page{}
	err := col.Pipe(pipeline).One(&pageDoc)

	return pageDoc, err
}

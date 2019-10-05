package model

import "github.com/globalsign/mgo/bson"

// Project : proyecto de un usuario
type Project struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Name  string        `json:"name" bson:"name"`
	Giro   string           `json:"giro" bson:"giro,omitempty"`
	Description   string           `json:"description" bson:"description,omitempty"`
	Object   string           `json:"object" bson:"object,omitempty"`
	Public_object   string           `json:"public_object" bson:"public_object,omitempty"`
	Owner bson.ObjectId `json:"owner" bson:"owner"`
	Sponsor bson.ObjectId `json:"sponsor" bson:"sponsor,omitempty"`
}

// Page : Pagina de resultado
type Page struct {
	Metadata []map[string]int `json:"metadata" bson:"metadata,omitempty"`
	Data     []interface{}    `json:"data" bson:"data,omitempty"`
}

// Create : Crear proyecto por ID
func (projectModel *Project) Create(projectDoc *Project) error {
	col, session := GetCollection(CollectionNameProject)
	defer session.Close()
	projectDoc.ID = bson.NewObjectId()
	err := col.Insert(projectDoc)

	return err
}

// Get : Obtener proyecto por ID
func (projectModel *Project) Get(id string) (*Project, error) {
	col, session := GetCollection(CollectionNameProject)
	defer session.Close()
	var projectDoc Project
	err := col.FindId(bson.ObjectIdHex(id)).One(&projectDoc)

	return &projectDoc, err
}

// Update : Actualizar proyecto por ID
func (projectModel *Project) Update(id string, projectDoc Project) error {

	col, session := GetCollection(CollectionNameProject)
	defer session.Close()
	err := col.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": projectDoc})
	return err
}

// Delete : Eliminar proyecto por ID
func (projectModel *Project) Delete(id string) error {

	col, session := GetCollection(CollectionNameProject)
	defer session.Close()
	err := col.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return err
}

// Find : Obtener proyecto
func (projectModel *Project) Find(query bson.M) ([]Project, error) {

	col, session := GetCollection(CollectionNameProject)
	defer session.Close()
	projects := []Project{}

	err := col.Find(query).All(&projects)
	return projects, err
}

// FindPaginate : Obtener proyecto
func (projectModel *Project) FindPaginate(query bson.M, limit int, offset int) (Page, error) {

	col, session := GetCollection(CollectionNameProject)
	defer session.Close()
	pag := []bson.M{{"$skip": offset}}
	if limit > 0 {
		pag = append(pag, bson.M{"$limit": limit})
	}
	pipeline := []bson.M{
		bson.M{"$match": query},
		bson.M{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total"}},
			"data":     pag, // add projection here wish you re-shape the docs
		}},
	}

	pageDoc := Page{}
	err := col.Pipe(pipeline).One(&pageDoc)

	return pageDoc, err
}

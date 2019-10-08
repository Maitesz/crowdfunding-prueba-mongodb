package controller

import (
"fmt"
"net/http"
jwt "github.com/appleboy/gin-jwt/v2"
"Crowdfunding_memoria/crowdfunding prueba mongodb/middleware"
"Crowdfunding_memoria/crowdfunding prueba mongodb/model"
"Crowdfunding_memoria/crowdfunding prueba mongodb/util"
"github.com/gin-gonic/gin"

"github.com/globalsign/mgo/bson"
)

// DogController : Controlador de perro
type ProjectController struct {
}

// Routes : Define las rutas del controlador
func (projectController *ProjectController) Routes(base *gin.RouterGroup, authNormal *jwt.GinJWTMiddleware) *gin.RouterGroup {

	// Projects - Rutas
	projectRouter := base.Group("/projects") //, middleware.SetRoles(RolAdmin, RolUser), authNormal.MiddlewareFunc())
	{
		projectRouter.GET("", projectController.GetAll())
		// Al agregar asociar con usuario
		projectRouter.POST("", authNormal.MiddlewareFunc(), projectController.Create())
		projectRouter.GET("/:id", projectController.One())
		// Verificar en handler que el perro sea dueño de usuario
		projectRouter.PUT("/:id", authNormal.MiddlewareFunc(), projectController.Update())
		// Solo admin puede eliminar
		projectRouter.DELETE("/:id", middleware.SetRoles(RolAdmin), authNormal.MiddlewareFunc(), projectController.Delete())
	}
	return projectRouter
}

// GetAll : Obtener todos los perros
func (projectController *ProjectController) GetAll() func(c *gin.Context) {
	return func(c *gin.Context) {
		/* obtener parametros de paginacion*/
		pagination := PaginationParams{}
		err := c.ShouldBind(&pagination)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se puedieron encontrar los parametros limit, offset", err))
			return
		}
		page, err := projectModel.FindPaginate(bson.M{}, pagination.Limit, pagination.Offset)

		if err != nil {
			c.JSON(http.StatusNotFound, util.GetError("No se pudo obtener la lista de perros", err))
		}
		// c.Header("",page.Metadata.)
		if len(page.Metadata) != 0 {
			c.Header("Pagination-Count", fmt.Sprintf("%d", page.Metadata[0]["total"]))
		}

		c.JSON(http.StatusOK, page.Data)
	}
}

// Create : Crear perro
func (projectController *ProjectController) Create() func(c *gin.Context) {
	return func(c *gin.Context) {

		// Traer Usuario
		user := userModel.LoadFromContext(c)
		var project model.Project
		err := c.Bind(&project)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo decodificar json", err))
			return
		}
		// Asignar owner
		project.Owner = user.ID
		err = projectModel.Create(&project)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo insertar perro", err))
			return
		}

		c.JSON(http.StatusOK, project)
	}
}

// One : Obtener perro por _id
func (projectController *ProjectController) One() func(c *gin.Context) {
	return func(c *gin.Context) {

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusNotFound, util.GetError("No se encuentra parametro :id", nil))
			return
		}
		if !bson.IsObjectIdHex(id) {
			c.JSON(http.StatusInternalServerError, util.GetError("El id ingresado no es válido", nil))
			return
		}
		group, err := projectModel.Get(id)
		if err != nil {
			c.JSON(http.StatusNotFound, util.GetError("No se encontró perro", err))
			return
		}
		c.JSON(http.StatusOK, group)
	}
}

// Update : Actualizar perro con _id
func (projectController *ProjectController) Update() func(c *gin.Context) {
	return func(c *gin.Context) {

		var project model.Project
		err := c.Bind(&project)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo convertir collection json", err))
			return
		}
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, util.GetError("No se encuentra parametro :id", nil))
			return
		}

		if !bson.IsObjectIdHex(id) {
			c.JSON(http.StatusInternalServerError, util.GetError("El id ingresado no es válido", nil))
			return
		}
		// Update
		err = projectModel.Update(id, project)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo actualizar perro", err))
			return
		}

		c.String(http.StatusOK, "")
	}
}

// Delete : Eliminar perro por _id
func (projectController *ProjectController) Delete() func(c *gin.Context) {
	return func(c *gin.Context) {

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, util.GetError("No se encuentra parametro :id", nil))
			return
		}
		if !bson.IsObjectIdHex(id) {
			c.JSON(http.StatusInternalServerError, util.GetError("El id ingresado no es válido", nil))
			return
		}
		err := projectModel.Delete(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo encontrar perro", err))
			return
		}
		c.String(http.StatusOK, "")
	}
}


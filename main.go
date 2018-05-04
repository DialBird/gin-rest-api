package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

// User is
type User struct {
	Gender int    `json:"gender"`
	Name   string `json:"name"`
}

func getRoot(c *gin.Context) {
	c.String(200, "OK %s well well done!", "Keisuke")
}

func main() {
	session, err := mgo.Dial("172.17.0.2:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	cl := session.DB("test").C("users")

	r := gin.Default()

	r.GET("/", getRoot)

	r.GET("/users", func(c *gin.Context) {
		results := []User{}
		err = cl.Find(bson.M{}).All(&results)
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(200, results)
	})

	r.POST("/users", func(c *gin.Context) {
		user := User{}
		user.Gender, _ = strconv.Atoi(c.PostForm("gender"))
		user.Name = c.PostForm("name")

		err := cl.Insert(&user)
		if err != nil {
			if mgo.IsDup(err) {
				fmt.Println("Duplicate key error")
			}
			if v, ok := err.(*mgo.LastError); ok {
				fmt.Println("Code:%d N:%d UpdatedExisting:%t WTimeout:%t Waited:%d", v.Code, v.N, v.UpdatedExisting, v.WTimeout, v.Waited)
			} else {
				fmt.Println("%+v", err)
			}
		}
	})

	r.PUT("/users", func(c *gin.Context) {
		selector := bson.M{"name": c.PostForm("name")}
		update := bson.M{"$set": bson.M{"name": c.PostForm("newName")}}
		err := cl.Update(selector, update)
		if v, ok := err.(*mgo.LastError); ok {
			fmt.Println("Code:%d N:%d UpdatedExisting:%t WTimeout:%t Waited:%d", v.Code, v.N, v.UpdatedExisting, v.WTimeout, v.Waited)
		} else {
			fmt.Println("%+v", err)
		}
	})

	r.DELETE("/users", func(c *gin.Context) {
		selector := bson.M{"name": c.PostForm("name")}
		err := cl.Remove(selector)
		if v, ok := err.(*mgo.LastError); ok {
			fmt.Println("Code:%d N:%d UpdatedExisting:%t WTimeout:%t Waited:%d", v.Code, v.N, v.UpdatedExisting, v.WTimeout, v.Waited)
		} else {
			fmt.Println("%+v", err)
		}
	})

	r.Run()
}

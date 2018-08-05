/*

Copyright (C) 2017-2018  Ettore Di Giacinto <mudler@gentoo.org>
Credits goes also to Gogs authors, some code portions and re-implemented design
are also coming from the Gogs project, which is using the go-macaron framework
and was really source of ispiration. Kudos to them!

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/

package userapi

import (
	"github.com/MottainaiCI/mottainai-server/pkg/context"
	database "github.com/MottainaiCI/mottainai-server/pkg/db"

	user "github.com/MottainaiCI/mottainai-server/pkg/user"

	macaron "gopkg.in/macaron.v1"
)

func SetManager(ctx *context.Context, db *database.Database) error {
	id := ctx.ParamsInt(":id")

	u, err := db.GetUser(id)
	if err != nil {
		ctx.NotFound()
		return err
	}

	u.MakeManager()

	err = db.UpdateUser(id, u.ToMap())
	if err != nil {
		ctx.NotFound()
		return err
	}

	return nil
}
func SetAdmin(ctx *context.Context, db *database.Database) error {
	id := ctx.ParamsInt(":id")

	u, err := db.GetUser(id)
	if err != nil {
		ctx.NotFound()
		return err
	}

	u.MakeAdmin()

	err = db.UpdateUser(id, u.ToMap())
	if err != nil {
		ctx.NotFound()
		return err
	}

	return nil
}

func SetManagerUser(ctx *context.Context, db *database.Database) string {
	err := SetManager(ctx, db)
	if err != nil {
		return ":("
	}

	return "OK"
}

func SetAdminUser(ctx *context.Context, db *database.Database) string {
	err := SetAdmin(ctx, db)
	if err != nil {
		return ":("
	}

	return "OK"
}

func UnSetManager(ctx *context.Context, db *database.Database) error {
	id := ctx.ParamsInt(":id")

	u, err := db.GetUser(id)
	if err != nil {
		ctx.NotFound()
		return err
	}

	u.RemoveManager()

	err = db.UpdateUser(id, u.ToMap())
	if err != nil {
		ctx.NotFound()
		return err
	}
	return nil
}

func UnSetAdmin(ctx *context.Context, db *database.Database) error {
	id := ctx.ParamsInt(":id")

	u, err := db.GetUser(id)
	if err != nil {
		ctx.NotFound()
		return err
	}

	u.RemoveAdmin()

	err = db.UpdateUser(id, u.ToMap())
	if err != nil {
		ctx.NotFound()
		return err
	}
	return nil
}
func UnSetAdminUser(ctx *context.Context, db *database.Database) string {
	err := UnSetAdmin(ctx, db)
	if err != nil {
		return ":("
	}

	return "OK"
}
func UnSetManagerUser(ctx *context.Context, db *database.Database) string {
	err := UnSetManager(ctx, db)
	if err != nil {
		return ":("
	}

	return "OK"
}
func Delete(ctx *context.Context, db *database.Database) error {
	id := ctx.ParamsInt(":id")

	user, err := db.GetUser(id)
	if err != nil {
		ctx.NotFound()
		return err
	}

	tokens, _ := db.GetTokensByUserID(user.ID)

	for _, t := range tokens {
		db.DeleteToken(t.ID)
	}

	err = db.DeleteUser(id)
	if err != nil {
		ctx.NotFound()
		return err
	}
	return nil
}

func DeleteUser(ctx *context.Context, db *database.Database) string {

	err := Delete(ctx, db)
	if err != nil {
		return ":("
	}

	return "OK"
}

func List(c *context.Context, db *database.Database) []user.User {
	us := db.AllUsers()
	return us
}

func Show(c *context.Context, db *database.Database) (user.User, error) {
	id := c.ParamsInt(":id")
	u, err := db.GetUser(id)
	u.Password = ""
	if err != nil {
		return u, nil
	}
	return u, nil
}

func ShowUser(c *context.Context, db *database.Database) {
	u, err := Show(c, db)

	if err != nil {
		c.NotFound()
		return
	}
	c.JSON(200, u)

}

func ListUsers(c *context.Context, db *database.Database) {
	c.JSON(200, List(c, db))
}

func Setup(m *macaron.Macaron) {
	reqSignIn := context.Toggle(&context.ToggleOptions{SignInRequired: true})
	reqAdmin := context.Toggle(&context.ToggleOptions{AdminRequired: true})
	reqManager := context.Toggle(&context.ToggleOptions{ManagerRequired: true})

	m.Get("/api/user/list", reqManager, reqSignIn, ListUsers)
	m.Get("/api/user/show/:id", reqManager, reqSignIn, ShowUser)

	m.Get("/api/user/set/admin/:id", reqSignIn, reqAdmin, SetAdminUser)
	m.Get("/api/user/unset/admin/:id", reqSignIn, reqAdmin, UnSetAdminUser)
	m.Get("/api/user/set/manager/:id", reqSignIn, reqAdmin, SetManagerUser)
	m.Get("/api/user/unset/manager/:id", reqSignIn, reqAdmin, UnSetManagerUser)
	m.Get("/api/user/delete/:id", reqSignIn, reqAdmin, DeleteUser)
}

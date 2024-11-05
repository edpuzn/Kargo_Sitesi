package middleeware

import (
	"Cargo_Dash/database"
	"Cargo_Dash/models"
	"Cargo_Dash/util"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func IsAuthorized(c *fiber.Ctx, page string) error {
	cookie := c.Cookies("jwt")

	Id, err := util.ParseJWT(cookie)
	if err != nil {
		return err
	}

	fmt.Println(Id)

	userId, _ := strconv.Atoi(Id)

	user := models.User{
		Id: uint(userId),
	}

	database.DB.Preload("Role").Find(&user)

	role := models.Role{
		Id: user.RoleId,
	}

	err = database.DB.Preload("Permissions").Where("id = ?", Id).Find(&role).Error
	if err != nil {
		return err
	}

	fmt.Println(role.Permissions)
	if c.Method() == "GET" {
		for _, permission := range role.Permissions {
			if permission.Name == "view_"+page || permission.Name == "edit_"+page {
				return nil
			}
		}
	} else {
		for _, permission := range role.Permissions {
			fmt.Println(permission.Name)
			if permission.Name == "edit_"+page {
				return nil
			}
		}
	}

	c.Status(fiber.StatusUnauthorized)
	return errors.New("unauthorized")
}

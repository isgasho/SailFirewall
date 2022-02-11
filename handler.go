// Copyright (c) 2022 Hevienz
// Full license can be found in the LICENSE file.

package main

import (
    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
    "go.uber.org/zap"
)

var (
    validate = validator.New()
)

func AddRule(c *fiber.Ctx) error {
    key := &Key{}

    if err := c.BodyParser(key); err != nil {
        return ErrorJSON(c, err)
    }

    err := validate.Struct(key)
    if err != nil {
        return ErrorJSON(c, err)
    }

    err = aclist.Insert(key.GetBytes(), 0)
    if err != nil {
        return ErrorJSON(c,  err)
    }

    logger.Info("add rule", zap.Any("key", key))

	return SuccessJSON(c)
}

func GetRule(c *fiber.Ctx) error {
    key := &Key{}

    if err := c.BodyParser(key); err != nil {
        return ErrorJSON(c, err)
    }

    err := validate.Struct(key)
    if err != nil {
        return ErrorJSON(c, err)
    }

	count, err := aclist.LookupUint64(key.GetBytes())
	if err != nil {
        return ErrorJSON(c, err)
	}

    logger.Info("get rule", zap.Any("key", key), zap.Uint64("count", count))

	return SingleJSON(c, count)
}

func DeleteRule(c *fiber.Ctx) error {
    key := &Key{}

    if err := c.BodyParser(key); err != nil {
        return ErrorJSON(c, err)
    }

    err := validate.Struct(key)
    if err != nil {
        return ErrorJSON(c, err)
    }

    err = aclist.Delete(key.GetBytes())
    if err != nil {
        return ErrorJSON(c, err)
    }

    logger.Info("delete rule", zap.Any("key", key))

	return SuccessJSON(c)
}

// -----------------------------------------------

func SingleJSON(c *fiber.Ctx, data interface{}) error {
    return c.JSON(fiber.Map{
        "code": 200,
        "data": data,
    })
}

func ErrorJSON(c *fiber.Ctx, err error) error {
	return c.JSON(fiber.Map{
		"code":  500,
		"error": err.Error(),
	})
}

func SuccessJSON(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"code": 200,
	})
}

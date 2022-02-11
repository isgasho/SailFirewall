// Copyright (c) 2022 Hevienz
// Full license can be found in the LICENSE file.

package main

import (
    "github.com/gofiber/fiber/v2"
    "go.uber.org/zap"
)

func runApiServer() {
    app := fiber.New(fiber.Config{
        AppName: "SailFirewall",
    })

    apiv1 := app.Group("/api/v1")

    apiv1.Post("/rule", AddRule)
    apiv1.Get("/rule", GetRule)
    apiv1.Delete("/rule", DeleteRule)


    logger.Info("API server started", zap.String("bind_address", *apiAddr))
    err := app.Listen(*apiAddr)
    if err != nil {
        logger.Fatal("API server start failed", zap.Error(err))
    }
}

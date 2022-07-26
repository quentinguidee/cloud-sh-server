package com.quentinguidee.plugins

import io.ktor.server.application.*
import org.jetbrains.exposed.sql.Database

fun Application.configureDatabase() {
    Database.connect(
        "jdbc:postgresql://localhost:5432/cloudsh",
        driver = "org.postgresql.Driver",
        user = "cloudsh",
        password = "cloudsh"
    )
}

package com.quentinguidee.plugins

import com.quentinguidee.models.*
import io.ktor.server.application.*
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.SchemaUtils
import org.jetbrains.exposed.sql.transactions.transaction

fun Application.configureDatabase() {
    Database.connect(
        "jdbc:postgresql://localhost:5432/cloudsh",
        driver = "org.postgresql.Driver",
        user = "cloudsh",
        password = "cloudsh"
    )

    transaction {
        SchemaUtils.create(
            Buckets,
            GitHubUsers,
            Nodes,
            Sessions,
            Users,
            UsersBuckets,
        )
    }
}

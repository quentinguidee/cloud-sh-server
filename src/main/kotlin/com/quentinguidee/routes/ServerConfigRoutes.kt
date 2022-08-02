package com.quentinguidee.routes

import com.quentinguidee.plugins.DB_CONFIG_PATH
import com.quentinguidee.plugins.DatabaseConfig
import com.quentinguidee.plugins.connectDatabase
import com.quentinguidee.services.adminsServices
import com.quentinguidee.utils.ServerAlreadyConfiguredException
import com.quentinguidee.utils.ok
import io.ktor.server.application.*
import io.ktor.server.routing.*
import io.ktor.server.util.*

fun Route.serverConfigRoute() {
    route("/config") {
        post("/database") {
            if (DB_CONFIG_PATH.toFile().exists())
                throw ServerAlreadyConfiguredException()

            val config = DatabaseConfig(
                host = call.parameters.getOrFail("host"),
                name = call.parameters.getOrFail("name"),
                user = call.parameters.getOrFail("user"),
                password = call.parameters.getOrFail("password"),
            )

            connectDatabase(config)

            adminsServices.saveDatabaseConfig(config)

            call.ok()
        }
    }
}

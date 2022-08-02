package com.quentinguidee.routes

import com.quentinguidee.models.OAuthMethodPrivate
import com.quentinguidee.plugins.DB_CONFIG_PATH
import com.quentinguidee.plugins.DatabaseConfig
import com.quentinguidee.plugins.connectDatabase
import com.quentinguidee.services.adminsServices
import com.quentinguidee.services.authServices
import com.quentinguidee.utils.ServerAlreadyConfiguredException
import com.quentinguidee.utils.ok
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.util.*
import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put

fun Route.serverConfigRoute() {
    route("/config") {
        route("/database") {
            get {
                call.respond(buildJsonObject {
                    put("already_done", DB_CONFIG_PATH.toFile().exists())
                })
            }

            post {
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

        route("/oauth") {
            get {
                call.respond(buildJsonObject {
                    put("already_done", authServices.methods().isNotEmpty())
                })
            }

            post {
                if (authServices.methods().isNotEmpty())
                    throw ServerAlreadyConfiguredException()

                val config = OAuthMethodPrivate(
                    name = call.parameters.getOrFail("name").lowercase(),
                    displayName = call.parameters.getOrFail("name"),
                    color = call.parameters.getOrFail("color"),
                    clientID = call.parameters.getOrFail("client_id"),
                    clientSecret = call.parameters.getOrFail("client_secret"),
                    authorizeURL = call.parameters.getOrFail("authorize_url"),
                    accessTokenURL = call.parameters.getOrFail("access_token_url"),
                    redirectURL = call.parameters.getOrFail("redirect_url"),
                )

                authServices.createMethod(config)

                call.ok()
            }
        }
    }
}

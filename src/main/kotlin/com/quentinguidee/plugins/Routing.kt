package com.quentinguidee.plugins

import com.quentinguidee.routes.adminRoutes
import com.quentinguidee.routes.authRoutes
import com.quentinguidee.routes.storage.storageRoutes
import com.quentinguidee.routes.userRoutes
import com.quentinguidee.utils.authenticated
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Application.configureRouting() {
    install(IgnoreTrailingSlash)

    routing {
        get("/ping") {
            call.respond("pong")
        }

        authRoutes()

        authenticated {
            adminRoutes()
            storageRoutes()
            userRoutes()
        }
    }
}

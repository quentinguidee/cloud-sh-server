package com.quentinguidee.plugins

import com.quentinguidee.routes.authRoutes
import com.quentinguidee.routes.bucketRoutes
import com.quentinguidee.routes.userRoutes
import io.ktor.server.application.*
import io.ktor.server.routing.*

fun Application.configureRouting() {
    routing {
        route("/auth") {
            authRoutes()
        }

        route("/bucket") {
            bucketRoutes()
        }

        route("/user") {
            userRoutes()
        }
    }
}

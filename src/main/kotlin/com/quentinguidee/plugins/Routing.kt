package com.quentinguidee.plugins

import com.quentinguidee.routes.userRouting
import io.ktor.server.application.*
import io.ktor.server.routing.*

fun Application.configureRouting() {
    routing {
        route("/user") {
            userRouting()
        }
    }
}

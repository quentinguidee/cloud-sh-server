package com.quentinguidee.routes

import io.ktor.server.routing.*

fun Route.authRoutes() {
    route("/github") {
        get("/login") {
        }

        get("/callback") {
        }
    }
}

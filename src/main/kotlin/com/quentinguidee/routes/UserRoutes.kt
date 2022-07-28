package com.quentinguidee.routes

import com.quentinguidee.services.userService
import com.quentinguidee.utils.user
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.util.*

fun Route.userRoutes() {
    route("/user") {
        get {
            call.respond(call.user.toJSON())
        }

        get("/{username}") {
            val username = call.parameters.getOrFail("username")

            val user = userService.get(username)

            call.respond(user.toJSON())
        }
    }
}

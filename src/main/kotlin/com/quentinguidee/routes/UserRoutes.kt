package com.quentinguidee.routes

import com.quentinguidee.services.usersServices
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

            val user = usersServices.get(username)

            call.respond(user.toJSON())
        }
    }
}

package com.quentinguidee.routes

import com.quentinguidee.services.userService
import com.quentinguidee.utils.session
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Route.userRoutes() {
    get {
        call.respond(call.session.user.toJSON())
    }

    get("/{username}") {
        val username = call.parameters["username"] ?: return@get call.respondText(
            "missing username",
            status = HttpStatusCode.BadRequest
        )

        val user = userService.get(username) ?: return@get call.respondText(
            "user not found",
            status = HttpStatusCode.NotFound
        )

        call.respond(user.toJSON())
    }
}

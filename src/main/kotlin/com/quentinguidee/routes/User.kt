package com.quentinguidee.routes

import com.quentinguidee.services.userService
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Route.userRouting() {
    get("/{username}") {
        val username = call.parameters["username"] ?: return@get call.respondText(
            "missing username",
            status = HttpStatusCode.BadRequest
        )

        val user = userService.getUser(username) ?: return@get call.respondText(
            "user not found",
            status = HttpStatusCode.NotFound
        )

        call.respond(user.toJSON())
    }
}

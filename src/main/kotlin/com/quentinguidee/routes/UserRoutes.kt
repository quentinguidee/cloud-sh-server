package com.quentinguidee.routes

import com.quentinguidee.models.UserSession
import com.quentinguidee.services.userService
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.sessions.*

fun Route.userRoutes() {
    get {
        val session = call.sessions.get<UserSession>() ?: return@get call.respondText(
            "not authenticated",
            status = HttpStatusCode.Unauthorized
        )

        val user = userService.get(session.username) ?: return@get call.respondText(
            "user not found",
            status = HttpStatusCode.NotFound
        )

        call.respond(user.toJSON())
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

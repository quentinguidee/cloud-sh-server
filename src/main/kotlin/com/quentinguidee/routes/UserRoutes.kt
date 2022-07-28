package com.quentinguidee.routes

import com.quentinguidee.services.sessionServices
import com.quentinguidee.services.userService
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Route.userRoutes() {
    get {
        val token = call.request.header(HttpHeaders.Authorization) ?: return@get call.respondText(
            "missing authentication token",
            status = HttpStatusCode.BadRequest
        )

        val session = sessionServices.session(token) ?: return@get call.respondText(
            "user session not found",
            status = HttpStatusCode.NotFound
        )

        call.respond(session.user.toJSON())
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

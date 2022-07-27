package com.quentinguidee.routes

import com.quentinguidee.models.UserSession
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.sessions.*

fun Route.bucketRoutes() {
    get {
        val session = call.sessions.get<UserSession>() ?: call.respondText(
            "failed to retrieve user session",
            status = HttpStatusCode.InternalServerError
        )
        call.respond("not implemented")
    }
}

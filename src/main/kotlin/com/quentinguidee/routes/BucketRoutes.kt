package com.quentinguidee.routes

import com.quentinguidee.services.bucketService
import com.quentinguidee.services.sessionServices
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Route.bucketRoutes() {
    get {
        val token = call.request.header(HttpHeaders.Authorization) ?: return@get call.respondText(
            "missing authentication token",
            status = HttpStatusCode.BadRequest
        )

        val session = sessionServices.session(token) ?: return@get call.respondText(
            "user session not found",
            status = HttpStatusCode.NotFound
        )

        val userID = session.user.id.value

        var bucket = bucketService.bucket(userID)
        if (bucket == null) {
            bucket = bucketService.createBucket(userID)
        }

        call.respond(bucket.toJSON())
    }
}

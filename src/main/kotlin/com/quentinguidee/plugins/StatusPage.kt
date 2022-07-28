package com.quentinguidee.plugins

import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.plugins.*
import io.ktor.server.plugins.statuspages.*
import io.ktor.server.response.*
import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put

fun Application.configureStatusPage() {
    install(StatusPages) {
        exception<Throwable> { call, cause ->
            val code = when (cause) {
                is BadRequestException -> HttpStatusCode.BadRequest
                is NoSuchElementException -> HttpStatusCode.NotFound
                else -> HttpStatusCode.InternalServerError
            }

            val message = buildJsonObject {
                put("message", cause.message.toString())
            }

            call.respondText(message.toString(), status = code)
        }
    }
}

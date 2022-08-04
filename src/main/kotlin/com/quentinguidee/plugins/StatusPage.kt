package com.quentinguidee.plugins

import com.quentinguidee.utils.DatabaseConnectionFailedException
import com.quentinguidee.utils.NotAuthenticatedException
import com.quentinguidee.utils.ServerAlreadyConfiguredException
import com.quentinguidee.utils.UnauthorizedException
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
                is DatabaseConnectionFailedException -> HttpStatusCode.NotFound
                is NoSuchElementException -> HttpStatusCode.NotFound
                is ServerAlreadyConfiguredException -> HttpStatusCode.Unauthorized
                is UnauthorizedException -> HttpStatusCode.Unauthorized
                is NotAuthenticatedException -> HttpStatusCode.Unauthorized
                else -> HttpStatusCode.InternalServerError
            }

            val message = buildJsonObject {
                put("type", cause.javaClass.simpleName)
                put("message", cause.message.toString())
            }

            call.respondText(message.toString(), status = code)
        }
    }
}

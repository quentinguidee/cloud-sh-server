package com.quentinguidee.utils

import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.response.*
import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put

suspend fun ApplicationCall.ok() {
    respondText(buildJsonObject {
        put("message", "OK")
    }.toString(), status = HttpStatusCode.OK)
}

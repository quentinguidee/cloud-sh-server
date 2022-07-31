package com.quentinguidee.routes.storage

import com.quentinguidee.models.AccessType
import com.quentinguidee.services.storage.bucketsServices
import com.quentinguidee.services.storage.nodesServices
import com.quentinguidee.utils.UnauthorizedException
import com.quentinguidee.utils.user
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.util.*
import java.util.*

fun Route.recentRoutes() {
    get("/recent") {
        val bucketUUID = call.parameters.getOrFail("bucket_uuid")
        val userID = call.user.id

        if (!bucketsServices.authorize(AccessType.READ, UUID.fromString(bucketUUID), userID))
            throw UnauthorizedException(call.user)

        val nodes = nodesServices.getRecent(UUID.fromString(bucketUUID), userID)

        call.respond(nodes)
    }
}

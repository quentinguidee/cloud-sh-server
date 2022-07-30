package com.quentinguidee.routes.storage

import com.quentinguidee.models.AccessType
import com.quentinguidee.services.storage.bucketsServices
import com.quentinguidee.services.storage.nodesServices
import com.quentinguidee.utils.UnauthorizedException
import com.quentinguidee.utils.ok
import com.quentinguidee.utils.user
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.util.*
import java.util.*

fun Route.binRoutes() {
    route("/bin") {
        get {
            val bucketUUID = call.parameters.getOrFail("bucket_uuid")

            if (!bucketsServices.authorize(AccessType.READ, UUID.fromString(bucketUUID), call.user.id))
                throw UnauthorizedException(call.user)

            val nodes = nodesServices.getBin(UUID.fromString(bucketUUID))

            call.respond(nodes)
        }

        delete {
            val bucketUUID = call.parameters.getOrFail("bucket_uuid")

            if (!bucketsServices.authorize(AccessType.WRITE, UUID.fromString(bucketUUID), call.user.id))
                throw UnauthorizedException(call.user)

            nodesServices.emptyBin(UUID.fromString(bucketUUID))

            call.ok()
        }
    }
}

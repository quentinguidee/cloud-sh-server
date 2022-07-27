package com.quentinguidee.models.db

import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put
import org.jetbrains.exposed.dao.IntEntity
import org.jetbrains.exposed.dao.IntEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import org.jetbrains.exposed.dao.id.IntIdTable
import org.jetbrains.exposed.sql.transactions.transaction

object Sessions : IntIdTable() {
    val user = reference("user_id", Users)
    val token = varchar("token", 63)
}

class Session(id: EntityID<Int>) : IntEntity(id) {
    companion object : IntEntityClass<Session>(Sessions)

    var user by User referencedOn Sessions.user
    var token by Sessions.token

    fun toJSON() = transaction {
        return@transaction buildJsonObject {
            put("user", user.toJSON())
            put("token", token)
        }
    }
}

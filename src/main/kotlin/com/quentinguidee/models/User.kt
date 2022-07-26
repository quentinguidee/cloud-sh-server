package com.quentinguidee.models

import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put
import org.jetbrains.exposed.dao.IntEntity
import org.jetbrains.exposed.dao.IntEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import org.jetbrains.exposed.dao.id.IntIdTable

object Users : IntIdTable() {
    val username = varchar("username", 127).uniqueIndex()
    val name = varchar("name", 127)
    val email = varchar("email", 255)
    val profilePicture = varchar("profile_picture", 255)
    val role = varchar("role", 63).nullable()
}

class User(id: EntityID<Int>) : IntEntity(id) {
    companion object : IntEntityClass<User>(Users)

    val username by Users.username
    val name by Users.name
    val email by Users.email
    val profilePicture by Users.profilePicture
    val role by Users.role

    fun toJSON(): JsonObject {
        return buildJsonObject {
            put("id", id.value)
            put("username", username)
            put("name", name)
            put("email", email)
            put("profile_picture", profilePicture)
            put("role", role)
        }
    }
}
